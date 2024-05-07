package crate.schema;

import io.vertx.core.Future;

import java.util.List;
import java.util.stream.StreamSupport;

public class SchemaService {

    private final SchemaRepository repo;

    public SchemaService(SchemaRepository repo) {
        this.repo = repo;
    }

    public Future<List<String>> listSchemas() {
        return repo.listSchemas()
            .map(rows -> StreamSupport.stream(rows.spliterator(), false)
                .map(row -> row.getString("schema_name"))
                .toList());
    }

    public List<String> listTables(String schema) {
        return repo.listTables(schema);
    }
}
