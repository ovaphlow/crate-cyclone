package crate.schema;

import crate.schema.SchemaRepository;
import crate.schema.SchemaService;
import io.vertx.core.Future;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

import static org.mockito.Mockito.*;

public class SchemaServiceTest {

    @Mock
    private SchemaRepository mockRepo;

    private SchemaService service;

    @BeforeEach
    public void setup() {
        MockitoAnnotations.openMocks(this);
        service = new SchemaService(mockRepo);
    }

    @Test
    @DisplayName("Should list schemas successfully")
    public void listSchemasSuccessfully() {
        when(mockRepo.retrieveSchemas()).thenReturn(Future.succeededFuture(Arrays.asList("schema1", "schema2")));

        service.listSchemas();

        verify(mockRepo, times(1)).retrieveSchemas();
    }

    @Test
    @DisplayName("Should handle no schemas")
    public void handleNoSchemas() {
        when(mockRepo.retrieveSchemas()).thenReturn(Future.succeededFuture(Collections.emptyList()));

        service.listSchemas();

        verify(mockRepo, times(1)).retrieveSchemas();
    }

    @Test
    @DisplayName("Should list tables successfully")
    public void listTablesSuccessfully() {
        when(mockRepo.retrieveTables(anyString())).thenReturn(Future.succeededFuture(Arrays.asList("table1", "table2")));

        service.listTables("schema1");

        verify(mockRepo, times(1)).retrieveTables("schema1");
    }

    @Test
    @DisplayName("Should handle no tables")
    public void handleNoTables() {
        when(mockRepo.retrieveTables(anyString())).thenReturn(Future.succeededFuture(Collections.emptyList()));

        service.listTables("schema1");

        verify(mockRepo, times(1)).retrieveTables("schema1");
    }

    @Test
    @DisplayName("Should list columns successfully")
    public void listColumnsSuccessfully() {
        Map<String, String> column1 = new HashMap<>();
        column1.put("name", "column1");
        column1.put("type", "integer");

        Map<String, String> column2 = new HashMap<>();
        column2.put("name", "column2");
        column2.put("type", "varchar");

        when(mockRepo.retrieveColumns(anyString(), anyString())).thenReturn(Future.succeededFuture(Arrays.asList(column1, column2)));

        service.listColumns("schema1", "table1");

        verify(mockRepo, times(1)).retrieveColumns("schema1", "table1");
    }

    @Test
    @DisplayName("Should handle no columns")
    public void handleNoColumns() {
        when(mockRepo.retrieveColumns(anyString(), anyString())).thenReturn(Future.succeededFuture(Collections.emptyList()));

        service.listColumns("schema1", "table1");

        verify(mockRepo, times(1)).retrieveColumns("schema1", "table1");
    }
}
