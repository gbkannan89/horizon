from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime

class GoalOut(BaseModel):
    id: int
    user_id: int
    name: str
    target_amount: float
    current_amount: float
    status: str
    color: str
    created_at: datetime
    projected_completion_date: Optional[str] = None
    monthly_saving_needed: Optional[float] = None
