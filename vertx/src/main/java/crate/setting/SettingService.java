package crate.setting;

import io.vertx.core.Future;

import java.util.List;

public class SettingService {
    private final SettingRepository repo;

    public SettingService(SettingRepository repo) {
        this.repo = repo;
    }

    public Future<List<Setting>> listSettings() {
        return repo.retrieve();
    }
}
