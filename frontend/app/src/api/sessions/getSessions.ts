import { api, API_URL } from '@shared/api';
import { type Session, type SessionManual } from '@shared/types/sessions';

export const getSessions = async (): Promise<Session[]> => {
  return await api.get<Session[]>({ url: API_URL.SESSIONS });
};

export const CreateManualSession = async (body: SessionManual): Promise<Session> => {
  return await api.post<Session>({
    url: API_URL.CREATE_SESSION,
    data: body,
  });
};
