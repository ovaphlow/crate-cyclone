package crate;

import crate.infrastructure.ErrorResponse;
import io.vertx.core.Future;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.RoutingContext;
import io.vertx.sqlclient.Pool;
import io.vertx.sqlclient.Row;
import io.vertx.sqlclient.RowSet;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.Map;
import java.util.stream.StreamSupport;

public class HandlerSetting {

    private static final Logger logger = LoggerFactory.getLogger(HandlerSetting.class);

    private final Pool pool;

    public HandlerSetting(Pool pool) {
        this.pool = pool;
    }

    public void setupRoutes(Router router) {
    }

    private void retrieve(RoutingContext context) {
        String option = context.request().getParam("option", "");
        if ("default".equals(option)) {
            Future<RowSet<Row>> future = new HandlerSchema(this.pool)
                .retrieve(List.of("*"),
                    "crate",
                    "setting",
                    List.of(),
                    Map.of("take", context.request().getParam("take", "20"),
                        "page", context.request().getParam("page", "1")));
            future.onSuccess(rows -> {
                JsonArray response = new JsonArray(StreamSupport.stream(rows.spliterator(), false)
                    .map(row -> new JsonObject()
                        .put("id", row.getLong("id"))
                        .put("rootId", row.getLong("root_id"))
                        .put("parentId", row.getLong("parent_id"))
                        .put("tags", row.getJsonArray("tags").toString())
                        .put("detail", row.getJsonObject("detail").toString())
                        .put("state", row.getJsonObject("state").toString())
                        .put("_id", row.getLong("id").toString())
                        .put("_rootId", row.getLong("root_id").toString())
                        .put("_parentId", row.getLong("parent_id").toString()))
                    .toList());
                context.response()
                    .putHeader("content-type", "application/json")
                    .end(response.encode());
            }).onFailure(err -> {
                logger.error(err.getMessage());
                JsonObject response = JsonObject.mapFrom(new ErrorResponse.Builder()
                    .status(500)
                    .title("Internal Server Error")
                    .detail(err.getMessage())
                    .instance(context.request().uri())
                    .build());
                context.response().setStatusCode(500).end(response.encode());
            });
        } else {
            JsonObject response = JsonObject.mapFrom(new ErrorResponse.Builder()
                .status(406)
                .title("Not Acceptable")
                .detail("参数错误")
                .instance(context.request().uri())
                .build());
            context.response().setStatusCode(406).end(response.encode());
        }
    }
}
