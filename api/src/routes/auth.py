import logging
import hashlib
from datetime import datetime, timedelta
from typing import Optional
from fastapi import APIRouter, Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from jose import JWTError, jwt
import bcrypt

from ..core.config import settings
from ..core.database import get_db
from ..schemas.auth import UserRegister, UserLogin, UserOut, Token, TokenRefresh, TokenData, DeleteAccountRequest

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/auth", tags=["auth"])

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="/api/auth/swagger-login")

def verify_password(plain_password: str, hashed_password: str) -> bool:
    try:
        return bcrypt.checkpw(plain_password.encode("utf-8"), hashed_password.encode("utf-8"))
    except Exception as e:
        logger.error(f"Password verification error: {e}")
        return False

def get_password_hash(password: str) -> str:
    pwd_bytes = password.encode("utf-8")
    salt = bcrypt.gensalt()
    return bcrypt.hashpw(pwd_bytes, salt).decode("utf-8")

def hash_token(token: str) -> str:
    return hashlib.sha256(token.encode("utf-8")).hexdigest()

def create_access_token(data: dict) -> str:
    to_encode = data.copy()
    expire = datetime.utcnow() + timedelta(minutes=settings.ACCESS_TOKEN_EXPIRE_MINUTES)
    to_encode.update({"exp": expire})
    return jwt.encode(to_encode, settings.JWT_SECRET_KEY, algorithm="HS256")

def create_refresh_token(data: dict) -> str:
    to_encode = data.copy()
    expire = datetime.utcnow() + timedelta(days=settings.REFRESH_TOKEN_EXPIRE_DAYS)
    to_encode.update({"exp": expire})
    return jwt.encode(to_encode, settings.JWT_REFRESH_SECRET_KEY, algorithm="HS256")

def get_current_user(token: str = Depends(oauth2_scheme), db = Depends(get_db)) -> UserOut:
    credentials_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Could not validate credentials",
        headers={"WWW-Authenticate": "Bearer"},
    )
    try:
        payload = jwt.decode(token, settings.JWT_SECRET_KEY, algorithms=["HS256"])
        user_id: int = payload.get("userId")
        email: str = payload.get("email")
        if user_id is None or email is None:
            raise credentials_exception
        token_data = TokenData(user_id=user_id, email=email)
    except JWTError:
        raise credentials_exception
    
    with db.cursor() as cur:
        cur.execute(
            'SELECT id, email, name, phone, user_type, risk_profile, household_id FROM users WHERE id = %s',
            (token_data.user_id,)
        )
        row = cur.fetchone()
        if row is None:
            raise credentials_exception
        
        return UserOut(
            id=row[0],
            email=row[1],
            name=row[2],
            phone=row[3],
            user_type=row[4],
            risk_profile=row[5],
            household_id=row[6]
        )

# Register
@router.post("/register", response_model=Token, status_code=status.HTTP_201_CREATED)
def register(user_in: UserRegister, db = Depends(get_db)):
    with db.cursor() as cur:
        # Check if email exists
        cur.execute("SELECT id FROM users WHERE email = %s", (user_in.email,))
        if cur.fetchone():
            raise HTTPException(
                status_code=status.HTTP_409_CONFLICT,
                detail="User with this email already exists"
            )

        # Determine household
        household_id = None
        if user_in.invite_code:
            # User provided an invite code — join existing household
            code = user_in.invite_code.strip().upper()
            cur.execute(
                "SELECT id, name FROM households WHERE UPPER(invite_code) = %s",
                (code,)
            )
            household_row = cur.fetchone()
            if not household_row:
                raise HTTPException(
                    status_code=status.HTTP_404_NOT_FOUND,
                    detail="Invalid invite code. Please check and try again."
                )
            household_id = household_row[0]
        else:
            # No invite code — check if email has a pending family invite
            cur.execute(
                """
                SELECT fi.household_id, fi.invite_code, h.name, fi.name
                FROM family_invites fi
                JOIN households h ON h.id = fi.household_id
                WHERE fi.email = %s
                """,
                (user_in.email,)
            )
            pending = cur.fetchone()
            if pending:
                raise HTTPException(
                    status_code=status.HTTP_409_CONFLICT,
                    detail={
                        "message": f"This email has been invited to join {pending[3] or pending[2]}'s household. Enter the invite code to join or sign up without a code to create a new household.",
                        "householdName": pending[2],
                        "invitedBy": pending[3],
                        "inviteCode": pending[1]
                    }
                )
        
        try:
            if household_id:
                # Join existing household
                cur.execute(
                    "SELECT name FROM households WHERE id = %s",
                    (household_id,)
                )
                h_row = cur.fetchone()
                household_name = h_row[0] if h_row else "Household"
            else:
                # Create a new household
                household_name = f"{user_in.name or user_in.email}'s Household"
                cur.execute(
                    "INSERT INTO households (name) VALUES (%s) RETURNING id",
                    (household_name,)
                )
                household_id = cur.fetchone()[0]
            
            # Hash password
            password_hash = get_password_hash(user_in.password)
            
            # Insert User
            cur.execute(
                """
                INSERT INTO users (email, password_hash, name, phone, user_type, risk_profile, household_id)
                VALUES (%s, %s, %s, %s, %s, %s, %s)
                RETURNING id, email, name, phone, user_type, risk_profile, household_id
                """,
                (user_in.email, password_hash, user_in.name, user_in.phone, user_in.user_type, user_in.risk_profile, household_id)
            )
            user_row = cur.fetchone()

            # If they joined via invite code, remove the pending invite
            if user_in.invite_code:
                cur.execute(
                    "DELETE FROM family_invites WHERE email = %s",
                    (user_in.email,)
                )

            db.commit()
        except Exception as e:
            db.rollback()
            logger.error(f"Registration database transaction failed: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Could not complete registration"
            )
        
        user_out = UserOut(
            id=user_row[0],
            email=user_row[1],
            name=user_row[2],
            phone=user_row[3],
            user_type=user_row[4],
            risk_profile=user_row[5],
            household_id=user_row[6]
        )
        
        # Generate tokens
        access_token = create_access_token({"userId": user_out.id, "email": user_out.email})
        refresh_token = create_refresh_token({"userId": user_out.id, "email": user_out.email})
        
        # Save refresh token
        token_hash = hash_token(refresh_token)
        expires_at = datetime.utcnow() + timedelta(days=settings.REFRESH_TOKEN_EXPIRE_DAYS)
        
        with db.cursor() as token_cur:
            token_cur.execute(
                "INSERT INTO refresh_tokens (token_hash, user_id, expires_at) VALUES (%s, %s, %s)",
                (token_hash, user_out.id, expires_at)
            )
            db.commit()
            
        return Token(
            access_token=access_token,
            refresh_token=refresh_token,
            user=user_out
        )

# Login
@router.post("/login", response_model=Token)
def login(login_in: UserLogin, db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT id, email, password_hash, name, user_type, risk_profile, household_id FROM users WHERE email = %s",
            (login_in.email,)
        )
        row = cur.fetchone()
        if not row or not verify_password(login_in.password, row[2]):
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Incorrect email or password",
                headers={"WWW-Authenticate": "Bearer"},
            )
        
        user_out = UserOut(
            id=row[0],
            email=row[1],
            name=row[3],
            user_type=row[4],
            risk_profile=row[5],
            household_id=row[6]
        )
        
        # Generate tokens
        access_token = create_access_token({"userId": user_out.id, "email": user_out.email})
        refresh_token = create_refresh_token({"userId": user_out.id, "email": user_out.email})
        
        # Save refresh token hash
        token_hash = hash_token(refresh_token)
        expires_at = datetime.utcnow() + timedelta(days=settings.REFRESH_TOKEN_EXPIRE_DAYS)
        
        cur.execute(
            "INSERT INTO refresh_tokens (token_hash, user_id, expires_at) VALUES (%s, %s, %s)",
            (token_hash, user_out.id, expires_at)
        )
        db.commit()
        
        return Token(
            access_token=access_token,
            refresh_token=refresh_token,
            user=user_out
        )

# Swagger Login for docs authorize button
@router.post("/swagger-login")
def swagger_login(form_data: OAuth2PasswordRequestForm = Depends(), db = Depends(get_db)):
    with db.cursor() as cur:
        cur.execute(
            "SELECT id, email, password_hash, name, user_type, risk_profile, household_id FROM users WHERE email = %s",
            (form_data.username,)
        )
        row = cur.fetchone()
        if not row or not verify_password(form_data.password, row[2]):
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Incorrect email or password",
                headers={"WWW-Authenticate": "Bearer"},
            )
        
        user_out = UserOut(
            id=row[0],
            email=row[1],
            name=row[3],
            user_type=row[4],
            risk_profile=row[5],
            household_id=row[6]
        )
        
        # Generate tokens
        access_token = create_access_token({"userId": user_out.id, "email": user_out.email})
        refresh_token = create_refresh_token({"userId": user_out.id, "email": user_out.email})
        
        # Save refresh token hash
        token_hash = hash_token(refresh_token)
        expires_at = datetime.utcnow() + timedelta(days=settings.REFRESH_TOKEN_EXPIRE_DAYS)
        
        cur.execute(
            "INSERT INTO refresh_tokens (token_hash, user_id, expires_at) VALUES (%s, %s, %s)",
            (token_hash, user_out.id, expires_at)
        )
        db.commit()
        
        return {
            "access_token": access_token,
            "token_type": "bearer"
        }

# Refresh Token
@router.post("/refresh")
def refresh(refresh_in: TokenRefresh, db = Depends(get_db)):
    credentials_exception = HTTPException(
        status_code=status.HTTP_403_FORBIDDEN,
        detail="Invalid or expired refresh token"
    )
    
    try:
        payload = jwt.decode(refresh_in.refresh_token, settings.JWT_REFRESH_SECRET_KEY, algorithms=["HS256"])
        user_id: int = payload.get("userId")
        email: str = payload.get("email")
        if user_id is None or email is None:
            raise credentials_exception
    except JWTError:
        raise credentials_exception
        
    token_hash = hash_token(refresh_in.refresh_token)
    
    with db.cursor() as cur:
        # Verify active token in DB
        cur.execute(
            "SELECT user_id FROM refresh_tokens WHERE token_hash = %s AND expires_at > %s",
            (token_hash, datetime.utcnow())
        )
        row = cur.fetchone()
        if not row:
            raise credentials_exception
            
        # Delete old token
        cur.execute("DELETE FROM refresh_tokens WHERE token_hash = %s", (token_hash,))
        
        # Fetch user
        cur.execute(
            "SELECT id, email, name, user_type, risk_profile, household_id FROM users WHERE id = %s",
            (user_id,)
        )
        user_row = cur.fetchone()
        
        user_out = UserOut(
            id=user_row[0],
            email=user_row[1],
            name=user_row[2],
            user_type=user_row[3],
            risk_profile=user_row[4],
            household_id=user_row[5]
        )
        
        # Generate new tokens
        access_token = create_access_token({"userId": user_out.id, "email": user_out.email})
        refresh_token = create_refresh_token({"userId": user_out.id, "email": user_out.email})
        
        # Save new refresh token
        new_token_hash = hash_token(refresh_token)
        expires_at = datetime.utcnow() + timedelta(days=settings.REFRESH_TOKEN_EXPIRE_DAYS)
        cur.execute(
            "INSERT INTO refresh_tokens (token_hash, user_id, expires_at) VALUES (%s, %s, %s)",
            (new_token_hash, user_out.id, expires_at)
        )
        db.commit()
        
        return {
            "access_token": access_token,
            "refresh_token": refresh_token,
            "token_type": "bearer",
            "user": user_out
        }

# Logout
@router.post("/logout", status_code=status.HTTP_200_OK)
def logout(refresh_in: TokenRefresh, db = Depends(get_db)):
    token_hash = hash_token(refresh_in.refresh_token)
    with db.cursor() as cur:
        cur.execute("DELETE FROM refresh_tokens WHERE token_hash = %s", (token_hash,))
        db.commit()
    return {"message": "Successfully logged out"}

# Get Current User Profile (Protected)
@router.get("/me", response_model=UserOut)
def get_me(current_user: UserOut = Depends(get_current_user)):
    return current_user

# Delete Account (Protected)
@router.delete("/me", status_code=status.HTTP_200_OK)
def delete_account(
    delete_in: DeleteAccountRequest,
    current_user: UserOut = Depends(get_current_user),
    db = Depends(get_db)
):
    with db.cursor() as cur:
        cur.execute("SELECT password_hash FROM users WHERE id = %s", (current_user.id,))
        row = cur.fetchone()
        if not row or not verify_password(delete_in.password, row[0]):
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Incorrect password",
            )

        if current_user.household_id:
            cur.execute(
                "SELECT COUNT(*) FROM users WHERE household_id = %s AND id != %s",
                (current_user.household_id, current_user.id)
            )
            other_members = cur.fetchone()[0]
            if other_members == 0:
                cur.execute("DELETE FROM households WHERE id = %s", (current_user.household_id,))

        cur.execute("DELETE FROM users WHERE id = %s", (current_user.id,))
        db.commit()

    return {"message": "Account deleted successfully"}
