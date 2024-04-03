import mysql from "mysql2";

import {
  DB_HOST,
  DB_USERNAME,
  DB_PASSWORD,
  DB_NAME,
  PROC,
} from "../configuration.js";

export const pool = mysql.createPool({
  host: DB_HOST,
  user: DB_USERNAME,
  password: DB_PASSWORD,
  database: DB_NAME,
  waitForConnections: true,
  connectionLimit: PROC * 2 + 1,
  queueLimit: 0,
});
