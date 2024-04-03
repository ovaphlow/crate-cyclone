import * as constants from "../utilities/constants.js";
import * as repository from "./repository.js";

/**
 * @param {import("koa").Context} ctx
 */
export const get = async (ctx) => {
  const { id } = ctx.params;
  if (id) {
    ctx.response.body = { id };
    return;
  }
  const { option } = ctx.request.query || "";
  if (option === "default") {
    const result = await repository.defaultFilter(
      {
        skip:
          (parseInt(ctx.request.query["page"]?.toString() || "1") - 1) *
          parseInt(ctx.request.query["take"]?.toString() || "10"),
        take: parseInt(ctx.request.query["take"]?.toString() || "10"),
      },
      {
        equal: (ctx.request.query["equal"]?.toString() || "[]").split(","),
        objectContain: (
          ctx.request.query["object-contain"]?.toString() || "[]"
        ).split(","),
        arrayContain: (
          ctx.request.query["array-contain"]?.toString() || "[]"
        ).split(","),
        like: (ctx.request.query["like"]?.toString() || "[]").split(","),
        objectLike: (
          ctx.request.query["object-like"]?.toString() || "[]"
        ).split(","),
        inList: (ctx.request.query["in"]?.toString() || "[]").split(","),
        lesser: (ctx.request.query["lesser"]?.toString() || "[]").split(","),
        greater: (ctx.request.query["greater"]?.toString() || "[]").split(","),
      },
    );
    ctx.set(constants.HEADER_API_VERSION, "2024-01-06");
    ctx.response.body = result;
    return;
  }
  ctx.response.status = 406;
  return;
};
