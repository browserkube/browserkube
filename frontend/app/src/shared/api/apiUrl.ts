export const API_URL = {
  STATUS: '/status',
  BROWSERS: '/browsers',
  SESSIONS: '/sessions/',
  CREATE_SESSION: '/sessions',
  SESSION: '/wd/hub/session',
  TERMINATED_SESSIONS: '/results/',
  SESSION_FILE: (sessionId: string) => `/sessions/${sessionId}/files/`,
  SCREENSHOTS: '/screenshots',
};
