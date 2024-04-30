package crate.schema.service;

import crate.schema.repository.SchemaRepositoryImpl;
import io.vertx.core.Future;

import java.util.List;
import java.util.stream.StreamSupport;

public class SchemaService {

    private final SchemaRepositoryImpl repo;

    public SchemaService(SchemaRepositoryImpl repo) {
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
