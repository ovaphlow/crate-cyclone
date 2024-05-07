package crate.application.handler;

import crate.infrastructure.ErrorResponse;
import crate.setting.SettingService;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.Router;

import java.util.List;
import java.util.Map;

public class SettingHandler {
    private final SettingService service;

    public SettingHandler(SettingService service) {
        this.service = service;
    }

    public void setupRoutes(Router router) {
        router.get("/crate-api/setting").handler(context -> {
            String option = context.request().getParam("option", "");
            int take = Integer.parseInt(context.request().getParam("take", "20"));
            long page = Long.parseLong(context.request().getParam("page", "1"));
            List<Map<String, Object>> result = service.defaultList(Map.of("take", take, "page", page));
            if ("default".equals(option)) {
                JsonArray response = new JsonArray(result);
                context.response()
                    .putHeader("content-type", "application/json")
                    .end(response.encode());
                return;
            }
            JsonObject response = JsonObject.mapFrom(new ErrorResponse.Builder()
                .status(406)
                .title("Not Acceptable")
                .detail("参数错误")
                .instance(context.request().uri())
                .build());
            context.response().setStatusCode(406).end(response.encode());
        });
    }
}
