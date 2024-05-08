package crate.setting;

import io.vertx.sqlclient.Pool;
import io.vertx.sqlclient.Tuple;

import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.stream.Collectors;
import java.util.stream.StreamSupport;

public class SettingRepository {
    private final Pool pool;

    public SettingRepository(Pool pool) {
        this.pool = pool;
    }

}
