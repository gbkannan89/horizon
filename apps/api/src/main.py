import time
import logging
from fastapi import FastAPI, Depends
from fastapi.middleware.cors import CORSMiddleware
from .database import init_db, close_pool, get_db
from .routes.auth import router as auth_router
from .routes.household import router as household_router
from .routes.income import router as income_router
from .routes.assets import router as assets_router
from .routes.liabilities import router as liabilities_router
from .routes.bills import router as bills_router
from .routes.score import router as score_router
from .routes.dashboard import router as dashboard_router
from .routes.discipline import router as discipline_router
from .routes.transactions import router as transactions_router
from .routes.insurance import router as insurance_router

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(title="FinScore API")
app.include_router(auth_router)
app.include_router(household_router)
app.include_router(income_router)
app.include_router(assets_router)
app.include_router(liabilities_router)
app.include_router(bills_router)
app.include_router(score_router)
app.include_router(dashboard_router)
app.include_router(discipline_router)
app.include_router(transactions_router)
app.include_router(insurance_router)

# Add CORS Middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.on_event("startup")
async def startup_event():
    logger.info("Starting up FinScore API...")
    try:
        init_db()
        logger.info("Startup complete.")
    except Exception as e:
        logger.error(f"Error during startup: {e}")

@app.on_event("shutdown")
async def shutdown_event():
    logger.info("Shutting down FinScore API...")
    close_pool()
    logger.info("Shutdown complete.")

@app.get("/health")
def health_check(conn = Depends(get_db)):
    start_time = time.time()
    try:
        with conn.cursor() as cur:
            cur.execute("SELECT 1;")
            cur.fetchone()
        db_status = "connected"
    except Exception as e:
        logger.error(f"Database health check failed: {e}")
        db_status = f"error: {str(e)}"
        
    return {
        "status": "OK",
        "timestamp": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "database": db_status,
        "latency_ms": round((time.time() - start_time) * 1000, 2)
    }

@app.get("/")
def read_root():
    return {
        "message": "FinScore Personal Finance Intelligence API",
        "healthCheck": "/health"
    }
