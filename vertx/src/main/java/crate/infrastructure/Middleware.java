package crate.infrastructure;

import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.RoutingContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Middleware {

    private static final Logger logger = LoggerFactory.getLogger(Middleware.class);

    public static void logRequestHandler(RoutingContext context) {
        logger.info("{} {} {}",
            context.request().remoteAddress().host(),
            context.request().method().name(),
            context.request().uri());

        JsonObject message = new JsonObject()
            .put("method", context.request().method().name())
            .put("endpoint", context.request().uri())
            .put("body", context.body().asString())
            .put("client_ip", context.request().remoteAddress().host());
        new EventPublisher().publishEvent(context.vertx(), message);

        context.next();
    }
}
