import datetime

from pydantic import BaseModel


class SubscriberBase(BaseModel):
    email: str
    name: str
    phone: str
    tags: str
    detail: str
    relation_id: int
    reference_id: int


class Subscriber(SubscriberBase):
    id: int
    time: datetime.datetime
    state: str

    class Config:
        orm_mode = True
