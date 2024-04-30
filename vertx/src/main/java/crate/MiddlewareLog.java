package crate;

import io.vertx.ext.web.RoutingContext;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class MiddlewareLog {

    private static final Logger logger = LoggerFactory.getLogger(MiddlewareLog.class);

    public static void logRequestHandler(RoutingContext context) {
        String method = context.request().method().toString();
        String uri = context.request().uri();

        logger.info("{} {}", method, uri);

        context.next();
    }
}
