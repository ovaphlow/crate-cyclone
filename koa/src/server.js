import cluster from "cluster";
import http from "http";

import { logger } from "./utilities/logger-pino.js";

let port = 8000;

import("dotenv").then((dotenv) => {
  dotenv.config();
  port = parseInt(process.env.PORT?.toString() || "8000", 10);
});

if (cluster.isPrimary) {
  logger.info(`主进程 PID:${process.pid}`);

  for (let i = 0; i < parseInt(process.env.PROC || "1", 10); i += 1) {
    cluster.fork();
  }

  cluster.on("online", (worker) => {
    logger.info(`子进程 PID:${worker.process.pid}, 端口:${port}`);
  });

  cluster.on("exit", (worker, code, signal) => {
    logger.error(
      `子进程 PID:${worker.process.pid}终止，错误代码:${code}，信号:${signal}`,
    );
    logger.info(`由主进程(PID:${process.pid})创建新的子进程`);
    cluster.fork();
  });
} else {
  import("./app.js").then(({ app }) => {
    http.createServer(app.callback()).listen(port);
  });
}
