export const arrayToObject = (array) => {
  const obj = {};
  for (let i = 0; i < array.length; i += 2) {
    obj[array[i]] = array[i + 1];
  }
  return obj;
};
