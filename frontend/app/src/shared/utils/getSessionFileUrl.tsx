import { API_URL } from '@shared/api';

export const getSessionFileUrl = (sessionId: string, fileName: string): string => {
  const url: string = API_URL.SESSION_FILE(sessionId);

  return `${url}${fileName}`;
};
