from datetime import datetime

from event.model import Event
from event.repository import EventRepository


class EventService:
    def __init__(self):
        self.repo = EventRepository()

    def create(self, relation_id: int, reference_id: int, detail: str, time: datetime):
        event = Event(relation_id=relation_id, reference_id=reference_id, detail=detail, time=time)
        return self.repo.create(event)
