import datetime

from pydantic import BaseModel


class SettingBase(BaseModel):
    root_id: int
    parent_id: int
    tags: str
    detail: str


class Setting(SettingBase):
    id: int
    created_at: datetime.datetime
    updated_at: datetime.datetime
    state: str

    class Config:
        orm_mode = True
