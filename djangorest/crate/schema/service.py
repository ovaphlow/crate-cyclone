from .repository import retrieve_schemas, retrieve_tables, retrieve_columns, create, retrieve


def get_all_schemas():
    return retrieve_schemas()


def get_all_tables(schema):
    return retrieve_tables(schema)


def get_columns(schema: str, table: str):
    return retrieve_columns(schema, table)


def save(schema: str, table: str, data: dict):
    return create(schema, table, data)


def list_data(schema: str, table: str, filters: list, options: dict):
    columns = get_columns(schema, table)
    column_list = [column["column_name"] for column in columns]
    result = retrieve(schema, table, column_list, filters, options)
    return result if result else None
