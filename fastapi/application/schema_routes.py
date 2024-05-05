from sqlalchemy.orm import Session

from fastapi import APIRouter, Depends, Response, Request
from infrastructure.postgres import SessionLocal
from schema.repository import SchemaRepository
from schema.service import SchemaService

router = APIRouter()


def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


@router.get("/db-schema")
async def schemas(db: Session = Depends(get_db)):
    repo = SchemaRepository(db)
    service = SchemaService(repo)
    return service.list_schemas()


@router.get("/{schema}/db-table")
async def tables(db: Session = Depends(get_db), schema: str = None):
    repo = SchemaRepository(db)
    service = SchemaService(repo)
    return service.list_tables(schema)


@router.get("/{schema}/{table}/db-column")
async def columns(db: Session = Depends(get_db), schema: str = None, table: str = None):
    repo = SchemaRepository(db)
    service = SchemaService(repo)
    return service.list_columns(schema, table)


@router.post("/{schema}/{table}")
async def save_data(request: Request, db: Session = Depends(get_db), data: dict = None, schema: str = None, table: str = None):
    repo = SchemaRepository(db)
    service = SchemaService(repo)
    try:
        service.save_data(schema, table, data)
        return Response(status_code=201)
    except:
        return Response(status_code=500, content=dict(
            type="about:blank",
            status=500,
            title="Internal Server Error",
            detail="An error occurred while saving data",
            instance=str(request.url)
        ))
