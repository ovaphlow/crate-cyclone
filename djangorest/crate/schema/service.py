from .repository import list_schemas, list_tables

def get_all_schemas():
    return list_schemas()

def get_all_tables(schema):
    return list_tables(schema)
