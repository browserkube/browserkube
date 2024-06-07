import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import { type SessionStatus } from '@shared/types/sessionStatus';
import { REDUCER_STATUS, type ReducerType } from '@shared/types/reducerType';
import { fetchSessionStatus } from './sessionStatusThunk';

type SessionStatusState = ReducerType<SessionStatus>;

const initialState: SessionStatusState = {
  status: REDUCER_STATUS.IDLE,
  error: null,
  data: {
    quotesLimit: 0,
    maxTimeout: 0,
    stats: {
      all: 0,
      running: 0,
      connecting: 0,
      queued: 0,
    },
  },
};

const sessionStatusSlice = createSlice({
  name: 'sessionStatus',
  initialState,
  reducers: {
    saveStats(state, action: PayloadAction<SessionStatus>) {
      state.data = action.payload;
    },
  },
  extraReducers(builder) {
    builder
      .addCase(fetchSessionStatus.pending, (state) => {
        state.status = REDUCER_STATUS.PENDING;
        state.error = null;
      })
      .addCase(fetchSessionStatus.fulfilled, (state, action: PayloadAction<SessionStatus>) => {
        state.status = REDUCER_STATUS.FULFILLED;
        state.error = null;
        state.data = action.payload;
      })
      .addCase(fetchSessionStatus.rejected, (state) => {
        state.status = REDUCER_STATUS.REJECTED;
      });
  },
});

export const { reducer: sessionStatusReducer, actions } = sessionStatusSlice;
export const { saveStats } = actions;
