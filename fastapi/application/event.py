from datetime import datetime

from pydantic import BaseModel

from event.service import EventService


class Event(BaseModel):
    relation_id: int
    reference_id: int
    detail: str
    time: datetime


class Singleton(type):
    _instance = {}

    def __call__(cls, *args, **kwargs):
        if cls not in cls._instance:
            cls._instance[cls] = super(Singleton, cls).__call__(*args, **kwargs)
        return cls._instance[cls]


class Publisher(metaclass=Singleton):
    def __init__(self):
        self.subscribers = set()

    def register(self, who):
        self.subscribers.add(who)

    def unregister(self, who):
        self.subscribers.discard(who)

    async def dispatch(self, message):
        for it in self.subscribers:
            await it.update(message)


class Subscriber(metaclass=Singleton):
    def __init__(self, name):
        self.name = name
        self.event_service = EventService()

    async def update(self, message):
        print(f"触发订阅事件 {self.name} {message}")
        self.event_service.create(message["relation_id"], message["reference_id"], message["detail"], message["time"])


publisher = Publisher()
subscriber = Subscriber("event")
publisher.register(subscriber)
