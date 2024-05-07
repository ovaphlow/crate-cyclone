package crate.infrastructure;

import io.vertx.core.Vertx;
import io.vertx.core.json.JsonObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.time.Instant;

public class EventPublisher {

    private final Logger logger = LoggerFactory.getLogger(EventPublisher.class);
    private boolean eventLog = false;

    public EventPublisher() {
        if ("true".equals(new Configuration(".env").get("EVENT_LOG"))) {
            eventLog = true;
        }
    }

    public void setupEventConsumer(Vertx vertx) {
        vertx.eventBus().consumer("endpoint", message -> {
            if (eventLog) {
                JsonObject data = (JsonObject) message.body();
                logger.info("Received message: {}", data.encodePrettily());
                logger.info("time {}", Instant.now().toString());
            }
        });
    }

    public void publishEvent(Vertx vertx, JsonObject message) {
        vertx.eventBus().publish("endpoint", message);
    }
}
