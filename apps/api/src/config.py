import os

class Settings:
    PROJECT_NAME: str = "FinScore API"
    DATABASE_URL: str = os.getenv(
        "DATABASE_URL",
        "postgresql://postgres:password@localhost:5435/horizon"
    )
    JWT_SECRET_KEY: str = os.getenv("JWT_SECRET_KEY", "finscore_jwt_access_secret_2026")
    JWT_REFRESH_SECRET_KEY: str = os.getenv("JWT_REFRESH_SECRET_KEY", "finscore_jwt_refresh_secret_2026")
    ACCESS_TOKEN_EXPIRE_MINUTES: int = 15
    REFRESH_TOKEN_EXPIRE_DAYS: int = 7

settings = Settings()
