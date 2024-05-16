from django.db import connection


def retrieve_schemas():
    try:
        with connection.cursor() as cursor:
            cursor.execute("select schema_name from information_schema.schemata")
            rows = cursor.fetchall()
            return [row[0] for row in rows] if rows else None
    except Exception as e:
        raise e


def retrieve_tables(schema: str):
    try:
        with connection.cursor() as cursor:
            cursor.execute("select table_name from information_schema.tables where table_schema = %s", (schema,))
            rows = cursor.fetchall()
            return [row[0] for row in rows] if rows else None
    except Exception as e:
        raise e


def retrieve_columns(schema: str, table: str):
    try:
        with connection.cursor() as cursor:
            cursor.execute(
                '''
                select column_name, data_type from information_schema.columns
                where table_schema = %s and table_name = %s
                ''',
                (schema, table)
            )
            rows = cursor.fetchall()
            return [{'column_name': row[0], 'data_type': row[1]} for row in rows] if rows else None
    except Exception as e:
        raise e


def create(schema: str, table: str, data: dict):
    try:
        with connection.cursor() as cursor:
            columns = ', '.join(data.keys())
            values = ', '.join(['%s'] * len(data))
            cursor.execute(
                f'''
                insert into {schema}.{table} ({columns})
                values ({values})
                ''',
                tuple(data.values())
            )
    except Exception as e:
        raise e
