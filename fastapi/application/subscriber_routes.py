import asyncio
import json

from fastapi import APIRouter, Depends, Request
from sqlalchemy.orm import Session
from starlette.responses import Response

from application.event import save_event
from infrastructure.postgres import SessionLocal
from subscriber.repository import SubscriberRepository
from subscriber.service import SubscriberService

router = APIRouter()


def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


@router.post("/sign-up")
async def sign_up(request: Request, db: Session = Depends(get_db)):
    repo = SubscriberRepository(db)
    service = SubscriberService(repo)
    body = await request.json()
    username = body.get("username")
    password = body.get("password")
    if not username or not password:
        return dict(type="about:blank",
                    status=400,
                    title="参数错误",
                    detail="The request body is invalid",
                    instance=str(request.url))
    subscriber = service.sign_up(username, password)
    if not subscriber:
        return dict(type="about:blank",
                    status=409,
                    title="用户已存在",
                    detail="The username already exists",
                    instance=str(request.url))
    asyncio.ensure_future(save_event(subscriber.id, 0, json.dumps(dict(ip=request.client.host, event="sign-up"))))
    return Response(status_code=201)
