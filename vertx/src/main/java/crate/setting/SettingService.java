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
}
