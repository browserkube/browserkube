import { api, API_URL } from '@shared/api';
import { type SessionStatus } from '@shared/types/sessionStatus';

export const getSessionStatus = async (): Promise<SessionStatus> => {
  return await api.get<SessionStatus>({ url: API_URL.STATUS });
};
