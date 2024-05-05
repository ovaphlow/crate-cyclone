import os

from dotenv import load_dotenv
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

load_dotenv()

DATABASE_URL = "postgresql://" + os.getenv("PG_USER") + ":" + os.getenv("PG_PASSWORD") + "@" + os.getenv("PG_HOST") + ":" + os.getenv("PG_PORT") + "/" + os.getenv("PG_DATABASE")

# DATABASE_URL = "postgresql://ovaphlow:@172.24.192.1:5432/postgres"

engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
