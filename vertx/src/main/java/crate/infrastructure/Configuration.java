package crate.infrastructure;

import io.vertx.pgclient.PgConnectOptions;

import java.io.FileInputStream;
import java.io.IOException;
import java.util.Properties;

public class Configuration {

    private final Properties properties;

    public Configuration(String envFilePath) {
        properties = new Properties();
        try {
            FileInputStream fis = new FileInputStream(envFilePath);
            properties.load(fis);
        } catch (IOException e) {
            System.out.println(e.getMessage());
        }
    }

    public String get(String key) {
        String value = System.getenv(key);
        if (null == value) {
            value = properties.getProperty(key);
        }
        return properties.getProperty(key);
    }

    public PgConnectOptions getPgConnectOptions() {
        return new PgConnectOptions()
            .setUser(get("PGSQL_USER"))
            .setPassword(get("PGSQL_PASSWORD"))
            .setHost(get("PGSQL_HOST"))
            .setPort(Integer.parseInt(get("PGSQL_PORT")))
            .setDatabase(get("PGSQL_DATABASE"));
    }
}
