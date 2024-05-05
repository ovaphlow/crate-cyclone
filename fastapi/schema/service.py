from .repository import SchemaRepository


class SchemaService:
    def __init__(self, repo: SchemaRepository):
        self.repo = repo

    def list_schemas(self):
        result = self.repo.get_all_schemas()
        list_result = [row[0] for row in result]
        return list_result

    def list_tables(self, schema):
        rows = self.repo.get_tables(schema)
        result = [row[0] for row in rows]
        return result

    def list_columns(self, schema, table):
        rows = self.repo.get_columns(schema, table)
        result = [dict(zip(["column_name", "data_type"], row)) for row in rows]
        return result

    def save_data(self, schema, table, data):
        self.repo.save(schema, table, data)
