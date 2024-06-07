import { createAsyncThunk } from '@reduxjs/toolkit';
import { getBrowsers } from '@api/browsers/getBrowsers';

export const fetchBrowsers = createAsyncThunk('fetchBrowsers', async () => {
  return await getBrowsers();
});
