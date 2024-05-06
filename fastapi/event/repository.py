from event.model import Event
from infrastructure.postgres import SessionLocal
from infrastructure.snowflake import SnowflakeIdGenerator


class EventRepository:
    def __init__(self):
        self.db = SessionLocal()

    def create(self, event: Event):
        generator = SnowflakeIdGenerator(datacenter_id=1, worker_id=1)
        event.id = generator.generate()
        self.db.add(event)
        self.db.commit()
        self.db.refresh(event)
        return event
