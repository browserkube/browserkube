export interface SessionStats {
  all: number;
  running: number;
  connecting?: number;
  queued?: number;
}

export interface SessionStatus {
  quotesLimit: number;
  maxTimeout: number;
  stats: SessionStats;
}
