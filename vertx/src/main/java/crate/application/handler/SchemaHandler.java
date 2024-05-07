package crate.application.handler;

import crate.infrastructure.ErrorResponse;
import crate.schema.SchemaService;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.Router;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;

public class SchemaHandler {

    private final Logger logger = LoggerFactory.getLogger(SchemaHandler.class);
    private final SchemaService service;

    public SchemaHandler(SchemaService service) {
        this.service = service;
    }

    public void setupRoutes(Router router) {
        router.get("/crate-api/db-schema").handler(ctx -> {
            service.listSchemas()
                .onSuccess(schemas -> {
                    JsonArray response = new JsonArray(schemas);
                    ctx.response()
                        .putHeader("content-type", "application/json")
                        .end(response.encode());
                })
                .onFailure(err -> {
                    logger.error("{}", err.getMessage());
                    JsonObject response = JsonObject.mapFrom(new ErrorResponse.Builder()
                        .type("about:blank")
                        .status(500)
                        .title("服务器错误")
                        .detail(err.getMessage())
                        .instance(ctx.request().uri())
                        .build());
                    ctx.response()
                        .putHeader("content-type", "application/json")
                        .setStatusCode(500)
                        .end(response.encode());
                });
        });

        router.get("/crate-api/:schema/db-table").handler(ctx -> {
            List<String> tables = service.listTables(ctx.pathParam("schema"));
            JsonArray response = new JsonArray(tables);
            ctx.response()
                .putHeader("content-type", "application/json")
                .end(response.encode());
        });
    }
}
