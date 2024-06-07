import { api, API_URL } from '@shared/api';
import { type Browser } from '@shared/types/browsers';

export const getBrowsers = async () => {
  return await api.get<Browser[]>({ url: API_URL.BROWSERS });
};
