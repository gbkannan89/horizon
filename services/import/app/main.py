from fastapi import FastAPI
from app.routers import upload, health

app = FastAPI(
    title="Horizon Import Service",
    description="CSV, Excel, PDF, OCR processing. No business logic.",
    version="0.1.0-dev",
)

app.include_router(health.router)
app.include_router(upload.router)
