import { api, API_URL } from '@shared/api';
import { type Session, type SessionManual } from '@shared/types/sessions';
import { type CreateSessionResponse } from '@shared/types/webDriver';

export const getSessions = async (): Promise<Session[]> => {
  return await api.get<Session[]>({ url: API_URL.SESSIONS });
};

export const CreateManualSession = async (body: SessionManual) => {
  return await api.post<CreateSessionResponse>({
    url: API_URL.TEST,
    data: body,
    params: { timeout: Number(process.env.REACT_APP_CREAT_SESSION_TIMEOUT) },
  });
};
// return await api.post<CreateSessionResponse>({ url: API_URL.TEST, data: body, timeout: Number(process.env.REACT_APP_CREAT_SESSION_TIMEOUT) });
