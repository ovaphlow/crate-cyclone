package crate;

import cn.hutool.core.date.DateUtil;
import cn.hutool.core.util.IdUtil;
import crate.infrastructure.model.ErrorResponse;
import io.vertx.core.Future;
import io.vertx.core.Vertx;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.auth.PubSecKeyOptions;
import io.vertx.ext.auth.jwt.JWTAuth;
import io.vertx.ext.auth.jwt.JWTAuthOptions;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.RoutingContext;
import io.vertx.sqlclient.Pool;
import io.vertx.sqlclient.Row;
import io.vertx.sqlclient.RowSet;
import io.vertx.sqlclient.Tuple;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;
import java.util.*;

public class HandlerSubscriber {

    private static final String SECRET = "-=ovaphlow.crate.cyclone=-";
    private static final ArrayList<String> JWT_CLAIM_TAGS = new ArrayList<>(Arrays.asList("ovaphlow.crate", "cyclone", "subscriber"));
    private static final Logger logger = LoggerFactory.getLogger(HandlerSubscriber.class);

    private final Vertx vertx;
    private final Pool pool;

    public HandlerSubscriber(Vertx vertx, Pool pool) {
        this.vertx = vertx;
        this.pool = pool;
    }

    public void setupRoutes(Router router) {
        router.route().handler(MiddlewareLog::logRequestHandler);
        router.post("/crate-api/subscriber/sign-up").handler(this::signUp);
        router.post("/crate-api/subscriber/log-in").handler(this::logIn);
        router.post("/crate-api/subscriber/refresh").handler(this::refresh);
    }

    private String hashPassword(String password, String salt, String encoder) throws Exception {
        try {
            Mac sha256HMAC = Mac.getInstance("HmacSHA256");
            SecretKeySpec secretKey = new SecretKeySpec(salt.getBytes(StandardCharsets.UTF_8), "HmacSHA256");
            sha256HMAC.init(secretKey);

            byte[] hashedBytes = sha256HMAC.doFinal(password.getBytes(StandardCharsets.UTF_8));
            sha256HMAC.doFinal(password.getBytes(StandardCharsets.UTF_8));
            if ("hex".equals(encoder)) {
                return HexFormat.of().formatHex(hashedBytes);
            } else {
                return Base64.getEncoder().encodeToString(hashedBytes);
            }
        } catch (Exception e) {
            throw new Exception("Failed to create HMAC-SHA256 instance");
        }
    }

    private String generateToken(String id, String email, String name) {
        JWTAuth provider = JWTAuth.create(vertx, new JWTAuthOptions()
            .addPubSecKey(new PubSecKeyOptions()
                .setAlgorithm("HS256")
                .setBuffer(SECRET)));
        return provider.generateToken(new JsonObject()
            .put("iss", JWT_CLAIM_TAGS.getFirst())
            .put("exp", System.currentTimeMillis() + 1000 * 60 * 60 * 24)
            .put("sub", JWT_CLAIM_TAGS.get(1))
            .put("aud", JWT_CLAIM_TAGS.getLast())
            .put("nbf", System.currentTimeMillis())
            .put("iat", System.currentTimeMillis())
            .put("jti", IdUtil.simpleUUID())
            .put("id", id)
            .put("email", email)
            .put("name", name));
    }

    private void signUp(RoutingContext context) {
        JsonObject body = new JsonObject(context.body().asString());
        Future<RowSet<Row>> future = pool.preparedQuery("""
                select count(*) as qty from crate.subscriber
                where name = $1 or email = $2 or phone = $3
                """)
            .execute(Tuple.of(body.getString("username"), body.getString("username"), body.getString("username")));
        future.compose(rows -> {
                Row row = rows.iterator().next();
                if (row.getLong("qty") > 0) {
                    JsonObject response = JsonObject.mapFrom(ErrorResponse.builder()
                        .status(409)
                        .title("Conflict")
                        .detail("用户已存在")
                        .instance(context.request().uri())
                        .build());
                    context.response().setStatusCode(409).end(response.encode());
                    return Future.failedFuture("用户已存在");
                }
                return Future.succeededFuture();
            })
            .compose(v -> {
                Long id = IdUtil.getSnowflake(1, 1).nextId();
                String salt = IdUtil.simpleUUID();
                String detail = "{}";
                try {
                    JsonObject d = new JsonObject()
                        .put("salt", salt)
                        .put("hash", hashPassword(body.getString("password"), salt, "hex"));
                    detail = d.encode();
                } catch (Exception e) {
                    logger.error(e.getMessage());
                }
                JsonObject state = new JsonObject().put("uuid", IdUtil.randomUUID()).put("created_at", DateUtil.formatDateTime(new Date()));
                return pool.preparedQuery("""
                        insert into crate.subscriber (id, email, name, phone, tags, detail, time, state)
                        values ($1, $2, $3, $4, $5, $6, $7, $8)
                        """)
                    .execute(Tuple.of(id, body.getString("email"), body.getString("email"), "", "[]", detail, new Date(), state.encode()));
            })
            .onSuccess(result -> context.response().setStatusCode(201).end())
            .onFailure(err -> {
                JsonObject response = JsonObject.mapFrom(ErrorResponse.builder()
                    .status(500)
                    .title("Internal Server Error")
                    .detail(err.getMessage())
                    .instance(context.request().uri())
                    .build());
                context.response().setStatusCode(500).end(response.encode());
            });
        /*
        pool.preparedQuery("""
                select count(*) as qty from crate.subscriber
                where name = $1 or email = $2 or phone = $3
                """)
            .execute(Tuple.of(body.getString("username"), body.getString("username"), body.getString("username")))
            .onSuccess(rows -> {
                Row row = rows.iterator().next();
                if (row.getLong("qty") > 0) {
                    JsonObject response = JsonObject.mapFrom(new ErrorResponse(409, "Conflict", "用户已存在", context.request().uri()));
                    context.response().setStatusCode(409).end(response.encode());
                    return;
                }
                Long id = IdUtil.getSnowflake(1, 1).nextId();
                String salt = IdUtil.simpleUUID();
                String detail = "{}";
                try {
                    JsonObject d = new JsonObject()
                        .put("salt", salt)
                        .put("hash", hashPassword(body.getString("password"), salt, "hex"));
                    detail = d.encode();
                } catch (Exception e) {
                    System.out.println(e.getMessage());
                }
                JsonObject state = new JsonObject().put("uuid", IdUtil.randomUUID()).put("created_at", DateUtil.formatDateTime(new Date()));
                pool.preparedQuery("""
                        insert into crate.subscriber (id, email, name, phone, tags, detail, time, state)
                        values ($1, $2, $3, $4, $5, $6, $7, $8)
                        """)
                    .execute(Tuple.of(id, body.getString("email"), body.getString("email"), "", "[]", detail, new Date(), state.encode()))
                    .onSuccess(result -> context.response().setStatusCode(201).end())
                    .onFailure(err -> {
                        JsonObject response = JsonObject.mapFrom(new ErrorResponse(500, "Internal Server Error", err.getMessage(), context.request().uri()));
                        context.response().setStatusCode(500).end(response.encode());
                    });
            })
            .onFailure(err -> {
                JsonObject response = JsonObject.mapFrom(new ErrorResponse(500, "Internal Server Error", err.getMessage(), context.request().uri()));
                context.response().setStatusCode(500).end(response.encode());
            });
         */
    }

    private void logIn(RoutingContext context) {
        JsonObject body = new JsonObject(context.body().asString());
        String username = body.getString("username");
        pool.preparedQuery("""
                select * from crate.subscriber
                where name = $1 or email = $2 or phone = $3
                """)
            .execute(Tuple.of(username, username, username))
            .onSuccess(rows -> {
                if (rows.size() == 0) {
                    JsonObject response = JsonObject.mapFrom(ErrorResponse.builder()
                        .status(404)
                        .title("Not Found")
                        .detail("用户不存在")
                        .instance(context.request().uri())
                        .build());
                    context.response().setStatusCode(404).end(response.encode());
                    return;
                } else if (rows.size() > 1) {
                    JsonObject response = JsonObject.mapFrom(ErrorResponse.builder()
                        .status(500)
                        .title("Internal Server Error")
                        .detail("服务器错误")
                        .instance(context.request().uri())
                        .build());
                    context.response().setStatusCode(500).end(response.encode());
                    return;
                }
                Row row = rows.iterator().next();
                JsonObject detailJson = new JsonObject(row.getValue("detail").toString());
                try {
                    String hashedPassword = hashPassword(body.getString("password"), detailJson.getString("salt"), "hex");
                    if (detailJson.getString("hash").equals(hashedPassword)) {
                        String token = generateToken(row.getValue("id").toString(), row.getValue("email").toString(), row.getValue("name").toString());
                        JsonObject response = new JsonObject().put("token", token);
                        context.response().end(response.encode());
                    } else {
                        JsonObject response = JsonObject.mapFrom(ErrorResponse.builder()
                            .status(401)
                            .title("Unauthorized")
                            .detail("密码错误")
                            .instance(context.request().uri())
                            .build());
                        context.response().setStatusCode(401).end(response.encode());
                    }
                } catch (Exception e) {
                    JsonObject response = JsonObject.mapFrom(ErrorResponse.builder()
                        .status(500)
                        .title("Internal Server Error")
                        .detail(e.getMessage())
                        .instance(context.request().uri())
                        .build());
                    context.response().setStatusCode(500).end(response.encode());
                }
            })
            .onFailure(err -> {
                JsonObject response = JsonObject.mapFrom(ErrorResponse.builder()
                    .type("about:blank")
                    .status(500)
                    .title("Internal Server Error")
                    .detail(err.getMessage())
                    .instance(context.request().uri())
                    .build());
                context.response().setStatusCode(500).end(response.encode());
            });
    }

    private void refresh(RoutingContext context) {
        //
    }
}
