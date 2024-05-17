from sqlalchemy import Column, BigInteger, String, DateTime

from infrastructure.postgres import Base


class Setting(Base):
    __tablename__ = "setting"
    __table_args__ = {"schema": "crate"}

    id = Column(BigInteger, primary_key=True, index=True)
    root_id = Column(BigInteger)
    parent_id = Column(BigInteger)
    tags = Column(String)
    detail = Column(String)
    created_at = Column(DateTime)
    updated_at = Column(DateTime)
    state = Column(String)
