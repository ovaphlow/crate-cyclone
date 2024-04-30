package crate.schema.repository;

import io.vertx.core.Future;
import io.vertx.sqlclient.Row;
import io.vertx.sqlclient.RowSet;

import java.util.List;

public interface SchemaRepository {

    Future<RowSet<Row>> listSchemas();

    List<String> listTables(String schema);
}
