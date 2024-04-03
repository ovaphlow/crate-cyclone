/**
 * @param {string[]} equal
 */
export function equal(equal) {
  if (equal.length % 2 !== 0) {
    console.debug("equal length is not even");
    return { c: [], p: [] };
  }
  const c = [];
  const p = [];
  for (let i = 0; i < equal.length; i += 2) {
    c.push(`${equal[i]} = ?`);
    p.push(equal[i + 1]);
  }
  return { c, p };
}

/**
 * @param {string[]} equal
 * @param {string[]} conditions
 * @param {string[]} params
 */
export function equalBuilder(equal, conditions, params) {
  for (let i = 0; i < equal.length; i += 2) {
    conditions.push(`${equal[i]} = ?`);
    params.push(equal[i + 1]);
  }
}

/**
 * @param {string[]} objectContain
 */
export function objectContain(objectContain) {
  if (objectContain.length % 3 !== 0) {
    console.debug("objectContain length is not multiple of 3");
    return { c: [], p: [] };
  }
  const c = [];
  const p = [];
  for (let i = 0; i < objectContain.length; i += 3) {
    c.push(
      `json_contains(${objectContain[i]}, json_object('${
        objectContain[i + 1]
      }', ?))`,
    );
    p.push(objectContain[i + 2]);
  }
  return { c, p };
}

/**
 * @param {string[]} objectContain
 * @param {string[]} conditions
 * @param {string[]} params
 */
export function objectContainBuilder(objectContain, conditions, params) {
  for (let i = 0; i < objectContain.length; i += 3) {
    conditions.push(
      `json_contains(${objectContain[i]}, json_object('${
        objectContain[i + 1]
      }', ?))`,
    );
    params.push(objectContain[i + 2]);
  }
}

/**
 * @param {string[]} arrayContain
 * @returns
 */
export function arrayContain(arrayContain) {
  if (arrayContain.length % 2 !== 0) {
    console.debug("arrayContain length is not even");
    return { c: [], p: [] };
  }
  const c = [];
  const p = [];
  for (let i = 0; i < arrayContain.length; i += 2) {
    c.push(`json_contains(${arrayContain[i]}, json_array(?))`);
    p.push(arrayContain[i + 1]);
  }
  return { c, p };
}

/**
 * @param {string[]} arrayContain
 * @param {string[]} conditions
 * @param {string[]} params
 */
export function arrayContainBuilder(arrayContain, conditions, params) {
  for (let i = 0; i < arrayContain.length; i += 2) {
    conditions.push(`json_contains(${arrayContain[i]}, json_array(?))`);
    params.push(arrayContain[i + 1]);
  }
}

/**
 * @param {string[]} like
 */
export function like(like) {
  if (like.length % 2 !== 0) {
    console.debug("like length is not even");
    return { c: [], p: [] };
  }
  const c = [];
  const p = [];
  for (let i = 0; i < like.length; i += 2) {
    c.push(`position(? in ${like[i]})`);
    p.push(like[i + 1]);
  }
  return { c, p };
}

/**
 * @param {string[]} like
 * @param {string[]} conditions
 * @param {string[]} params
 */
export function likeBuilder(like, conditions, params) {
  for (let i = 0; i < like.length; i += 2) {
    conditions.push(`position(? in ${like[i]})`);
    params.push(like[i + 1]);
  }
}

/**
 * @param {string[]} objectLike
 */
export function objectLike(objectLike) {
  if (objectLike.length % 3 !== 0) {
    console.debug("objectLike length is not multiple of 3");
    return { c: [], p: [] };
  }
  const c = [];
  const p = [];
  for (let i = 0; i < objectLike.length; i += 3) {
    c.push(`position(? in ${objectLike[i]}->>'$.${objectLike[i + 1]}')`);
    p.push(objectLike[i + 2]);
  }
  return { c, p };
}

/**
 * @param {string[]} objectLike
 * @param {string[]} conditions
 * @param {string[]} params
 */
export function objectLikeBuilder(objectLike, conditions, params) {
  for (let i = 0; i < objectLike.length; i += 3) {
    conditions.push(
      `position(? in ${objectLike[i]}->>'$.${objectLike[i + 1]}')`,
    );
    params.push(objectLike[i + 2]);
  }
}

/**
 * @param {string[]} inList
 */
export function inList(inList) {
  if (inList.length < 2) {
    console.debug("inList length is less than 2");
    return { c: [], p: [] };
  }
  const c_ = new Array(inList.length - 1).fill("?");
  const c = [`${inList[0]} in (${c_.join(",")})`];
  const p = inList.slice(1);
  return { c, p };
}

/**
 * @param {string[]} inList
 * @param {string[]} conditions
 * @param {string[]} params
 */
export function inBuilder(inList, conditions, params) {
  const c = new Array(inList.length - 1).fill("?");
  conditions.push(`${inList[0]} in (${c.join(",")})`);
  params.push(...inList.slice(1));
}

/**
 * @param {string[]} lesser
 */
export function lesser(lesser) {
  if (lesser.length % 2 !== 0) {
    console.debug("lesser length is not even");
    return { c: [], p: [] };
  }
  const c = [];
  const p = [];
  for (let i = 0; i < lesser.length; i += 2) {
    c.push(`${lesser[i]} <= ?`);
    p.push(lesser[i + 1]);
  }
  return { c, p };
}

/**
 * @param {string[]} lesser
 * @param {string[]} conditions
 * @param {string[]} params
 */
export function lesserBuilder(lesser, conditions, params) {
  for (let i = 0; i < lesser.length; i += 2) {
    conditions.push(`${lesser[i]} <= ?`);
    params.push(lesser[i + 1]);
  }
}

/**
 * @param {string[]} greater
 */
export function greater(greater) {
  if (greater.length % 2 !== 0) {
    console.debug("greater length is not even");
    return { c: [], p: [] };
  }
  const c = [];
  const p = [];
  for (let i = 0; i < greater.length; i += 2) {
    c.push(`${greater[i]} >= ?`);
    p.push(greater[i + 1]);
  }
  return { c, p };
}

/**
 * @param {string[]} greater
 * @param {string[]} conditions
 * @param {string[]} params
 */
export function greaterBuilder(greater, conditions, params) {
  for (let i = 0; i < greater.length; i += 2) {
    conditions.push(`${greater[i]} >= ?`);
    params.push(greater[i + 1]);
  }
}
