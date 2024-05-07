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

    public List<Setting> retrieve(Map<String, Object> options) {
        CompletableFuture<List<Setting>> future = new CompletableFuture<>();
        int take = (int) options.getOrDefault("take", 20);
        long offset = (long) options.getOrDefault("page", 1) * take - take;
        pool.preparedQuery("select * from crate.setting order by id desc limit $1 offset $2")
            .execute(Tuple.of(take, offset))
            .onSuccess(rows -> future.complete(StreamSupport.stream(rows.spliterator(), false)
                .map(row -> new Setting(
                    row.getLong("id"),
                    row.getLong("root_id"),
                    row.getLong("parent_id"),
                    row.getString("tags"),
                    row.getString("detail"),
                    row.getString("state")))
                .collect(Collectors.toList())))
            .onFailure(future::completeExceptionally);
        try {
            return future.get();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }
}
