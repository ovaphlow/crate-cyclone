package crate.schema;

import cn.hutool.core.util.IdUtil;
import io.vertx.core.Future;
import io.vertx.core.Promise;
import io.vertx.core.json.JsonObject;
import io.vertx.sqlclient.Pool;
import io.vertx.sqlclient.Tuple;

import java.time.OffsetDateTime;
import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.StreamSupport;

public class SchemaRepository {

    private final Pool pool;

    public SchemaRepository(Pool pool) {
        this.pool = pool;
    }

    public Future<List<String>> retrieveSchemas() {
        Promise<List<String>> promise = Promise.promise();
        pool.query("select schema_name from information_schema.schemata")
            .execute()
            .onSuccess(rows -> promise.complete(StreamSupport.stream(rows.spliterator(), false)
                .map(row -> row.getString("schema_name"))
                .collect(Collectors.toList())))
            .onFailure(promise::fail);
        return promise.future();
    }

    public Future<List<String>> retrieveTables(String schema) {
        Promise<List<String>> promise = Promise.promise();
        pool.preparedQuery("select table_name from information_schema.tables where table_schema = $1")
            .execute(Tuple.of(schema))
            .onSuccess(rows -> promise.complete(StreamSupport.stream(rows.spliterator(), false)
                .map(row -> row.getString("table_name"))
                .collect(Collectors.toList())))
            .onFailure(promise::fail);
        return promise.future();
    }

    public Future<List<Map<String, String>>> retrieveColumns(String schema, String table) {
        Promise<List<Map<String, String>>> promise = Promise.promise();
        pool.preparedQuery("""
                select column_name, data_type from information_schema.columns
                where table_schema = $1 and table_name = $2
                order by ordinal_position asc
                """)
            .execute(Tuple.of(schema, table))
            .onSuccess(rows -> promise.complete(StreamSupport.stream(rows.spliterator(), false)
                .map(row -> Map.of("column_name", row.getString("column_name"),
                    "data_type", row.getString("data_type")))
                .collect(Collectors.toList())))
            .onFailure(promise::fail);
        return promise.future();
    }

    public Future<Void> create(String schema, String table, Map<String, Object> data) {
        Promise<Void> promise = Promise.promise();
        System.out.println(data);
        this.retrieveColumns(schema, table)
            .onSuccess(columnList -> {
                List<String> columnNames = columnList.stream().map(column -> column.get("column_name")).toList();
                if (data.keySet().stream().noneMatch(columnNames::contains)) {
                    throw new IllegalArgumentException("Invalid column name");
                }
                OffsetDateTime time = OffsetDateTime.now();
                if (columnNames.contains("created_at")) {
                    data.put("created_at", time);
                }
                if (columnNames.contains("updated_at")) {
                    data.put("updated_at", time);
                }
                if (columnNames.contains("state")) {
                    Map<String, Object> state = Map.of("uuid", UUID.randomUUID().toString());
                    JsonObject stateJson = new JsonObject(state);
                    data.put("state", stateJson);
                }
                data.put("id", IdUtil.getSnowflake(1, 1).nextId());
                StringBuilder columns = new StringBuilder();
                StringBuilder values = new StringBuilder();
                Tuple params = Tuple.tuple();
                for (Map.Entry<String, Object> entry : data.entrySet()) {
                    columns.append(entry.getKey()).append(",");
                    Optional<Map<String, String>> col = columnList.stream()
                        .filter(column -> column.get("column_name").equals(entry.getKey()))
                        .findFirst();
                    if (col.isPresent() && "jsonb".equals(col.get().get("data_type")) && entry.getValue() instanceof String) {
                        values.append("'").append(entry.getValue()).append("'::jsonb,");
                    } else {
                        values.append("$").append(params.size() + 1).append(",");
                        params.addValue(entry.getValue());
                    }
                }
                String query = String.format("insert into %s.%s (%s) values (%s)", schema, table,
                    columns.substring(0, columns.length() - 1), values.substring(0, values.length() - 1));
                System.out.println(query);
                pool.preparedQuery(query).execute(params)
                    .onSuccess(result -> promise.complete())
                    .onFailure(promise::fail);
            })
            .onFailure(promise::fail);
        return promise.future();
    }

    public Future<Map<String, Object>> retrieveById(String schema, String table, long id, String uuid) {
        Promise<Map<String, Object>> promise = Promise.promise();
        pool.preparedQuery("select * from " + schema + "." + table + " where id = $1 and state->>'uuid' = $2")
            .execute(Tuple.of(id, uuid))
            .onSuccess(rows -> {
                if (rows.size() == 0) {
                    promise.fail("无数据");
                } else {
                    promise.complete(rows.iterator().next().toJson().getMap());
                }
            })
            .onFailure(promise::fail);
        return promise.future();
    }

    public Future<Void> update(String schema, String table, long id, String uuid, Map<String, Object> data) {
        Promise<Void> promise = Promise.promise();
        this.retrieveColumns(schema, table)
            .onSuccess(columnList -> {
                List<String> columnNames = columnList.stream().map(column -> column.get("column_name")).toList();
                if (columnList.isEmpty()) {
                    throw new IllegalArgumentException("表不存在");
                }
                if (data.keySet().stream().noneMatch(columnNames::contains)) {
                    throw new IllegalArgumentException("字段错误");
                }
                OffsetDateTime time = OffsetDateTime.now();
                if (columnNames.contains("updated_at")) {
                    data.put("updated_at", time);
                }
                StringBuilder set = new StringBuilder();
                Tuple params = Tuple.tuple();
                for (Map.Entry<String, Object> entry : data.entrySet()) {
                    set.append(entry.getKey()).append(" = ");
                    Optional<Map<String, String>> col = columnList.stream()
                        .filter(column -> column.get("column_name").equals(entry.getKey()))
                        .findFirst();
                    if (col.isPresent() && "jsonb".equals(col.get().get("data_type")) && entry.getValue() instanceof String) {
                        set.append("'").append(entry.getValue()).append("'::jsonb,");
                    } else {
                        set.append("$").append(params.size() + 1).append(",");
                        params.addValue(entry.getValue());
                    }
                }
                String query = String.format("update %s.%s set %s where id = $%d and state->>'uuid' = $%d",
                    schema, table, set.substring(0, set.length() - 1), params.size() + 1, params.size() + 2);
                params.addValue(id);
                params.addValue(uuid);
                pool.preparedQuery(query).execute(params)
                    .onSuccess(result -> promise.complete())
                    .onFailure(promise::fail);
            })
            .onFailure(promise::fail);
        return promise.future();
    }
}
