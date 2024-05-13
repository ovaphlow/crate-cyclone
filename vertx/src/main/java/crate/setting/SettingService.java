package crate.setting;

import io.vertx.core.Future;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

public class SettingService {
    private final SettingRepository repo;

    public SettingService(SettingRepository repo) {
        this.repo = repo;
    }

    public Future<List<Setting>> listSettings() {
        return repo.retrieve();
    }
}
