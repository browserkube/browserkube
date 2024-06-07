export const getSessionId = () => {
  const hashArray = window.location.hash.split('/');
  return hashArray[hashArray.length - 1];
};
