from pydantic import BaseModel
from enum import Enum

class ImportStatus(str, Enum):
    PENDING = "pending"
    PROCESSING = "processing"
    COMPLETED = "completed"
    FAILED = "failed"

class ImportResult(BaseModel):
    status: ImportStatus
    events_created: int = 0
    errors: list[str] = []
