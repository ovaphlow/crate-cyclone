package crate;

import cn.hutool.core.date.DateUtil;
import cn.hutool.core.util.IdUtil;
import crate.infrastructure.ErrorResponse;
import io.vertx.core.Future;
import io.vertx.core.Promise;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.RoutingContext;
import io.vertx.sqlclient.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;
import java.util.stream.Collectors;
import java.util.stream.StreamSupport;

public class HandlerSchema {

    private static final Logger logger = LoggerFactory.getLogger(HandlerSchema.class);

    private final Pool pool;

    public HandlerSchema(Pool pool) {
        this.pool = pool;
    }

    public void setupRoutes(Router router) {
//        router.route().handler(Middleware::logRequestHandler);
//        router.get("/crate-api/:schema/:table/:uuid/:id").handler(this::retrieve);
//        router.put("/crate-api/:schema/:table/:uuid/:id").handler(this::update);
//        router.delete("/crate-api/:schema/:table/:uuid/:id").handler(this::remove);
//        router.get("/crate-api/:schema/db-table").handler(this::retrieveTables);
//        router.post("/crate-api/:schema/:table").handler(this::create);
//        router.get("/crate-api/db-schema").handler(this::retrieveSchemas);
    }

    public Future<RowSet<Row>> retrieve(List<String> columns,
                                        String schema,
                                        String table,
                                        List<List<String>> filters,
                                        Map<String, String> options) {
        int take = Integer.parseInt(options.getOrDefault("take", "10"));
        long skip = (Long.parseLong(options.getOrDefault("page", "1")) - 1) * take;
        String where = "";
        List<String> conditions = new ArrayList<>(List.of());
        List<String> params = new ArrayList<>(List.of());
        for (List<String> filter : filters) {
            if ("equal".equals(filter.getFirst())) {
                if (Objects.isNull(filter.get(1))) continue;
                if (Objects.isNull(filter.get(2))) continue;
                conditions.add("%s = $%d".formatted(filter.get(1), params.size() + 1));
                params.add(filter.get(2));
            } else if ("objectContain".equals(filter.getFirst())) {
                if (Objects.isNull(filter.get(1))) continue;
                if (Objects.isNull(filter.get(2))) continue;
                if (Objects.isNull(filter.get(3))) continue;
                conditions.add("%s @> '{\"%s\": $%d}'::jsonb".formatted(filter.get(1), filter.get(2), params.size() + 1));
                params.add(filter.get(3));
            } else if ("arrayContain".equals(filter.getFirst())) {
                if (Objects.isNull(filter.get(1))) continue;
                if (Objects.isNull(filter.get(2))) continue;
                conditions.add("%s @> '[$%d]'::jsonb".formatted(filter.get(1), params.size() + 1));
                params.add(filter.get(2));
            } else if ("greater".equals(filter.getFirst())) {
                if (Objects.isNull(filter.get(1))) continue;
                if (Objects.isNull(filter.get(2))) continue;
                conditions.add("%s >= $%d".formatted(filter.get(1), params.size() + 1));
                params.add(filter.get(2));
            } else if ("lesser".equals(filter.getFirst())) {
                if (Objects.isNull(filter.get(1))) continue;
                if (Objects.isNull(filter.get(2))) continue;
                conditions.add("%s <= $%d".formatted(filter.get(1), params.size() + 1));
                params.add(filter.get(2));
            } else if ("like".equals(filter.getFirst())) {
                if (Objects.isNull(filter.get(1))) continue;
                if (Objects.isNull(filter.get(2))) continue;
                conditions.add("%s like $%d".formatted(filter.get(1), params.size() + 1));
                params.add(filter.get(2));
            } else if ("in".equals(filter.getFirst())) {
                if (Objects.isNull(filter.get(1))) continue;
                if (filter.size() < 3) continue;
                for (int i = 1; i < filter.size(); i++) {
                    List<String> f = new ArrayList<>();
                    for (int j = 2; j < filter.subList(2, filter.size()).size(); j++) {
                        f.add("$%d".formatted(params.size() + j + 1));
                        params.add(filter.get(j));
                    }
                    conditions.add("%s in (%s)".formatted(filter.get(1), String.join(", ", f)));
                }
            }
        }
        if (!conditions.isEmpty()) {
            where = "where " + String.join(" and ", conditions);
        }
        String query = """
            select %s from %s.%s
            %s
            order by id desc
            limit %s offset %s
            """.formatted(String.join(", ", columns), schema, table, where, take, skip);
        return pool.preparedQuery(query).execute(Tuple.wrap(params));
    }

    private Future<List<Map<String, String>>> retrieveColumns(String schema, String table) {
        Promise<List<Map<String, String>>> promise = Promise.promise();
        List<Map<String, String>> columns = new ArrayList<>(List.of());
        PreparedQuery<RowSet<Row>> query = pool.preparedQuery("""
            select column_name, data_type
            from information_schema.columns
            where table_schema = $1 and table_name = $2
            order by ordinal_position
            """);
        query.execute(Tuple.of(schema, table))
            .onSuccess(rows -> {
                columns.addAll(StreamSupport.stream(rows.spliterator(), false)
                    .map(row -> new HashMap<>(Map.of(
                        "column_name", row.getValue("column_name").toString(),
                        "data_type", row.getValue("data_type").toString()
                    )))
                    .toList());
                promise.complete(columns);
            })
            .onFailure(err -> {
                logger.error(err.getMessage());
                promise.fail(err);
            });
        return promise.future();
    }

    private void retrieve(RoutingContext context) {
        String _id = context.pathParam("id");
        Long id = Long.parseLong(_id);
        String uuid = context.pathParam("uuid");
        retrieveColumns(context.pathParam("schema"), context.pathParam("table"))
            .onSuccess(columns -> {
                if (columns.isEmpty()) {
                    JsonObject response = JsonObject.mapFrom(new ErrorResponse(null, 404, "No schema.", "", context.request().uri()));
                    context.response().setStatusCode(404).putHeader("content-type", "application/json").end(response.encode());
                    return;
                }
                pool.preparedQuery("""
                        select %s from %s.%s
                        where id = $1 and state->>'uuid' = $2
                        limit 1
                        """.formatted(String.join(", ", columns.stream().map(it -> it.get("column_name")).toList()), context.pathParam("schema"), context.pathParam("table"))
                    ).execute(Tuple.of(id, uuid))
                    .onSuccess(rows -> {
                        Row row = rows.iterator().next();
                        if (null == row) {
                            JsonObject response = JsonObject.mapFrom(new ErrorResponse(null, 404, "No data.", "", context.request().uri()));
                            context.response().setStatusCode(404).putHeader("content-type", "application/json").end(response.encode());
                            return;
                        }
                        Map<String, Object> data = row.toJson().getMap();
                        for (Map.Entry<String, Object> entry : new HashSet<>(data.entrySet())) {
                            if ("id".equals(entry.getKey()) || entry.getKey().contains("_id")) {
                                data.put("_" + entry.getKey(), entry.getValue().toString());
                            }
                        }
                        JsonObject response = JsonObject.mapFrom(data);
                        context.response().putHeader("content-type", "application/json").end(response.encode());
                    })
                    .onFailure(err -> {
                        JsonObject response = JsonObject.mapFrom(new ErrorResponse(null, 500, "服务器错误", err.getMessage(), context.request().uri()));
                        context.response().setStatusCode(500).putHeader("content-type", "application/json").end(response.encode());
                    });
            })
            .onFailure(err -> {
                logger.error(err.getMessage());
                JsonObject response = JsonObject.mapFrom(new ErrorResponse(null, 500, "服务器错误", err.getMessage(), context.request().uri()));
                context.response().setStatusCode(500).putHeader("content-type", "application/json").end(response.encode());
            });
    }

    private void update(RoutingContext context) {
        String schema = context.pathParam("schema");
        String table = context.pathParam("table");
        String uuid = context.pathParam("uuid");
        long id = Long.parseLong(context.pathParam("id"));
        JsonObject body = new JsonObject(context.body().asString());
        Future<List<Map<String, String>>> future = retrieveColumns(schema, table);
        future.compose(columns -> {
            if (columns.isEmpty()) {
                return Future.failedFuture("表 不存在");
            }
            if (!new HashSet<>(columns.stream().map(it -> it.get("column_name")).toList()).containsAll(body.fieldNames())) {
                return Future.failedFuture("字段 不匹配");
            }
            List<String> conditions = new ArrayList<>(List.of());
            List<Object> params = new ArrayList<>(List.of());
            for (String field : body.fieldNames()) {
                Optional<Map<String, String>> col = columns.stream().filter(it -> it.get("column_name").equals(field)).findFirst();
                if (col.isEmpty()) {
                    continue;
                }
                if ("jsonb".equals(col.get().get("data_type"))) {
                    conditions.add("%s = '%s'::jsonb".formatted(field, body.getValue(field).toString()));
                } else {
                    conditions.add("%s = $%d".formatted(field, params.size() + 1));
                    params.add(body.getValue(field));
                }
            }
            conditions.add("state = state || '{\"updated_at\": \"%s\"}'::jsonb".formatted(DateUtil.formatDateTime(new Date())));
            String query = """
                update %s.%s
                set %s
                where id = $%d and state->>'uuid' = $%d
                """.formatted(schema, table, String.join(", ", conditions), params.size() + 1, params.size() + 2);
            logger.info(query);
            return pool.preparedQuery(query).execute(Tuple.wrap(params).addLong(id).addString(uuid));
        }).onSuccess(result -> context.response().setStatusCode(200).end()).onFailure(err -> {
            JsonObject response = JsonObject.mapFrom(new ErrorResponse(null, 500, "服务器错误", err.getMessage(), context.request().uri()));
            context.response().setStatusCode(500).putHeader("content-type", "application/json").end(response.encode());
        });
    }

    private void remove(RoutingContext context) {
        //
    }
}
