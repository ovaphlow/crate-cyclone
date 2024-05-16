from .repository import SchemaRepository


class SchemaService:
    def __init__(self, repo: SchemaRepository):
        self.repo = repo

    def list_schemas(self):
        result = [row[0] for row in self.repo.get_all_schemas()]
        return result

    def list_tables(self, schema):
        result = [row[0] for row in self.repo.get_tables(schema)]
        return result

    def list_columns(self, schema, table):
        result = [dict(zip(["column_name", "data_type"], row)) for row in self.repo.get_columns(schema, table)]
        return result

    def save_data(self, schema, table, data):
        self.repo.create(schema, table, data)

    def retrieve_data(self, schema, table, filters, options):
        columns = self.list_columns(schema, table)
        col_list = [column["column_name"] for column in columns]
        result = [dict(zip(col_list, row)) for row in self.repo.retrieve(col_list, schema, table, filters, options)]
        return result
