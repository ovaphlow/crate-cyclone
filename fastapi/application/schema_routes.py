import asyncio
import json
from datetime import datetime, timezone

from fastapi import APIRouter, Depends, Request, Response
from sqlalchemy.orm import Session

from application.event import publisher, Event
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


async def save_event(relation_id: int, reference_id: int, detail: str):
    event = Event(relation_id=relation_id, reference_id=reference_id, detail=detail,
                  time=datetime.now(timezone.utc).isoformat())
    await publisher.dispatch(event.dict())


@router.get("/db-schema")
async def schemas(request: Request, db: Session = Depends(get_db)):
    repo = SchemaRepository(db)
    service = SchemaService(repo)
    asyncio.ensure_future(save_event(0, 0, json.dumps(dict(ip=request.client.host))))
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
