from fastapi import APIRouter

router = APIRouter(tags=["health"])

@router.get("/health/live")
async def live():
    return {"status": "live"}

@router.get("/health/ready")
async def ready():
    return {"status": "ready"}
