import { type Session, SessionStates } from '@shared/types/sessions';
import {
  type SessionMessage,
  type StatusMessage,
  isSessionMessage,
  isSessionSnapshotMessage,
  isStatusMessage,
} from '@shared/types/events';
import { saveStats } from '@redux/sessionStatus/sessionStatusSlice';
import { addTerminatedSession } from '@redux/terminatedSessions/slice';
import { fetchTerminatedSessions } from '@redux/terminatedSessions/thunk';
import { type AppDispatch, type AppState } from '../store';
import { addSession, removeSession, reconcileSessions } from '../sessions/sessionsSlice';

export const EventMessageHandlers = (dispatch: AppDispatch, getState: () => AppState) => {
  const handleStatus = (message: StatusMessage) => {
    dispatch(saveStats(message.payload));
  };

  const handleSession = (message: SessionMessage) => {
    const { payload } = message;

    payload.forEach((session: Session) => {
      const { id, state } = session;
      const normalizedState = state.toLowerCase();
      if (normalizedState === SessionStates.TERMINATED) {
        dispatch(removeSession({ id }));
        dispatch(addTerminatedSession({ session }));
        void dispatch(fetchTerminatedSessions());
      } else {
        dispatch(addSession({ session }));
      }
    });
  };

  const handleSnapshot = (message: SessionMessage) => {
    dispatch(reconcileSessions(message.payload));
  };

  const handleMessage = (message: MessageEvent<string>) => {
    const parsedMessage = JSON.parse(message.data);
    if (isSessionSnapshotMessage(parsedMessage)) {
      handleSnapshot(parsedMessage);
    } else if (isSessionMessage(parsedMessage)) {
      handleSession(parsedMessage);
    }
    if (isStatusMessage(parsedMessage)) {
      handleStatus(parsedMessage);
    }
  };
  return { handleMessage };
};
