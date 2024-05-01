from django.urls import path
from .views import SchemaEndpoint, TableEndpoint

urlpatterns = [
    path('db-schema', SchemaEndpoint.as_view(), name='schema'),
    path('<str:schema>/db-table', TableEndpoint.as_view(), name='table')
]
