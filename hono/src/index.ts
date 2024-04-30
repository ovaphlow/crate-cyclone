import { serve } from "@hono/node-server";
import { Hono } from "hono";
import { logger } from "hono/logger";
import dotenv from "dotenv";
import { SchemaRepository } from "./schema/repository";
import { SchemaService } from "./schema/service";
import { createSchemaHandler } from "./schema/handler";

dotenv.config();

const app = new Hono();

app.use(logger());

const schemaRepository = new SchemaRepository();
const schemaService = new SchemaService(schemaRepository);

(() => {
    const handler = createSchemaHandler(schemaService);
    handler.setupRoutes(app);
})();

const port = parseInt(process.env.PORT || "8421", 10);
console.log(`Server is running on port ${port}`);

serve({
    fetch: app.fetch,
    port,
});
