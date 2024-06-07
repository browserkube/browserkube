import { type AppState } from '@redux/store';
import { type SessionStats } from '@shared/types/sessionStatus';
import { type REDUCER_STATUS } from '@shared/types/reducerType';

// TODO: Hack, need to check how BE send max timeout
export const getMaxTimeout = (state: AppState): number => {
  return (state.state.data.maxTimeout / 1000000000) * 60000;
};

export const getStats = (state: AppState): SessionStats => {
  return state.state.data.stats;
};

export const getQuotesLimit = (state: AppState): number => {
  return state.state.data.quotesLimit;
};

export const getStatus = (state: AppState): REDUCER_STATUS => {
  return state.state.status;
};
