from fastapi import APIRouter, UploadFile, File

router = APIRouter(prefix="/import", tags=["import"])

@router.post("/csv")
async def upload_csv(file: UploadFile = File(...)):
    return {"status": "not_implemented", "filename": file.filename}

@router.post("/excel")
async def upload_excel(file: UploadFile = File(...)):
    return {"status": "not_implemented", "filename": file.filename}

@router.post("/pdf")
async def upload_pdf(file: UploadFile = File(...)):
    return {"status": "not_implemented", "filename": file.filename}
