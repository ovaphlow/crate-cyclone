from django.db import connection

def list_schemas():
    with connection.cursor() as cursor:
        cursor.execute("select schema_name from information_schema.schemata")
        return [row[0] for row in cursor.fetchall()]

def list_tables(schema):
    with connection.cursor() as cursor:
        cursor.execute("select table_name from information_schema.tables where table_schema = %s", (schema,))
        return [row[0] for row in cursor.fetchall()]
