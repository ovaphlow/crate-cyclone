from fastapi import APIRouter, Depends, Request
from sqlalchemy.orm import Session

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


@router.get("")
async def get(request: Request, db: Session = Depends(get_db)):
    equal = request.query_params.get("equal", None)
    filters: list = []
    if equal:
        p = equal.split(',')
        if len(p) % 2 == 0:
            filters.append(['equal', i, p[p.index(i) + 1]] for i in p[::2])
    repo = SchemaRepository(db)
    service = SchemaService(repo)
    take = int(request.query_params.get("take", "10"))
    page = int(request.query_params.get("page", "1"))
    options = dict(take=take, page=page)
    return service.retrieve_data('crate', 'setting', filters, options)
