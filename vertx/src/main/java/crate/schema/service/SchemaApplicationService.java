package crate.schema.service;

import crate.schema.repository.SchemaRepositoryImpl;

public class SchemaApplicationService {

    private final SchemaRepositoryImpl repo;

    public SchemaApplicationService(SchemaRepositoryImpl repo) {
        this.repo = repo;
    }
}
