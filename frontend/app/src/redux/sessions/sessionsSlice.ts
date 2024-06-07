import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import { type Session, type Sessions, type SessionsByIdType } from '@shared/types/sessions';
import { REDUCER_STATUS, type ReducerType } from '@shared/types/reducerType';
import { createWdSession, fetchSessions } from './sessionsThunk';

interface SessionsStatusData {
  fetchSessionsStatus: REDUCER_STATUS;
  error: string | null;
  createSessionStatus: REDUCER_STATUS;
}

type SessionsState = ReducerType<Sessions, SessionsStatusData>;

const initialState: SessionsState = {
  data: {
    byId: {},
  },
  fetchSessionsStatus: REDUCER_STATUS.IDLE,
  createSessionStatus: REDUCER_STATUS.IDLE,
  error: null,
};

const sessionsSlice = createSlice({
  name: 'sessions',
  initialState,
  reducers: {
    addSession(state, action: PayloadAction<{ session: Session }>): SessionsState {
      const { session } = action.payload;
      state.data.byId[session.id] = {
        ...session,
        state: session.state.toLowerCase(),
        name: session.name || session.id,
      };
      // TODO: remove this line when create session flow will be fixed on BE
      state.createSessionStatus = REDUCER_STATUS.FULFILLED;
      return state;
    },
    removeSession(state, action: PayloadAction<{ id: string }>): SessionsState {
      const { payload } = action;
      Reflect.deleteProperty(state.data.byId, payload.id);
      return state;
    },
    updateSessionState(state, action: PayloadAction<{ id: string; newState: string }>): SessionsState {
      const { id, newState } = action.payload;
      const session = state.data.byId[id];
      state.data.byId[id] = { ...session, state: newState.toLowerCase() };
      return state;
    },
  },
  extraReducers(builder) {
    builder
      .addCase(fetchSessions.pending, (state) => {
        state.fetchSessionsStatus = REDUCER_STATUS.PENDING;
        state.error = null;
      })
      .addCase(fetchSessions.fulfilled, (state, action: PayloadAction<Session[]>) => {
        const sessions = action.payload;
        state.data.byId = sessions.reduce<SessionsByIdType>((resultMap, item) => {
          resultMap[item.id] = { ...item, state: item.state.toLowerCase(), name: item.name || item.id };
          return resultMap;
        }, {});
        state.fetchSessionsStatus = REDUCER_STATUS.FULFILLED;
        return state;
      })
      .addCase(createWdSession.pending, (state) => {
        state.createSessionStatus = REDUCER_STATUS.PENDING;
        state.error = null;
      })
      .addCase(createWdSession.fulfilled, (state) => {
        state.createSessionStatus = REDUCER_STATUS.FULFILLED;
      })
      .addCase(createWdSession.rejected, (state) => {
        state.createSessionStatus = REDUCER_STATUS.REJECTED;
      });
  },
});

export const { reducer: sessionsReducer, actions } = sessionsSlice;
export const { addSession, removeSession, updateSessionState } = actions;
