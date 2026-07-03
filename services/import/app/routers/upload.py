from fastapi import APIRouter, File, UploadFile

from app.services.excel_parser import parse_excel
from app.services.import_service import import_csv
from app.services.pdf_parser import parse_pdf

router = APIRouter(prefix="/import", tags=["import"])


@router.post("/csv")
async def upload_csv(file: UploadFile = File(...)):
    content = await file.read()
    result = import_csv(content)
    return {
        "status": result.status.value,
        "filename": file.filename,
        "total": result.total,
        "events_created": result.events_created,
        "errors": result.errors,
        "transactions": [t.model_dump() for t in result.transactions],
    }


@router.post("/excel")
async def upload_excel(file: UploadFile = File(...)):
    content = await file.read()
    transactions = parse_excel(content)
    return {
        "status": "completed",
        "filename": file.filename,
        "total": len(transactions),
        "events_created": len(transactions),
        "errors": [],
        "transactions": [t.model_dump() for t in transactions],
    }


@router.post("/pdf")
async def upload_pdf(file: UploadFile = File(...)):
    content = await file.read()
    transactions = parse_pdf(content)
    return {
        "status": "completed",
        "filename": file.filename,
        "total": len(transactions),
        "events_created": len(transactions),
        "errors": [],
        "transactions": [t.model_dump() for t in transactions],
    }
