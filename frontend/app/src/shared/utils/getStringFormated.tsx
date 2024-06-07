/**
 *
 * @param str
 * @returns the same string with the first capital letter, others toLowercase()
 */

export const getStringFormated = (str: string): string => {
  if (!str) {
    return '';
  }
  const res = `${str.charAt(0).toUpperCase()}${str.slice(1).toLowerCase()}`;
  return res;
};
