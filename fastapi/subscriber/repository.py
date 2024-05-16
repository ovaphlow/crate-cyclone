import json
import uuid
from datetime import datetime

from sqlalchemy import or_

from infrastructure.postgres import SessionLocal
from infrastructure.snowflake import SnowflakeIdGenerator
from subscriber.model import Subscriber


class SubscriberRepository:
    def __init__(self, db: SessionLocal = None):
        self.db = db

    def create(self, subscriber: Subscriber):
        generator = SnowflakeIdGenerator(datacenter_id=1, worker_id=1)
        subscriber.id = generator.generate()
        subscriber.time = datetime.now().isoformat()
        state = dict(uuid=str(uuid.uuid4()), created_at=datetime.now().isoformat())
        subscriber.state = json.dumps(state)
        self.db.add(subscriber)
        self.db.commit()
        self.db.refresh(subscriber)
        return subscriber

    def retrieve_by_username(self, username: str):
        # return self.db.query(Subscriber).filter(Subscriber.username == username).first()
        return self.db.query(Subscriber).filter(
            or_(
                Subscriber.email == username,
                Subscriber.name == username,
                Subscriber.phone == username
            )
        ).all()
