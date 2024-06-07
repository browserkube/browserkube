export const API_URL = {
  STATUS: '/status',
  BROWSERS: '/browsers',
  SESSIONS: '/sessions/',
  SESSION: '/wd/hub/session',
  TERMINATED_SESSIONS: '/results/',
  SESSION_FILE: (sessionId: string) => `/sessions/${sessionId}/files/`,
  SCREENSHOTS: '/screenshots',
  TEST: '/api/browsers', // TODO: temporary, while the [BE] will disable /api for the post method
};
