package crate.application.handler;

import crate.infrastructure.ErrorResponse;
import crate.schema.SchemaService;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.Router;

import java.util.HashMap;
import java.util.Map;

public class SchemaHandler {
    private final SchemaService service;

    public SchemaHandler(SchemaService service) {
        this.service = service;
    }

    public void setupRoutes(Router router) {
        router.get("/crate-api/database/schema").handler(context -> {
            service.listSchemas()
                .onSuccess(schemas -> {
                    JsonArray response = new JsonArray(schemas);
                    context.response()
                        .putHeader("content-type", "application/json")
                        .end(response.encode());
                })
                .onFailure(err -> {
                    JsonObject response = JsonObject.mapFrom(new ErrorResponse(null, 500, "服务器错误", err.getMessage(), context.request().uri()));
                    context.response()
                        .setStatusCode(500)
                        .putHeader("content-type", "application/json")
                        .end(response.encode());
                });
        });

        router.get("/crate-api/database/:schema/table").handler(context -> {
            service.listTables(context.pathParam("schema"))
                .onSuccess(tables -> {
                    JsonArray response = new JsonArray(tables);
                    context.response()
                        .putHeader("content-type", "application/json")
                        .end(response.encode());
                })
                .onFailure(err -> {
                    JsonObject response = JsonObject.mapFrom(new ErrorResponse(null, 500, "服务器错误", err.getMessage(), context.request().uri()));
                    context.response()
                        .setStatusCode(500)
                        .putHeader("content-type", "application/json")
                        .end(response.encode());
                });
        });

        router.get("/crate-api/database/:schema/:table/column").handler(context -> {
            service.listColumns(context.pathParam("schema"), context.pathParam("table"))
                .onSuccess(columns -> {
                    JsonArray response = new JsonArray(columns);
                    context.response()
                        .putHeader("content-type", "application/json")
                        .end(response.encode());
                })
                .onFailure(err -> {
                    JsonObject response = JsonObject.mapFrom(new ErrorResponse(null, 500, "服务器错误", err.getMessage(), context.request().uri()));
                    context.response()
                        .setStatusCode(500)
                        .putHeader("content-type", "application/json")
                        .end(response.encode());
                });
        });

        router.post("/crate-api/database/:schema/:table").handler(context -> {
            JsonObject body = context.body().asJsonObject();
            service.save(context.pathParam("schema"), context.pathParam("table"), body.getMap())
                .onSuccess(result -> context.response().setStatusCode(201).end())
                .onFailure(err -> {
                    JsonObject response = JsonObject.mapFrom(new ErrorResponse(null, 500, "服务器错误", err.getMessage(), context.request().uri()));
                    context.response()
                        .setStatusCode(500)
                        .putHeader("content-type", "application/json")
                        .end(response.encode());
                });
        });

        router.get("/crate-api/database/:schema/:table/:id").handler(context -> {
            service.getById(context.pathParam("schema"), context.pathParam("table"), Long.parseLong(context.pathParam("id")), context.request().getParam("uuid"))
                .onSuccess(result -> {
                    Map<String, Object> r = new HashMap<>(result);
                    for (Map.Entry<String, Object> entry : result.entrySet()) {
                        if (entry.getValue() instanceof Long) {
                            r.put("_" + entry.getKey(), entry.getValue().toString());
                        }
                    }
                    JsonObject response = JsonObject.mapFrom(r);
                    context.response()
                        .putHeader("content-type", "application/json")
                        .end(response.encode());
                })
                .onFailure(err -> {
                    JsonObject response = JsonObject.mapFrom(new ErrorResponse(null, 500, "服务器错误", err.getMessage(), context.request().uri()));
                    context.response()
                        .setStatusCode(500)
                        .putHeader("content-type", "application/json")
                        .end(response.encode());
                });
        });

        router.put("/crate-api/database/:schema/:table/:id").handler(context -> {
            service.update(context.pathParam("schema"), context.pathParam("table"),
                    Long.parseLong(context.pathParam("id")),
                    context.request().getParam("uuid"),
                    context.body().asJsonObject().getMap())
                .onSuccess(result -> context.response().setStatusCode(204).end())
                .onFailure(err -> {
                    JsonObject response = JsonObject.mapFrom(new ErrorResponse(null, 500, "服务器错误", err.getMessage(), context.request().uri()));
                    context.response()
                        .setStatusCode(500)
                        .putHeader("content-type", "application/json")
                        .end(response.encode());
                });
        });
    }
}
