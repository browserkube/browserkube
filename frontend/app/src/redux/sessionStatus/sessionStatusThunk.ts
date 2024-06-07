import { createAsyncThunk } from '@reduxjs/toolkit';
import { type SessionStatus } from '@shared/types/sessionStatus';
import { getSessionStatus } from '@api/sessionStatus/getSessionStatus';

export const fetchSessionStatus = createAsyncThunk<SessionStatus>('sessionStatus', async () => {
  return await getSessionStatus();
});
