import asyncio
import json

from fastapi import APIRouter, Depends, Request, Response
from sqlalchemy.orm import Session

from application.event import save_event
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


@router.get("/schema")
async def schemas(request: Request, db: Session = Depends(get_db)):
    repo = SchemaRepository(db)
    service = SchemaService(repo)
    asyncio.ensure_future(save_event(0, 0, json.dumps(dict(ip=request.client.host))))
    return service.list_schemas()


@router.get("/{schema}/table")
async def tables(db: Session = Depends(get_db), schema: str = None):
    repo = SchemaRepository(db)
    service = SchemaService(repo)
    return service.list_tables(schema)


@router.get("/{schema}/{table}/column")
async def columns(db: Session = Depends(get_db), schema: str = None, table: str = None):
    repo = SchemaRepository(db)
    service = SchemaService(repo)
    return service.list_columns(schema, table)


@router.post("/{schema}/{table}")
async def save_data(request: Request, db: Session = Depends(get_db), data: dict = None, schema: str = None,
                    table: str = None):
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


@router.get("/{schema}/{table}")
async def get(request: Request, db: Session = Depends(get_db), schema: str = None, table: str = None):
    equal = request.query_params.get("equal", None)
    filters: list = []
    if equal:
        p = equal.split(',')
        if len(p) % 2 == 0:
            for i in p[::2]:
                filters.append(['equal', i, p[p.index(i) + 1]])
    repo = SchemaRepository(db)
    service = SchemaService(repo)
    take = int(request.query_params.get("take", "10"))
    page = int(request.query_params.get("page", "1"))
    options = dict(take=take, page=page)
    return service.retrieve_data(schema, table, filters, options)
