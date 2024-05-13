package crate.application.handler;

import cn.hutool.core.bean.BeanUtil;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import crate.infrastructure.ErrorResponse;
import crate.setting.SettingService;
import io.vertx.core.json.JsonArray;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.Router;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class SettingHandler {
    private final Logger logger = LoggerFactory.getLogger(SettingHandler.class);
    private final SettingService service;
    private final ObjectMapper mapper;

    public SettingHandler(SettingService service) {
        this.service = service;
        this.mapper = new ObjectMapper();
        this.mapper.registerModule(new JavaTimeModule());
    }

    public void setupRoutes(Router router) {
        router.get("/crate-api/setting").handler(context -> {
            String option = context.request().getParam("option", "");
            if ("default".equals(option)) {
                service.listSettings()
                    .onSuccess(settings -> {
                        logger.info("{}", settings);
                        List<Map<String, Object>> result = new ArrayList<>();
                        for (var setting : settings) {
                            Map<String, Object> map = BeanUtil.beanToMap(setting);
                            map.put("createdAt", setting.getCreatedAt().toString());
                            map.put("updatedAt", setting.getUpdatedAt().toString());
                            map.put("_id", setting.getId().toString());
                            map.put("_rootId", setting.getRootId().toString());
                            map.put("_parentId", setting.getParentId().toString());
                            result.add(map);
                        }
                        JsonArray response = new JsonArray(result);
                        logger.info("{}", response);
                        context.response()
                            .putHeader("content-type", "application/json")
                            .end(response.encode());
                    })
                    .onFailure(err -> {
                        JsonObject response = JsonObject.mapFrom(new ErrorResponse(null, 500, "Internal Server Error", err.getMessage(), context.request().uri()));
                        context.response().setStatusCode(500).end(response.encode());
                    });
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
