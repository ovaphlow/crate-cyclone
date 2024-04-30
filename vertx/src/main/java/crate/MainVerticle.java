package crate;

import crate.application.handler.SchemaHandler;
import crate.schema.repository.SchemaRepositoryImpl;
import crate.schema.service.SchemaService;
import crate.utility.Config;
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
        Router router = Router.router(vertx);
        router.route().handler(BodyHandler.create().setBodyLimit(1024 * 1024 * 10));

        Config config = new Config(".env");
        int port = Integer.parseInt(config.get("PORT"));
        PgConnectOptions connectOptions = config.getPgConnectOptions();
        PoolOptions poolOptions = new PoolOptions().setMaxWaitQueueSize(5);
        Pool pool = Pool.pool(vertx, connectOptions, poolOptions);

        new HandlerSetting(pool).setupRoutes(router);
        new HandlerSubscriber(vertx, pool).setupRoutes(router);
        new HandlerSchema(pool).setupRoutes(router);

        SchemaRepositoryImpl schemaRepository = new SchemaRepositoryImpl(pool);
        SchemaService schemaService = new SchemaService(schemaRepository);

        SchemaHandler schemaHandler = new SchemaHandler(schemaService);
        schemaHandler.setupRoutes(router);

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
