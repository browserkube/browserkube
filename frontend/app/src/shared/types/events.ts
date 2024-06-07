import { type Session } from '@shared/types/sessions';
import { type SessionStatus } from '@shared/types/sessionStatus';

export enum WebSocketStatus {
  DISCONNECTED = 'disconnected',
  ESTABLISHING = 'establishing',
  CONNECTED = 'connected',
}

export interface EventsPayload<T = unknown> {
  name: string;
  payload: T;
}

export type SessionMessage = EventsPayload<Session[]>;

export type StatusMessage = EventsPayload<SessionStatus>;

export function isSessionMessage(e: EventsPayload<unknown>): e is EventsPayload<Session[]> {
  return 'name' in e && 'payload' in e && e.name === 'session';
}

export function isStatusMessage(e: EventsPayload<unknown>): e is EventsPayload<SessionStatus> {
  return 'name' in e && 'payload' in e && e.name === 'status';
}
