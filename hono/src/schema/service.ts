import { SchemaRepository } from "./repository";

export class SchemaService {
    private repo: SchemaRepository;

    constructor(repo: SchemaRepository) {
        this.repo = repo;
    }

    async listSchemas() {
        return await this.repo.listSchemas();
    }

    async listTables(schema: string) {
        return await this.repo.listTables(schema);
    }

    async listColumns(schema: string, table: string) {
        return await this.repo.listColumns(schema, table);
    }
}
