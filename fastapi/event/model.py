from datetime import datetime

from sqlalchemy import Column, BigInteger, String, DateTime

from infrastructure.postgres import Base


class Event(Base):
    __tablename__ = "event"
    __table_args__ = {"schema": "crate"}

    id = Column(BigInteger, primary_key=True, index=True)
    relation_id = Column(BigInteger)
    reference_id = Column(BigInteger)
    detail = Column(String)
    time = Column(DateTime, default=datetime.now().isoformat())
