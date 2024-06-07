import { useSelector } from 'react-redux';
import { useEffect } from 'react';
import { useAppDispatch } from '@hooks/useAppDispatch';
import { isSessionsFetched, getSession } from '@redux/sessions/sessionsSelectors';
import { getVncSession } from '@redux/VncSession/VncSessionSelectors';
import { initVncSession, createVncSession, deleteVncSession } from '@redux/VncSession/VncSessionSlice';
import { fetchSessions } from '@redux/sessions/sessionsThunk';

export const useVncPanel = (id: string) => {
  const dispatch = useAppDispatch();
  const sessionsFetched = useSelector(isSessionsFetched);
  const session = useSelector((state) => getSession(state, id));
  const vncSession = useSelector((state) => getVncSession(state));

  useEffect(() => {
    if (!sessionsFetched) {
      void dispatch(fetchSessions());
    }
  }, [dispatch, sessionsFetched]);

  useEffect(() => {
    if (sessionsFetched && session) {
      dispatch(createVncSession({ sessionId: id }));
    } else {
      dispatch(initVncSession());
    }
    return () => {
      dispatch(deleteVncSession());
    };
  }, [session, id, sessionsFetched, dispatch]);

  return vncSession;
};
