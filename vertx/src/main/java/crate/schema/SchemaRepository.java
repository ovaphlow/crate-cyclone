package crate.schema;

import io.vertx.core.Future;
import io.vertx.sqlclient.Pool;
import io.vertx.sqlclient.Row;
import io.vertx.sqlclient.RowSet;
import io.vertx.sqlclient.Tuple;

import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;
import java.util.stream.StreamSupport;

public class SchemaRepository {

    private final Pool pool;

    public SchemaRepository(Pool pool) {
        this.pool = pool;
    }

    public Future<RowSet<Row>> listSchemas() {
        return pool.query("select schema_name from information_schema.schemata")
            .execute();
    }

    public List<String> listTables(String schema) {
        CompletableFuture<List<String>> future = new CompletableFuture<>();
        pool.preparedQuery("select table_name from information_schema.tables where table_schema = $1")
            .execute(Tuple.of(schema))
            .onSuccess(rows -> future.complete(StreamSupport.stream(rows.spliterator(), false)
                .map(row -> row.getString("table_name"))
                .collect(Collectors.toList()))).onFailure(future::completeExceptionally);
        try {
            return future.get();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }
}
