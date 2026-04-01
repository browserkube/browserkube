import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import {
  type Session,
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
    addTerminatedSession(state, action: PayloadAction<{ session: Session }>): TerminatedSessionsState {
      const { session } = action.payload;
      state.data.byId[session.id] = {
        ...session,
        state: session.state.toLowerCase(),
        name: session.name || session.id,
      };
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
          resultMap[item.id] = { ...item, state: item.state.toLowerCase(), name: item.name || item.id };
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
export const { addTerminatedSession } = actions;
