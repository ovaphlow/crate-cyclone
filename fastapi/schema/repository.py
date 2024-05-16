from sqlalchemy import text
from sqlalchemy.orm import Session


class SchemaRepository:
    def __init__(self, db: Session):
        self.db = db

    def get_all_schemas(self):
        return self.db.execute(text("select schema_name from information_schema.schemata")).fetchall()

    def get_tables(self, schema):
        q = text("select table_name from information_schema.tables where table_schema = :schema")
        return self.db.execute(q, params={"schema": schema}).fetchall()

    def get_columns(self, schema, table):
        q = text("""
                 select column_name, data_type from information_schema.columns
                 where table_schema = :schema and table_name = :table
                 """)
        return self.db.execute(q, params={"schema": schema, "table": table}).fetchall()

    def create(self, schema, table, data):
        columns = self.get_columns(schema, table)
        column_names = ", ".join(column[0] for column in columns)
        placeholders = ", ".join(":" + column[0] for column in columns)
        self.db.execute(
            text(f"""
                 insert into {schema}.{table} ({column_names}) values ({placeholders})
                 """),
            params=data
        )

    def retrieve(self, columns: list, schema: str, table: str, filters: list, options: dict):
        take: int = options.get("take", 10)
        skip: int = options.get("page", 1) * take - take
        where: str = ""
        conditions: list = []
        params: dict = {}
        for f in filters:
            if f[0] == "equal":
                if f[1] is None:
                    continue
                if f[2] is None:
                    continue
                conditions.append(f"{f[1]} = :{f[1]}")
                params[f[1]] = f[2]
        if len(conditions) > 0:
            where = "where " + " and ".join(conditions)
        q = text(f"""
                 select {", ".join(columns)} from {schema}.{table}
                 {where}
                 order by id desc
                 limit :take offset :skip
                 """)
        params['take'] = take
        params['skip'] = skip
        return self.db.execute(q, params=params).fetchall()
