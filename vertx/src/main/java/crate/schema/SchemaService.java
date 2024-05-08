package crate.schema;

import io.vertx.core.Future;
import io.vertx.core.Promise;

import java.util.List;
import java.util.Map;

public class SchemaService {

    private final SchemaRepository repo;

    public SchemaService(SchemaRepository repo) {
        this.repo = repo;
    }

    public Future<List<String>> listSchemas() {
        return repo.retrieveSchemas();
    }

    public Future<List<String>> listTables(String schema) {
        return repo.retrieveTables(schema);
    }

    public Future<List<Map<String, String>>> listColumns(String schema, String table) {
        return repo.retrieveColumns(schema, table);
    }

    public Future<Void> save(String schema, String table, Map<String, Object> data) {
        Promise<Void> promise = Promise.promise();
        repo.create(schema, table, data).onSuccess(promise::complete).onFailure(promise::fail);
        return promise.future();
    }

    public Future<Map<String, Object>> getById(String schema, String table, long id, String uuid) {
        return repo.retrieveById(schema, table, id, uuid);
    }

    public Future<Void> update(String schema, String table, long id, String uuid, Map<String, Object> data) {
        return repo.update(schema, table, id, uuid, data);
    }
}
