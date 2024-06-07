export const getMethodCommand = (method: string) => {
  if (!method) {
    return '';
  }

  return method.split('/')[1];
};
