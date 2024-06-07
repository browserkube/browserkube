import { api, API_URL } from '@shared/api';
import { type SessionDetails } from '@shared/types/sessions';

export const getBrowserKubeSessionDetails = async (id: string): Promise<SessionDetails> => {
  return await api.get<SessionDetails>({ url: API_URL.TERMINATED_SESSIONS + id });
};
