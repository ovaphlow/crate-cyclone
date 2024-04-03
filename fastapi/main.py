from fastapi import FastAPI, Request
from sqlalchemy import create_engine, text
from sqlalchemy.engine.url import URL

app = FastAPI()

url = URL(
    drivername="mysql+pymysql",
    username="root",
    password="",
    host="11",
    port=3306,
    database="crate",
    query={"charset": "utf8mb4"},
)
engine = create_engine(url, pool_size=3, max_overflow=2, pool_recycle=30)

event_columns = ["id", "relation_id", "reference_id", "tags", "detail", "time"]

@app.get("/crate-api/event")
def read_root(request: Request):
    query_params = dict(request.query_params)
    if query_params["option"] == "default":
        equal = query_params.get("equal", "[]").split(",")
        with engine.connect() as connection:
            conditions = []
            params = []
            if len(equal) > 0 and len(equal) % 2 == 0:
                for i in range(0, len(equal), 2):
                    conditions.append(equal[i] + " = :param" + str(i//2))
                    params.append({"param" + str(i//2): equal[i + 1]})
            q = "select " + ", ".join(event_columns) + " from events"
            if len(conditions) > 0:
                q = q + " where " + " and ".join(conditions)
            result = connection.execute(text(q), **{k: v for d in params for k, v in d.items()})
            data = []
            for row in result:
                row_dict = {
                    "id": row.id,
                    "relation_id": row.relation_id,
                    "reference_id": row.reference_id,
                    "tags": row.tags,
                    "detail": row.detail,
                    "time": row.time,
                }
                data.append(row_dict)
            return data
    return []
