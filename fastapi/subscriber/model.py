from sqlalchemy import Column, BigInteger, String, DateTime

from infrastructure.postgres import Base


class Subscriber(Base):
    __tablename__ = "subscriber"
    __table_args__ = {"schema": "crate"}

    id = Column(BigInteger, primary_key=True, index=True)
    relation_id = Column(BigInteger)
    reference_id = Column(BigInteger)
    time = Column(DateTime)
    email = Column(String)
    name = Column(String)
    phone = Column(String)
    tags = Column(String)
    detail = Column(String)
    state = Column(String)
