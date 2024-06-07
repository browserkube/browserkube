import { type AppState } from '@redux/store';
import { type TerminatedSession, type TerminatedSessions } from '@shared/types/sessions';

export const getTerminatedSessionsState = (state: AppState): TerminatedSessions => {
  return state.terminatedSessions;
};

export const getTerminatedSessions = (state: AppState): TerminatedSession[] => {
  return Object.values(state.terminatedSessions.data.byId);
};
