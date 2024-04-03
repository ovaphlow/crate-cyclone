import Koa from "koa";
import bodyParser from "koa-bodyparser";
import helmet from "koa-helmet";
import koaLogger from "koa-logger";
import Router from "@koa/router";
import { logger } from "./utilities/logger-pino.js";

export const app = new Koa();

app.use(helmet());

app.use(bodyParser({ jsonLimit: "16mb" }));

app.use(koaLogger((str, args) => logger.debug(str, args)));

app.use(async (ctx, next) => {
  try {
    await next();
  } catch (error) {
    if (error instanceof Error) {
      ctx.response.status = 500;
      ctx.response.body = { message: error.message };
    } else {
      ctx.response.status = 500;
      ctx.response.body = { message: "服务器错误" };
    }
  }
});

app.on("error", (err, ctx) => {
  logger.error(`${ctx.req.method} ${ctx.req.url}`);
  logger.error(err.stack);
});

const router = new Router();
const prefix = "/crate-api";

(() => {
  router.get(
    `${prefix}/`,
    /** @param {import("koa").Context} ctx */
    (ctx) => {
      ctx.body = "Hello, world!";
    },
  );
})();

(() => {
  import("./events/endpoint.js").then(({ get }) => {
    router.get(`/crate-api/event`, get);
  });
})();

app.use(router.routes());
app.use(router.allowedMethods());
