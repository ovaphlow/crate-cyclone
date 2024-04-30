import { Hono } from "hono";
import { SchemaService } from "./service";

export const createSchemaHandler = (service: SchemaService) => {
    const setupRoutes = (app: Hono) => {
        app.get("/crate-api/db-schema", async (c) => {
            const schemas = await service.listSchemas();
            return c.json(schemas);
        });

        app.get("/crate-api/:schema/db-table", async (c) => {
            const schema = c.req.param("schema");
            const tables = await service.listTables(schema || "");
            return c.json(tables);
        });

        app.get("/crate-api/:schema/:table/db-column", async (c) => {
            const schema = c.req.param("schema");
            const table = c.req.param("table");
            const columns = await service.listColumns(schema || "", table || "");
            return c.json(columns);
        });
    };

    return {
        setupRoutes,
    };
};
