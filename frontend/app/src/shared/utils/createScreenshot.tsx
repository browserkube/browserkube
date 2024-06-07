import { API_URL, api } from '@shared/api';

export const handleScreenshot = async (activeSessionId: string) => {
  if (activeSessionId) {
    const URL = `${API_URL.SESSIONS}${activeSessionId}${API_URL.SCREENSHOTS}`;
    try {
      await api.post({ url: URL, data: {} });
    } catch (e) {
      console.log('Error, while making the screenshot', e);
    }
  }
};
