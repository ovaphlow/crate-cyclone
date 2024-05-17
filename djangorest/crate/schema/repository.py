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


def retrieve(schema:str, table:str, columns: list, filters: list, options: dict):
    try:
        take = options.get("take", 10)
        skip = options.get("page", 1) * take - take
        where = ""
        conditions = []
        params = {}
        for f in filters:
            if f[0] == "equal":
                if f[1] is None:
                    continue
                if f[2] is None:
                    continue
                conditions.append(f"{f[1]} = %s")
                params[f[1]] = f[2]
        if len(conditions) > 0:
            where = "where " + " and ".join(conditions)
        with connection.cursor() as cursor:
            cursor.execute(
                f'''
                select {", ".join(columns)} from {schema}.{table}
                {where}
                order by id desc
                limit %s offset %s
                ''',
                tuple(list(params.values()) + [take, skip])
            )
            rows = cursor.fetchall()
            return [dict(zip(columns, row)) for row in rows] if rows else None
    except Exception as e:
        raise e
