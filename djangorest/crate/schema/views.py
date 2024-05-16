from rest_framework.response import Response
from rest_framework.views import APIView

from .service import get_all_schemas, get_all_tables, get_columns, save


class SchemaEndpoint(APIView):
    @staticmethod
    def get(request):
        schemas = get_all_schemas()
        return Response(schemas, status=200)


class TableEndpoint(APIView):
    @staticmethod
    def get(request: object, schema: object) -> object:
        tables = get_all_tables(schema)
        return Response(tables, status=200)


class ColumnEndpoint(APIView):
    @staticmethod
    def get(request: object, schema: str, table: str) -> object:
        columns = get_columns(schema, table)
        return Response(columns, status=200)


class ListEndpoint(APIView):
    @staticmethod
    def get(request: object, schema: str, table: str):
        return Response('Not implemented', status=501)

    @staticmethod
    def post(request, schema, table):
        data = request.data
        save(schema, table, data)
        return Response('Created', status=201)
