import { type Session, SessionStates } from '@shared/types/sessions';
import { type SessionMessage, type StatusMessage, isSessionMessage, isStatusMessage } from '@shared/types/events';
import { saveStats } from '@redux/sessionStatus/sessionStatusSlice';
import { type AppDispatch, type AppState } from '../store';
import { addSession, removeSession, updateSessionState } from '../sessions/sessionsSlice';

export const EventMessageHandlers = (dispatch: AppDispatch, getState: () => AppState) => {
  const handleStatus = (message: StatusMessage) => {
    dispatch(saveStats(message.payload));
  };
  const handleSession = (message: SessionMessage) => {
    const { payload } = message;
    const sessions = getState().sessions.data;

    payload.forEach((session: Session) => {
      const { id, state } = session;
      if (!sessions.byId[id]) {
        dispatch(addSession({ session }));
      } else if (state === SessionStates.TERMINATED) {
        dispatch(removeSession({ id }));
      } else if (state !== sessions.byId[id].state) {
        dispatch(updateSessionState({ id, newState: state }));
      }
    });
  };

  const handleMessage = (message: MessageEvent<string>) => {
    const parsedMessage = JSON.parse(message.data);
    if (isSessionMessage(parsedMessage)) {
      handleSession(parsedMessage);
    }
    if (isStatusMessage(parsedMessage)) {
      handleStatus(parsedMessage);
    }
  };
  return { handleMessage };
};
