from rest_framework.response import Response
from rest_framework.views import APIView
from .service import get_all_schemas, get_all_tables

class SchemaEndpoint(APIView):
    def get(self, request):
        schemas = get_all_schemas()
        return Response(schemas, status=200)

class TableEndpoint(APIView):
    def get(self, request, schema):
        tables = get_all_tables(schema)
        return Response(tables, status=200)
