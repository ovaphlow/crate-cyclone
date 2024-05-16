from .repository import retrieve_schemas, retrieve_tables, retrieve_columns, create


def get_all_schemas():
    return retrieve_schemas()


def get_all_tables(schema):
    return retrieve_tables(schema)


def get_columns(schema: str, table: str):
    return retrieve_columns(schema, table)


def save(schema: str, table: str, data: dict):
    return create(schema, table, data)
