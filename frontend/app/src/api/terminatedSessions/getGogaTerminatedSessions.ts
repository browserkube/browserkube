import { api, API_URL } from '@shared/api';
import { type TerminatedSessionsResponse } from '@shared/types/sessions';

export const getGogaTerminatedSessions = async (): Promise<TerminatedSessionsResponse> => {
  return await api.get<TerminatedSessionsResponse>({ url: API_URL.TERMINATED_SESSIONS });
};
