import * as conditionBuilder from "../utilities/condition-builder.js";
import { pool } from "../utilities/database-mysql.js";

const columns = ["id", "relation_id", "reference_id", "tags", "detail", "time"];

/**
 * @param {object} option
 * @param {number} option.skip
 * @param {number} option.take
 * @param {object} filter
 * @param {string[]} filter.equal
 * @param {string[]} filter.objectContain
 * @param {string[]} filter.arrayContain
 * @param {string[]} filter.like
 * @param {string[]} filter.objectLike
 * @param {string[]} filter.inList
 * @param {string[]} filter.lesser
 * @param {string[]} filter.greater
 */
export async function defaultFilter(
  { skip, take },
  {
    equal,
    objectContain,
    arrayContain,
    like,
    objectLike,
    inList,
    lesser,
    greater,
  },
) {
  let q = `select ${columns.join(",")} from events`;
  /** @type {string[]} */
  const conditions = [];
  /** @type {string[]} */
  const params = [];
  if (equal.length > 0) {
    const { c, p } = conditionBuilder.equal(equal);
    conditions.push(...c);
    params.push(...p);
  }
  if (objectContain.length) {
    const { c, p } = conditionBuilder.objectContain(objectContain);
    conditions.push(...c);
    params.push(...p);
  }
  if (arrayContain.length) {
    const { c, p } = conditionBuilder.arrayContain(arrayContain);
    conditions.push(...c);
    params.push(...p);
  }
  if (like.length) {
    const { c, p } = conditionBuilder.like(like);
    conditions.push(...c);
    params.push(...p);
  }
  if (objectLike.length) {
    const { c, p } = conditionBuilder.objectLike(objectLike);
    conditions.push(...c);
    params.push(...p);
  }
  if (inList.length) {
    const { c, p } = conditionBuilder.inList(inList);
    conditions.push(...c);
    params.push(...p);
  }
  if (lesser.length) {
    const { c, p } = conditionBuilder.lesser(lesser);
    conditions.push(...c);
    params.push(...p);
  }
  if (greater.length) {
    const { c, p } = conditionBuilder.greater(greater);
    conditions.push(...c);
    params.push(...p);
  }
  if (conditions.length > 0) {
    q = `${q} where ${conditions.join(" and ")}`;
  }
  q = `${q} order by id desc limit ${skip}, ${take}`;
  console.info(q);
  const client = pool.promise();
  const [result] = await client.execute(q, params);
  const rows = [];
  /** @ts-ignore */
  for (const it of result) {
    rows.push({
      ...it,
      tags: JSON.stringify(it.tags),
      detail: JSON.stringify(it.detail),
      _id: it.id?.toString(),
      _relationId: it.relation_id?.toString(),
      _referenceId: it.reference_id?.toString(),
    });
  }
  return rows;
}
