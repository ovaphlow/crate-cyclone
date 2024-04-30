import pg from "pg";
import dotenv from "dotenv";

const { Pool } = pg;

dotenv.config();

export class SchemaRepository {
    private pool: Pool;

    constructor() {
        this.pool = new Pool({
            user: process.env.PG_USER,
            password: process.env.PG_PASSWORD,
            host: process.env.PG_HOST,
            port: parseInt(process.env.PG_PORT || "5432", 10),
            database: process.env.PG_DATABASE,
            max: 3,
            idleTimeoutMillis: 20000,
            connectionTimeoutMillis: 3000,
        });
    }

    async listSchemas() {
        const client = await this.pool.connect();
        try {
            const result = await client.query("select schema_name from information_schema.schemata");
            return result.rows.map((row: { schema_name: string }) => row.schema_name);
        } finally {
            client.release();
        }
    }

    async listTables(schema: string) {
        const client = await this.pool.connect();
        try {
            const result = await client.query(
                "select table_name from information_schema.tables where table_schema = $1",
                [schema],
            );
            return result.rows.map((row: { table_name: string }) => row.table_name);
        } finally {
            client.release();
        }
    }

    async listColumns(schema: string, table: string) {
        const client = await this.pool.connect();
        try {
            const result = await client.query(
                `
                select column_name, data_type from information_schema.columns
                where table_schema = $1 and table_name = $2
                `,
                [schema, table],
            );
            return result.rows.map((row: { column_name: string; data_type: string }) => {
                return {
                    columnName: row.column_name,
                    dataType: row.data_type,
                };
            });
        } finally {
            client.release();
        }
    }
}
