from infrastructure.postgres import SessionLocal
from setting.model import Setting


class SettingRepository:
    def __init__(self, db: SessionLocal = None):
        self.db = db

    def create(self, setting: Setting):
        pass
