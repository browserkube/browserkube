import { type AppState } from '@redux/store';
import { type Session, type Sessions } from '@shared/types/sessions';
import { REDUCER_STATUS } from '@shared/types/reducerType';

export const getSessions = (state: AppState): Sessions => {
  return state.sessions.data;
};

export const getSession = (state: AppState, id: string): Session => {
  return state.sessions.data.byId[id];
};

export const isSessionsFetched = (state: AppState): boolean => {
  return state.sessions.fetchSessionsStatus === REDUCER_STATUS.FULFILLED;
};

export const isNewSessionCreating = (state: AppState): boolean => {
  return state.sessions.createSessionStatus === REDUCER_STATUS.PENDING;
};
