import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import {
  type TerminatedSessionsByIdType,
  type TerminatedSessionsResponse,
  type TerminatedSessions,
} from '@shared/types/sessions';
import { REDUCER_STATUS, type ReducerType } from '@shared/types/reducerType';
import { fetchTerminatedSessions } from './thunk';

type TerminatedSessionsState = ReducerType<TerminatedSessions>;

const initialState: TerminatedSessionsState = {
  data: {
    byId: {},
  },
  status: REDUCER_STATUS.IDLE,
  error: null,
};

const terminatedSessionsSlice = createSlice({
  name: 'terminatedSessions',
  initialState,
  reducers: {
    // TODO: delete after implementation a new real reducer
    stubReducer(state): TerminatedSessionsState {
      return state;
    },
  },
  extraReducers(builder) {
    builder
      .addCase(fetchTerminatedSessions.pending, (state) => {
        state.status = REDUCER_STATUS.PENDING;
        state.error = null;
      })
      .addCase(fetchTerminatedSessions.fulfilled, (state, action: PayloadAction<TerminatedSessionsResponse>) => {
        const sessions = action.payload.Items;
        state.data.byId = sessions.reduce<TerminatedSessionsByIdType>(function (resultMap, item) {
          resultMap[item.id] = { ...item, name: item.name || item.id };
          return resultMap;
        }, {});
        state.status = REDUCER_STATUS.FULFILLED;
      })
      .addCase(fetchTerminatedSessions.rejected, (state) => {
        state.status = REDUCER_STATUS.REJECTED;
      });
  },
});

export const { reducer: terminatedSessionsReducer, actions } = terminatedSessionsSlice;
export const { stubReducer } = actions;
