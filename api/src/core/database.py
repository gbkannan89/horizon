import os
import logging
from psycopg_pool import ConnectionPool
from .config import settings

logger = logging.getLogger(__name__)

# Connection pool instance
pool = None

def init_pool():
    global pool
    if pool is None:
        logger.info("Initializing psycopg connection pool...")
        pool = ConnectionPool(
            conninfo=settings.DATABASE_URL,
            min_size=1,
            max_size=10,
            open=True
        )

def close_pool():
    global pool
    if pool is not None:
        logger.info("Closing psycopg connection pool...")
        pool.close()
        pool = None

def init_db():
    init_pool()
    # Find schema.sql path
    current_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    schema_path = os.path.join(current_dir, "db", "schema.sql")
    
    logger.info(f"Reading database schema from {schema_path}...")
    try:
        with open(schema_path, "r", encoding="utf-8") as f:
            schema_sql = f.read()
        
        with pool.connection() as conn:
            with conn.cursor() as cur:
                cur.execute(schema_sql)
            conn.commit()
        logger.info("Database tables initialized successfully.")
    except Exception as e:
        logger.error(f"Failed to initialize database: {e}")
        raise e

def get_db():
    if pool is None:
        init_pool()
    with pool.connection() as conn:
        yield conn
