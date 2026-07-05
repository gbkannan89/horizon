import json
import logging
from fastapi import APIRouter, Depends, HTTPException, status
from ..core.database import get_db
from ..schemas.auth import UserOut
from .auth import get_current_user
from ..services.engine import calculate_financial_health_score

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/score", tags=["score"])

# GET Health Score (Protected)
@router.get("/health")
def get_health_score(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    if not current_user.household_id:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="User does not have a household"
        )
        
    score, max_possible, breakdown = calculate_financial_health_score(current_user.household_id, db)
    
    with db.cursor() as cur:
        try:
            cur.execute(
                """
                INSERT INTO health_scores (household_id, score, max_possible_score, pillar_scores)
                VALUES (%s, %s, %s, %s)
                RETURNING id, calculated_at
                """,
                (current_user.household_id, score, max_possible, json.dumps(breakdown["pillars"]))
            )
            row = cur.fetchone()
            db.commit()
            calculated_at = row[1]
        except Exception as e:
            db.rollback()
            logger.error(f"Failed to log health score to database: {e}")
            # Non-blocking error, return calculated score anyway
            import datetime
            calculated_at = datetime.datetime.utcnow()
            
    return {
        "score": score,
        "maxPossibleScore": max_possible,
        "calculatedAt": calculated_at,
        "breakdown": breakdown
    }

# GET Health Score History (Protected)
@router.get("/history")
def get_health_score_history(current_user: UserOut = Depends(get_current_user), db = Depends(get_db)):
    if not current_user.household_id:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="User does not have a household"
        )
        
    with db.cursor() as cur:
        cur.execute(
            """
            SELECT id, score, max_possible_score, pillar_scores, calculated_at 
            FROM health_scores 
            WHERE household_id = %s 
            ORDER BY calculated_at DESC 
            LIMIT 50
            """,
            (current_user.household_id,)
        )
        rows = cur.fetchall()
        
        return [
            {
                "id": r[0],
                "score": r[1],
                "maxPossibleScore": r[2],
                "pillars": r[3],
                "calculatedAt": r[4]
            }
            for r in rows
        ]
