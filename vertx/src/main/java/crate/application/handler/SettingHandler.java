package crate.application.handler;

import crate.HandlerSchema;
import crate.infrastructure.ErrorResponse;
import crate.setting.SettingService;
import io.vertx.core.Future;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.Router;
import io.vertx.sqlclient.Row;
import io.vertx.sqlclient.RowSet;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.Map;
import java.util.stream.StreamSupport;

public class SettingHandler {
    private final Logger logger = LoggerFactory.getLogger(SettingHandler.class);
    private final SettingService service;

    public SettingHandler(SettingService service) {
        this.service = service;
    }

    public void setupRoutes(Router router) {
        router.get("/crate-api/setting").handler(context -> {
            String option = context.request().getParam("option", "");
            if ("default".equals(option)) {
//                Future<RowSet<Row>> future = new HandlerSchema(this.pool)
//                    .retrieve(List.of("*"),
//                        "crate",
//                        "setting",
//                        List.of(),
//                        Map.of("take", context.request().getParam("take", "20"),
//                            "page", context.request().getParam("page", "1")));
//                future.onSuccess(rows -> {
//                    JsonArray response = new JsonArray(StreamSupport.stream(rows.spliterator(), false)
//                        .map(row -> new JsonObject()
//                            .put("id", row.getLong("id"))
//                            .put("rootId", row.getLong("root_id"))
//                            .put("parentId", row.getLong("parent_id"))
//                            .put("tags", row.getJsonArray("tags").toString())
//                            .put("detail", row.getJsonObject("detail").toString())
//                            .put("state", row.getJsonObject("state").toString())
//                            .put("_id", row.getLong("id").toString())
//                            .put("_rootId", row.getLong("root_id").toString())
//                            .put("_parentId", row.getLong("parent_id").toString()))
//                        .toList());
//                    context.response()
//                        .putHeader("content-type", "application/json")
//                        .end(response.encode());
//                }).onFailure(err -> {
//                    logger.error(err.getMessage());
//                    JsonObject response = JsonObject.mapFrom(new ErrorResponse(null, 500, "Internal Server Error", err.getMessage(), context.request().uri()));
//                    context.response().setStatusCode(500).end(response.encode());
//                });
//                JsonArray response = new JsonArray(result);
//                context.response()
//                    .putHeader("content-type", "application/json")
//                    .end(response.encode());
                return;
            }
            JsonObject response = JsonObject.mapFrom(new ErrorResponse(null, 406, "参数错误", "", context.request().uri()));
            context.response().setStatusCode(406).end(response.encode());
        });
    }
}
