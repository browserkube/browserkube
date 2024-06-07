import { api, API_URL } from '@shared/api';
import { createCapabilities } from './capabilityService';

export const createSession = async (browser: string, version: string) => {
  // TODO: update void generic
  return await api.post<unknown>({
    url: API_URL.SESSION,
    data: createCapabilities(browser, version),
    params: { timeout: Number(process.env.REACT_APP_CREAT_SESSION_TIMEOUT) },
  });
};
