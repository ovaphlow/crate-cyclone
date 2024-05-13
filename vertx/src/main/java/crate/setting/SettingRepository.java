package crate.setting;

import io.vertx.core.Future;
import io.vertx.core.Promise;
import io.vertx.sqlclient.Pool;

import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.StreamSupport;

public class SettingRepository {
    private final Pool pool;

    public SettingRepository(Pool pool) {
        this.pool = pool;
    }

    public Future<List<Setting>> retrieve() {
        Promise<List<Setting>> promise = Promise.promise();
        pool.query("select * from crate.setting order by id desc")
            .execute()
            .onSuccess(rows -> promise.complete(StreamSupport.stream(rows.spliterator(), false)
                .map(row -> new Setting.Builder()
                    .id(row.getLong("id"))
                    .rootId(row.getLong("root_id"))
                    .parentId(row.getLong("parent_id"))
                    .tags(row.getJsonArray("tags").toString())
                    .detail(row.getJsonObject("detail").toString())
                    .createdAt(row.getOffsetDateTime("created_at"))
                    .updatedAt(row.getOffsetDateTime("updated_at"))
                    .state(row.getJsonObject("state").toString())
                    .build())
                .collect(Collectors.toList())))
            .onFailure(promise::fail);
        return promise.future();
    }
}
