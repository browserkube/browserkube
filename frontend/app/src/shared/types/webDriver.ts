export interface CreateSessionResponse {
  value: {
    sessionId: string;
    capabilities: {
      acceptInsecureCerts?: boolean;
      browserName: string;
      browserVersion: string;
      chrome?: {
        chromedriverVersion: string;
        userDataDir: string;
      };
      'browserkube:chromeOptions'?: { debuggerAddress: string };
      networkConnectionEnabled?: boolean;
      pageLoadStrategy?: string;
      platformName: string;
      // TODO update proxy typing
      proxy?: unknown;
      setWindowRect?: boolean;
      strictFileInteractivity?: boolean;
      timeouts: { implicit: 0; pageLoad: number; script: number };
      unhandledPromptBehavior?: string;
      'webauthn:extension:credBlob'?: boolean;
      'webauthn:extension:largeBlob'?: boolean;
      'webauthn:virtualAuthenticators'?: boolean;
    };
  };
}
