import { configureStore, combineReducers } from '@reduxjs/toolkit';
import { sessionStatusReducer } from '@redux/sessionStatus/sessionStatusSlice';
import { webSocketReducer } from './webSocket/webSocketSlice';
import { terminatedSessionsReducer } from './terminatedSessions/slice';
import { UIReducer } from './UI/UISlice';
import { sessionDetailsReducer } from './sessionDetails/sessionSlice';
import { sessionsReducer } from './sessions/sessionsSlice';
import { VncSessionReducer } from './VncSession/VncSessionSlice';
import { eventWebsocketMiddleware } from './middleware/eventWebsocketMiddleware';
import { browsersReducer } from './browsers/browsersSlice';
import { chipsReducer } from './chips/chipsSlice';

export const appStore = configureStore({
  reducer: combineReducers({
    browsers: browsersReducer,
    chips: chipsReducer,
    state: sessionStatusReducer,
    sessions: sessionsReducer,
    webSocket: webSocketReducer,
    VncSession: VncSessionReducer,
    terminatedSessions: terminatedSessionsReducer,
    sessionDetails: sessionDetailsReducer,
    UI: UIReducer,
  }),
  middleware: (getDefaultMiddleware) => {
    return getDefaultMiddleware().concat(eventWebsocketMiddleware());
  },
});

export type AppState = ReturnType<typeof appStore.getState>;
export type AppDispatch = typeof appStore.dispatch;
