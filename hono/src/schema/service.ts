import { SchemaRepository } from "./repository";
import { v4 as uuidv4 } from "uuid";
import { Snowflake} from "../infrastructure/snowflake";

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

    async save(schema: string, table: string, data: Record<string, unknown>) {
        const snowflake = new Snowflake(1);
        data["id"] = snowflake.nextId();
        const state = {
            "uuid": uuidv4(),
            "created_at": new Date().toISOString(),
        }
        data["state"] = JSON.stringify(state);
        return await this.repo.save(schema, table, data);
    }
}
