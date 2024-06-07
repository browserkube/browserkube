import { createContext, useState, useContext, type Dispatch, type SetStateAction, useEffect, useRef } from 'react';
import { useSelector } from 'react-redux';
import { type TerminatedSession, type SessionLine } from '@shared/types/sessions';
import { useAppDispatch } from '@hooks/useAppDispatch';
import { getActiveSessionId } from '@redux/sessionDetails/selectors';
import { saveActiveSessionId } from '@redux/sessionDetails/sessionSlice';
import { fetchSessionDetails } from '@redux/sessionDetails/thunk';
import { getSession } from '@redux/sessions/sessionsSelectors';

interface SessionContextValue {
  session: TerminatedSession;
  setSession: Dispatch<SetStateAction<SessionLine>>;
}

const SessionContext = createContext<SessionContextValue | null>(null);

export const ChosenSessionProvider = (props: any) => {
  const dispatch = useAppDispatch();
  const [session, setSession] = useState<TerminatedSession | null>(null);
  const activeSessionId = useSelector(getActiveSessionId);
  const activeSessionById = useSelector((s) => getSession(s, activeSessionId));
  const prevSession = useRef('');

  const ACTIVE_STATE_TERMINATED = session && session.state === 'Terminated';

  useEffect(() => {
    if (activeSessionById) {
      setSession(activeSessionById);
    }
  }, [activeSessionById]);

  useEffect(() => {
    if (!session) {
      return;
    }
    const newSessionId = session.id;
    if (activeSessionId !== newSessionId) {
      dispatch(
        saveActiveSessionId({
          id: newSessionId,
        })
      );
      if (ACTIVE_STATE_TERMINATED) {
        void dispatch(fetchSessionDetails(newSessionId));
      }
      prevSession.current = newSessionId;
    }
  }, [dispatch, session?.id]);

  return <SessionContext.Provider value={{ session, setSession }} {...props} />;
};

export const useSession = () => {
  const context = useContext(SessionContext);
  if (!context) {
    throw new Error('useSession must be used within a ChosenSessionProvider');
  }
  return context;
};
