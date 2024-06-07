import { type AppState } from '@redux/store';
import { type VncSessionType } from '@shared/types/VncSession';

export const getVncSession = (state: AppState): VncSessionType => {
  return state.VncSession;
};

export const getVncSessionId = (state: AppState): string => {
  return state.VncSession.sessionId;
};

export const getVncSessionScreenStatus = (state: AppState): boolean => {
  return state.VncSession.locked;
};

export const getVncSessionLogStatus = (state: AppState): boolean => {
  return state.VncSession.swipeableOpened;
};
