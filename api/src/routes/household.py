import random
import string
import logging
from fastapi import APIRouter, Depends, HTTPException, status
from ..core.database import get_db
from ..schemas.auth import UserOut
from ..schemas.household import JoinHousehold, ContributingMemberCreate, ContributingMemberUpdate, FamilyInviteRequest
from .auth import get_current_user

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/household", tags=["household"])

def generate_unique_invite_code(db) -> str:
    with db.cursor() as cur:
        while True:
            code = "FIN-" + "".join(random.choices(string.ascii_uppercase + string.digits, k=6))
            cur.execute("SELECT id FROM households WHERE invite_code = %s", (code,))
            if not cur.fetchone():
                return code

# GET Household Summary (Protected)
@router.get("/summary")
def get_household_summary(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    if not current_user.household_id:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User does not belong to any household"
        )
        
    with db.cursor() as cur:
        # Fetch household info
        cur.execute(
            "SELECT id, name, invite_code FROM households WHERE id = %s",
            (current_user.household_id,)
        )
        household_row = cur.fetchone()
        if not household_row:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Household not found"
            )
            
        household_id, household_name, invite_code = household_row
        
        # Fetch members (linked partners)
        cur.execute(
            "SELECT id, email, name, user_type, risk_profile FROM users WHERE household_id = %s",
            (household_id,)
        )
        member_rows = cur.fetchall()
        members = [
            {
                "id": r[0],
                "email": r[1],
                "name": r[2],
                "userType": r[3],
                "riskProfile": r[4]
            }
            for r in member_rows
        ]
        
        # Fetch contributing members
        cur.execute(
            "SELECT id, name, monthly_income, contribution_to_household, relationship FROM contributing_members WHERE household_id = %s",
            (household_id,)
        )
        contrib_rows = cur.fetchall()
        contributing_members = [
            {
                "id": r[0],
                "name": r[1],
                "monthlyIncome": float(r[2]),
                "contributionToHousehold": float(r[3]),
                "relationship": r[4]
            }
            for r in contrib_rows
        ]
        
        return {
            "id": household_id,
            "name": household_name,
            "inviteCode": invite_code,
            "members": members,
            "contributingMembers": contributing_members
        }

# POST Create/Get Invite Code (Protected)
@router.post("/invite")
def create_invite_code(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    if not current_user.household_id:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="User does not have a household"
        )
        
    with db.cursor() as cur:
        # Check if code already exists
        cur.execute(
            "SELECT invite_code FROM households WHERE id = %s",
            (current_user.household_id,)
        )
        row = cur.fetchone()
        if row and row[0]:
            return {"inviteCode": row[0]}
            
        # Generate and save new code
        code = generate_unique_invite_code(db)
        cur.execute(
            "UPDATE households SET invite_code = %s WHERE id = %s",
            (code, current_user.household_id)
        )
        db.commit()
        
        return {"inviteCode": code}

# POST Join Household (Protected)
@router.post("/join")
def join_household(join_in: JoinHousehold, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    invite_code = join_in.invite_code.strip().upper()
    
    with db.cursor() as cur:
        # Find target household
        cur.execute(
            "SELECT id, name FROM households WHERE UPPER(invite_code) = %s",
            (invite_code,)
        )
        target_household = cur.fetchone()
        if not target_household:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Invalid invitation code"
            )
            
        target_id, target_name = target_household
        old_household_id = current_user.household_id
        
        if old_household_id == target_id:
            return {"message": "You are already a member of this household", "householdName": target_name}
            
        try:
            # Update user's household reference
            cur.execute(
                "UPDATE users SET household_id = %s WHERE id = %s",
                (target_id, current_user.id)
            )
            db.commit()
            
            # Clean up old empty household
            if old_household_id:
                cur.execute(
                    "SELECT COUNT(id) FROM users WHERE household_id = %s",
                    (old_household_id,)
                )
                users_left = cur.fetchone()[0]
                
                if users_left == 0:
                    logger.info(f"Household {old_household_id} has no users left. Deleting it.")
                    cur.execute("DELETE FROM households WHERE id = %s", (old_household_id,))
                    db.commit()
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to join household: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Failed to join household due to database error"
            )
        return {
            "message": f"Successfully joined household: {target_name}",
            "householdId": target_id,
            "householdName": target_name
        }

# POST Invite Family Member by Email (Protected)
@router.post("/invite-by-email", status_code=status.HTTP_201_CREATED)
def invite_family_member(invite: FamilyInviteRequest, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    if not current_user.household_id:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="User does not have a household"
        )

    with db.cursor() as cur:
        # Check if email already has a pending invite in this household
        cur.execute(
            "SELECT id FROM family_invites WHERE household_id = %s AND email = %s",
            (current_user.household_id, invite.email)
        )
        if cur.fetchone():
            raise HTTPException(
                status_code=status.HTTP_409_CONFLICT,
                detail="This email already has a pending invitation to your household"
            )

        # Check if email is already a registered user
        cur.execute("SELECT id FROM users WHERE email = %s", (invite.email,))
        existing_user = cur.fetchone()
        if existing_user:
            raise HTTPException(
                status_code=status.HTTP_409_CONFLICT,
                detail="A user with this email already exists. Ask them to sign in and use your invite code."
            )

        # Get or generate invite code for the household
        cur.execute("SELECT invite_code FROM households WHERE id = %s", (current_user.household_id,))
        row = cur.fetchone()
        invite_code = row[0] if row and row[0] else generate_unique_invite_code(db)

        if not row or not row[0]:
            cur.execute(
                "UPDATE households SET invite_code = %s WHERE id = %s",
                (invite_code, current_user.household_id)
            )

        # Store the pending invite
        cur.execute(
            """
            INSERT INTO family_invites (household_id, email, name, relationship, invite_code)
            VALUES (%s, %s, %s, %s, %s)
            """,
            (current_user.household_id, invite.email, invite.name, invite.relationship, invite_code)
        )
        db.commit()

        # Also add as a contributing member so their contribution is tracked
        cur.execute(
            """
            INSERT INTO contributing_members (household_id, name, monthly_income, contribution_to_household, relationship)
            VALUES (%s, %s, 0, 0, %s)
            RETURNING id
            """,
            (current_user.household_id, invite.name, invite.relationship or "family")
        )
        db.commit()

    return {
        "message": f"Invitation sent to {invite.email}",
        "inviteCode": invite_code,
        "householdName": current_user.name or current_user.email
    }


# POST Add Contributing Member (Protected)
@router.post("/contributing-members", status_code=status.HTTP_201_CREATED)
def add_contributing_member(member_in: ContributingMemberCreate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    if not current_user.household_id:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="User does not have a household"
        )
    
    with db.cursor() as cur:
        cur.execute(
            """
            INSERT INTO contributing_members (household_id, name, monthly_income, contribution_to_household, relationship)
            VALUES (%s, %s, %s, %s, %s)
            RETURNING id, name, monthly_income, contribution_to_household, relationship
            """,
            (current_user.household_id, member_in.name, member_in.monthly_income, member_in.contribution_to_household, member_in.relationship)
        )
        row = cur.fetchone()
        db.commit()
        
        return {
            "id": row[0],
            "name": row[1],
            "monthlyIncome": float(row[2]),
            "contributionToHousehold": float(row[3]),
            "relationship": row[4]
        }

# PUT Update Contributing Member (Protected)
@router.put("/contributing-members/{member_id}")
def update_contributing_member(member_id: int, member_in: ContributingMemberUpdate, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    if not current_user.household_id:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="User does not have a household"
        )
        
    with db.cursor() as cur:
        # Check if member belongs to current user's household
        cur.execute(
            "SELECT id FROM contributing_members WHERE id = %s AND household_id = %s",
            (member_id, current_user.household_id)
        )
        if not cur.fetchone():
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Contributing member not found in your household"
            )
            
        # Build dynamic update statement
        updates = []
        params = []
        if member_in.name is not None:
            updates.append("name = %s")
            params.append(member_in.name)
        if member_in.monthly_income is not None:
            updates.append("monthly_income = %s")
            params.append(member_in.monthly_income)
        if member_in.contribution_to_household is not None:
            updates.append("contribution_to_household = %s")
            params.append(member_in.contribution_to_household)
        if member_in.relationship is not None:
            updates.append("relationship = %s")
            params.append(member_in.relationship)
            
        if not updates:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="No fields to update"
            )
            
        params.append(member_id)
        query = f"UPDATE contributing_members SET {', '.join(updates)} WHERE id = %s RETURNING id, name, monthly_income, contribution_to_household, relationship"
        
        try:
            cur.execute(query, tuple(params))
            row = cur.fetchone()
            db.commit()
            return {
                "id": row[0],
                "name": row[1],
                "monthlyIncome": float(row[2]),
                "contributionToHousehold": float(row[3]),
                "relationship": row[4]
            }
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to update contributing member: {e}")
            raise HTTPException(
                status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
                detail="Failed to update contributing member"
            )

# DELETE Remove Contributing Member (Protected)
@router.delete("/contributing-members/{member_id}", status_code=status.HTTP_200_OK)
def delete_contributing_member(member_id: int, current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    if not current_user.household_id:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="User does not have a household"
        )
        
    with db.cursor() as cur:
        # Check if member belongs to current user's household
        cur.execute(
            "SELECT id FROM contributing_members WHERE id = %s AND household_id = %s",
            (member_id, current_user.household_id)
        )
        if not cur.fetchone():
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail="Contributing member not found in your household"
            )
            
        cur.execute("DELETE FROM contributing_members WHERE id = %s", (member_id,))
        db.commit()
        return {"message": "Contributing member removed successfully"}

