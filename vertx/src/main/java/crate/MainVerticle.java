package crate;

import crate.application.handler.SchemaHandler;
import crate.application.handler.SettingHandler;
import crate.infrastructure.Configuration;
import crate.infrastructure.EventPublisher;
import crate.infrastructure.Middleware;
import crate.schema.SchemaRepository;
import crate.schema.SchemaService;
import crate.setting.SettingRepository;
import crate.setting.SettingService;
import io.vertx.core.AbstractVerticle;
import io.vertx.core.Promise;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.handler.BodyHandler;
import io.vertx.pgclient.PgConnectOptions;
import io.vertx.sqlclient.Pool;
import io.vertx.sqlclient.PoolOptions;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class MainVerticle extends AbstractVerticle {

    private static final Logger logger = LoggerFactory.getLogger(MainVerticle.class);

    @Override
    public void start(Promise<Void> startPromise) throws Exception {
        Configuration config = new Configuration(".env");
        int port = Integer.parseInt(config.get("PORT"));
        PgConnectOptions connectOptions = config.getPgConnectOptions();
        PoolOptions poolOptions = new PoolOptions().setMaxWaitQueueSize(5);
        Pool pool = Pool.pool(vertx, connectOptions, poolOptions);

        EventPublisher event = new EventPublisher();
        event.setupEventConsumer(vertx);

        Router router = Router.router(vertx);
        router.route().handler(BodyHandler.create().setBodyLimit(1024 * 1024 * 10));
        router.route().handler(Middleware::logRequestHandler);

        SchemaRepository schemaRepository = new SchemaRepository(pool);
        SchemaService schemaService = new SchemaService(schemaRepository);
        SchemaHandler schemaHandler = new SchemaHandler(schemaService);
        schemaHandler.setupRoutes(router);

        SettingRepository settingRepository = new SettingRepository(pool);
        SettingService settingService = new SettingService(settingRepository);
        SettingHandler settingHandler = new SettingHandler(settingService);
        settingHandler.setupRoutes(router);

        vertx.createHttpServer()
            .requestHandler(router)
            .listen(port, http -> {
                if (http.succeeded()) {
                    startPromise.complete();
                    logger.info("HTTP server started on port {}", port);
                } else {
                    startPromise.fail(http.cause());
                }
            });
    }
}
