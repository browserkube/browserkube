import { api, API_URL } from '@shared/api';

export const deleteSession = async (id: string): Promise<void> => {
  await api.remove({
    url: `${API_URL.SESSION}/${id}`,
    params: { timeout: Number(process.env.REACT_APP_DELETE_SESSION_TIMEOUT) },
  });
};
