import {
  equalBuilder,
  objectContainBuilder,
  arrayContainBuilder,
  likeBuilder,
  objectLikeBuilder,
  inBuilder,
  lesserBuilder,
  greaterBuilder,
} from "./condition-builder.js";
import { pool } from "./database-mysql.js";

const schemas = [
  {
    name: "event",
    table: "events",
    columns: ["id", "relation_id", "reference_id", "tags", "detail", "time"],
  },
  {
    name: "setting",
    table: "settings",
    columns: ["id", "root_id", "parent_id", "tags", "time"],
  },
];

/**
 * @param {string} name
 * @param {number} skip
 * @param {number} take
 * @param {object} option
 * @param {string[]} option.equal
 * @param {string[]} option.objectContain
 * @param {string[]} option.arrayContain
 * @param {string[]} option.like
 * @param {string[]} option.objectLike
 * @param {string[]} option.inList
 * @param {string[]} option.lesser
 * @param {string[]} option.greater
 */
export async function retrieve(name, skip, take, option) {
  const schema = schemas.find((schema) => schema.name === name);
  if (!schema) {
    throw new Error(`schema ${name} not found`);
  }
  let q = `select ${schema.columns.join(", ")} from ${schema.table}`;
  /** @type {string[]} */
  const conditions = [];
  /** @type {string[]} */
  const params = [];
  /** @type {string[]} */
  const equal = option["equal"] || [];
  if (equal.length > 0 && equal.length % 2 === 0) {
    equalBuilder(equal, conditions, params);
  }
  const objectContain = option["objectContain"] || [];
  if (objectContain.length > 0 && objectContain.length % 3 === 0) {
    objectContainBuilder(objectContain, conditions, params);
  }
  const arrayContain = option["arrayContain"] || [];
  if (arrayContain.length > 0 && arrayContain.length % 2 === 0) {
    arrayContainBuilder(arrayContain, conditions, params);
  }
  const like = option["like"] || [];
  if (like.length > 0 && like.length % 2 === 0) {
    likeBuilder(like, conditions, params);
  }
  const objectLike = option["objectLike"] || [];
  if (objectLike.length > 0 && objectLike.length % 3 === 0) {
    objectLikeBuilder(objectLike, conditions, params);
  }
  const inList = option["inList"] || [];
  if (inList.length > 1) {
    inBuilder(inList, conditions, params);
  }
  const lesser = option["lesser"] || [];
  if (lesser.length > 0 && lesser.length % 2 === 0) {
    lesserBuilder(lesser, conditions, params);
  }
  const greater = option["greater"] || [];
  if (greater.length > 0 && greater.length % 2 === 0) {
    greaterBuilder(greater, conditions, params);
  }
  if (conditions.length > 0) {
    q = `${q} where ${conditions.join(" and ")}`;
  }
  q = `${q} order by id desc limit ${skip}, ${take}`;
  console.info(q);
  const client = pool.promise();
  const [result] = await client.execute(q, params);
  return result;
}
