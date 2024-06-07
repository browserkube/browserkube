import { createAsyncThunk } from '@reduxjs/toolkit';
import { getBrowserKubeSessionDetails } from '@api/sessionDetails/getBrowserKubeSessionDetails';
import { type CommandsResponse, type CommandsParams, type SessionDetails } from '@shared/types/sessions';
import { getSessionCommands, getSessionLogs } from '@api/sessionDetails/getSessionLogs';

export const fetchSessionDetails = createAsyncThunk<SessionDetails, string>('sessionDetails/fetch', async (id) => {
  return await getBrowserKubeSessionDetails(id);
});

export const fetchSessionDetailsLogs = createAsyncThunk<string, string>('sessionDetailsLogs/fetch', async (url) => {
  return await getSessionLogs(url);
});

export const fetchSessionDetailsCommands = createAsyncThunk<CommandsResponse, { url: string; params: CommandsParams }>(
  'sessionDetailsCommands/fetch',
  async ({ url, params }: { url: string; params: CommandsParams }) => {
    return await getSessionCommands(url, params);
  }
);
