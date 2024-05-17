from rest_framework.response import Response
from rest_framework.views import APIView

from .service import get_all_schemas, get_all_tables, get_columns, save, list_data


class SchemaEndpoint(APIView):
    @staticmethod
    def get(request):
        try:
            schemas = get_all_schemas()
            return Response(schemas, status=200)
        except Exception as e:
            response = dict({'type': 'about:blank',
                             'status': 500,
                             'title': '服务器错误',
                             'detail': str(e),
                             'instance': request.build_absolute_uri()})
            return Response(response, status=500)


class TableEndpoint(APIView):
    @staticmethod
    def get(request, schema: object) -> object:
        try:
            tables = get_all_tables(schema)
            return Response(tables, status=200)
        except Exception as e:
            response = dict({'type': 'about:blank',
                             'status': 500,
                             'title': '服务器错误',
                             'detail': str(e),
                             'instance': request.build_absolute_uri()})
            return Response(response, status=500)


class ColumnEndpoint(APIView):
    @staticmethod
    def get(request, schema: str, table: str) -> object:
        try:
            columns = get_columns(schema, table)
            return Response(columns, status=200)
        except Exception as e:
            response = dict({'type': 'about:blank',
                             'status': 500,
                             'title': '服务器错误',
                             'detail': str(e),
                             'instance': request.build_absolute_uri()})
            return Response(response, status=500)


class ListEndpoint(APIView):
    @staticmethod
    def get(request, schema: str, table: str):
        try:
            equal = request.query_params.get('equal', None)
            filters: list = []
            if equal:
                p = equal.split(',')
                if len(p) % 2 == 0:
                    filters.extend(['equal', i, p[p.index(i) + 1]] for i in p[::2])
            take = int(request.query_params.get("take", "10"))
            page = int(request.query_params.get("page", "1"))
            options = dict(take=take, page=page)
            response = list_data(schema, table, filters, options)
            return Response(response if response else [], 200)
        except Exception as e:
            response = dict({'type': 'about:blank',
                             'status': 500,
                             'title': '服务器错误',
                             'detail': str(e),
                             'instance': request.build_absolute_uri()})
            return Response(response, status=500)

    @staticmethod
    def post(request, schema, table):
        try:
            data = request.data
            save(schema, table, data)
            return Response('Created', status=201)
        except Exception as e:
            response = dict({'type': 'about:blank',
                             'status': 500,
                             'title': '服务器错误',
                             'detail': str(e),
                             'instance': request.build_absolute_uri()})
            return Response(response, status=500)
