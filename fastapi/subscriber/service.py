import base64
import hashlib
import json
import os

from subscriber.model import Subscriber
from subscriber.repository import SubscriberRepository


class SubscriberService:
    def __init__(self, repo: SubscriberRepository):
        self.repo = repo

    def sign_up(self, username: str, password: str):
        subs = self.list_by_username(username)
        if len(subs) > 0:
            return None
        salt = os.urandom(16).hex()
        key = hashlib.pbkdf2_hmac("sha256", password.encode("utf-8"), salt.encode(), 100000, dklen=32)
        hashed_password = base64.b64encode(key).decode()
        detail = dict(salt=salt, hash=hashed_password)
        subscriber = Subscriber(relation_id=0, reference_id=0,
                                tags="[]", detail=json.dumps(detail), email=username, name=username, phone="")
        return self.repo.create(subscriber)

    def list_by_username(self, username: str):
        return self.repo.retrieve_by_username(username)
