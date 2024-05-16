from django.urls import path

from .views import SchemaEndpoint, TableEndpoint, ColumnEndpoint, ListEndpoint

urlpatterns = [
    path('schema', SchemaEndpoint.as_view(), name='schema'),
    path('<str:schema>/table', TableEndpoint.as_view(), name='table'),
    path('<str:schema>/<str:table>/column', ColumnEndpoint.as_view(), name='column'),
    path('<str:schema>/<str:table>', ListEndpoint.as_view(), name='list'),
]
