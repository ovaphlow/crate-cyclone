package crate.setting;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

public class SettingService {
    private final SettingRepository repo;

    public SettingService(SettingRepository repo) {
        this.repo = repo;
    }

    public List<Map<String, Object>> defaultList(Map<String, Object> options) {
        List<Setting> settings = repo.retrieve(options);
        List<Map<String, Object>> result = new ArrayList<>();
        for (Setting setting : settings) {
            Map<String, Object> it = Map.of(
                "id", setting.getId(),
                "rootId", setting.getRootId(),
                "parentId", setting.getParentId(),
                "tags", setting.getTags(),
                "detail", setting.getDetail(),
                "state", setting.getState(),
                "_id", setting.getId().toString(),
                "_rootId", setting.getRootId().toString(),
                "_parentId", setting.getParentId().toString()
            );
            result.add(it);
        }
        return result;
    }
}
