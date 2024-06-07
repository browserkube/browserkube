export interface Session {
  id: string;
  name: string;
  image: string;
  type: string;
  state: string;
  browser: string;
  browserVersion: string;
  manual: boolean;
  vncOn: boolean;
  vncPsw: string;
  logsOn: boolean;
  createdAt: number;
  platformName: string;
  screenResolution: string;
  sessionType: string;
}

export interface Sessions {
  byId: SessionsByIdType;
}

export interface SessionManual {
  sessionName: string;
  platformName: string;
  browserName: string;
  browserVersion: string;
  screenResolution: string;
  recordVideo: boolean;
}

export interface SessionLinesArray {
  sessionArr: SessionLine[];
}

export enum SessionStates {
  PENDING = 'pending',
  RUNNING = 'running',
  TERMINATING = 'terminating',
  TERMINATED = 'terminated',
}

type ActiveSessionGridRowsBaseType = Record<
  | 'id'
  | 'name'
  | 'state'
  | 'image'
  | 'type'
  | 'availability'
  | 'browser'
  | 'platformName'
  | 'screenResolution'
  | 'browserVersion',
  string
>;

export type ActiveSessionGridRowsType = ActiveSessionGridRowsBaseType & {
  manual: boolean;
};

export type SessionLine = ActiveSessionGridRowsType | TerminatedSession;

export type SessionsByIdType = Record<string, Session>;

export type TerminatedSession = Session;

export type TerminatedSessionsByIdType = Record<string, TerminatedSession>;

export interface TerminatedSessions {
  byId: TerminatedSessionsByIdType;
}

export interface SessionDetails {
  browser: string;
  browserVersion: string;
  createdAt: number;
  id: string;
  image: string;
  logsRefAddr: string;
  videoRefAddr: string;
  state: string;
  type: string;
  vncOn: boolean;
  logsOn: boolean;
}

export interface SessionDetailsLogs {
  text: string;
  shouldRefetch: boolean;
}

export interface Commands {
  sessionId: string;
  commandId: string;
  command: string;
  method: string;
  request: null;
  statusCode: number;
  response: string;
  timestamp: string;
}

export interface SessionDetailsCommands {
  data: Commands[];
  shouldRefetch: boolean;
  newPageToken: string;
}

export interface CommandsParams {
  pageToken: string;
  pageSize: number;
}

export interface CommandsResponse {
  commands: Commands[];
  newPageToken: string;
}

export interface SessionDetailsType {
  activeSessionId: string;
  logs: SessionDetailsLogs;
  commands: SessionDetailsCommands;
  sessionData: SessionDetails;
}

export interface TerminatedSessionsResponse {
  Items: TerminatedSession[];
}
