package crate.utility;

import io.vertx.pgclient.PgConnectOptions;

import java.io.FileInputStream;
import java.io.IOException;
import java.util.Properties;

public class Config {

    private final Properties properties;

    public Config(String envFilePath) {
        properties = new Properties();
        try {
            FileInputStream fis = new FileInputStream(envFilePath);
            properties.load(fis);
        } catch (IOException e) {
            System.out.println(e.getMessage());
        }
    }

    public String get(String key) {
        return properties.getProperty(key);
    }

    public PgConnectOptions getPgConnectOptions() {
        return new PgConnectOptions()
            .setUser(get("DB_USER"))
            .setPassword(get("DB_PASSWORD"))
            .setHost(get("DB_HOST"))
            .setPort(Integer.parseInt(get("DB_PORT")))
            .setDatabase(get("DB_NAME"));
    }
}
