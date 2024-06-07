import { createAsyncThunk } from '@reduxjs/toolkit';
import { type Session, SessionStates } from '@shared/types/sessions';
import { getSessions } from '@api/sessions/getSessions';
import { deleteSession } from '@api/webDriver/deleteSession';
import { createSession } from '@api/webDriver/createSession';
import { updateSessionState } from './sessionsSlice';

interface CreateSessionArgs {
  browser: string;
  version: string;
  image: string;
}

interface DeleteSessionArgs {
  sessionId: string;
}

enum TYPE_THUNK_PREFIXES {
  CREATE_SESSION = 'wd/createSession',
  DELETE_SESSION = 'wd/deleteSession',
  FETCH_SESSIONS = 'sessions/fetch',
}

export const createWdSession = createAsyncThunk<Promise<void>, CreateSessionArgs>(
  TYPE_THUNK_PREFIXES.CREATE_SESSION,
  async (desiredBrowser: CreateSessionArgs): Promise<void> => {
    const { browser, version } = desiredBrowser;
    await createSession(browser, version);
  }
);

export const deleteWdSession = createAsyncThunk<Promise<void>, DeleteSessionArgs>(
  TYPE_THUNK_PREFIXES.DELETE_SESSION,
  async ({ sessionId }: DeleteSessionArgs, { dispatch }): Promise<void> => {
    dispatch(updateSessionState({ id: sessionId, newState: SessionStates.TERMINATING }));
    await deleteSession(sessionId);
  }
);

export const fetchSessions = createAsyncThunk<Session[]>(TYPE_THUNK_PREFIXES.FETCH_SESSIONS, async () => {
  return await getSessions();
});
