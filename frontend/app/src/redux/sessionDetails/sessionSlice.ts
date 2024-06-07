import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import { type CommandsResponse, type SessionDetails, type SessionDetailsType } from '@shared/types/sessions';
import { REDUCER_STATUS, type ReducerType } from '@shared/types/reducerType';
import { fetchSessionDetails, fetchSessionDetailsCommands, fetchSessionDetailsLogs } from './thunk';

type SessionDetailsState = ReducerType<SessionDetailsType>;

const initialState: SessionDetailsState = {
  data: {
    activeSessionId: '',
    logs: {
      text: '',
      shouldRefetch: false,
    },
    commands: {
      data: [
        {
          sessionId: '',
          commandId: '',
          command: '',
          method: '',
          request: null,
          statusCode: 0,
          response: '',
          timestamp: '',
        },
      ],
      shouldRefetch: false,
      newPageToken: 'first page',
    },
    sessionData: {
      browser: '',
      browserVersion: '',
      logsRefAddr: '',
      createdAt: new Date().getTime(),
      id: '',
      image: '',
      type: '',
      state: '',
      videoRefAddr: '',
      vncOn: false,
      logsOn: false,
    },
  },
  status: REDUCER_STATUS.IDLE,
  error: null,
};

const sessionDetails = createSlice({
  name: 'sessionDetails',
  initialState,
  reducers: {
    saveActiveSessionId(state, action: PayloadAction<{ id: string }>): ReducerType<SessionDetailsType> {
      state.data.activeSessionId = action.payload.id;
      state.data.sessionData = initialState.data.sessionData; // extra re-writing data, that forces re-render
      state.data.logs = initialState.data.logs;
      state.data.commands = initialState.data.commands;
      return state;
    },
  },
  extraReducers(builder) {
    builder
      .addCase(fetchSessionDetails.fulfilled, (state, action: PayloadAction<SessionDetails>) => {
        state.data.sessionData = action.payload;
        state.data.logs.shouldRefetch = true;
        state.status = REDUCER_STATUS.FULFILLED;
      })
      .addCase(fetchSessionDetails.pending, (state) => {
        state.status = REDUCER_STATUS.PENDING;
        state.error = null;
      })
      .addCase(fetchSessionDetails.rejected, (state) => {
        state.status = REDUCER_STATUS.REJECTED;
      })
      .addCase(fetchSessionDetailsLogs.pending, (state) => {
        state.status = REDUCER_STATUS.PENDING;
        state.error = null;
      })
      .addCase(fetchSessionDetailsLogs.fulfilled, (state, action: PayloadAction<string>) => {
        state.data.logs = {
          text: action.payload,
          shouldRefetch: false,
        };
        state.status = REDUCER_STATUS.FULFILLED;
      })
      .addCase(fetchSessionDetailsLogs.rejected, (state) => {
        state.data.logs = {
          text: initialState.data.logs.text,
          shouldRefetch: false,
        };
        state.status = REDUCER_STATUS.REJECTED;
      })
      .addCase(fetchSessionDetailsCommands.pending, (state) => {
        state.status = REDUCER_STATUS.PENDING;
        state.error = null;
      })
      .addCase(fetchSessionDetailsCommands.fulfilled, (state, action: PayloadAction<CommandsResponse>) => {
        const { commands, newPageToken } = action.payload;
        let mergedCommands = [];
        if (state.data.commands.data.length === 1 && !state.data.commands.data[0].sessionId) {
          mergedCommands = commands;
        } else {
          mergedCommands = [...state.data.commands.data, ...commands];
        }
        state.data.commands = {
          data: mergedCommands,
          newPageToken,
          shouldRefetch: false,
        };
        state.status = REDUCER_STATUS.FULFILLED;
      })
      .addCase(fetchSessionDetailsCommands.rejected, (state) => {
        state.data.commands = {
          data: initialState.data.commands.data,
          shouldRefetch: false,
          newPageToken: 'first page',
        };
        state.status = REDUCER_STATUS.REJECTED;
      });
  },
});

export const { reducer: sessionDetailsReducer, actions } = sessionDetails;
export const { saveActiveSessionId } = actions;
