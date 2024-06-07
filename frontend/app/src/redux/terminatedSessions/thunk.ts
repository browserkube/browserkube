import { createAsyncThunk } from '@reduxjs/toolkit';
import { type TerminatedSessionsResponse } from '@shared/types/sessions';
import { getGogaTerminatedSessions } from '@api/terminatedSessions/getGogaTerminatedSessions';

const TYPE_THUNK_PREFIX = 'terminatedSessions/fetch';
export const fetchTerminatedSessions = createAsyncThunk<TerminatedSessionsResponse>(TYPE_THUNK_PREFIX, async () => {
  return await getGogaTerminatedSessions();
});
