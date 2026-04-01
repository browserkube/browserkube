import { type Middleware } from 'redux';
import { type AppDispatch, type AppState } from '@redux/store';
import { connect, disconnect, connected, catchError } from '../webSocket/webSocketSlice';
import { EventMessageHandlers } from '../webSocket/handlers';

interface AppStateContext {
  dispatch: AppDispatch;
  getState: () => AppState;
}

const MAX_RECONNECT_DELAY_MS = 30_000;
const BASE_RECONNECT_DELAY_MS = 1_000;

export const eventWebsocketMiddleware = (): Middleware => {
  let socket: WebSocket | null = null;
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  let reconnectAttempt = 0;
  let intentionalClose = false;

  const clearReconnectTimer = () => {
    if (reconnectTimer !== null) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }
  };

  const createSocket = (url: string, dispatch: AppDispatch, getState: () => AppState) => {
    const Events = EventMessageHandlers(dispatch, getState);

    if (socket !== null) {
      intentionalClose = true;
      socket.close();
    }
    intentionalClose = false;
    socket = new WebSocket(url);

    socket.onopen = () => {
      reconnectAttempt = 0;
      dispatch(connected());
    };

    socket.onmessage = (message) => {
      Events.handleMessage(message);
    };

    socket.onclose = () => {
      dispatch(disconnect());
      if (!intentionalClose) {
        const delay = Math.min(BASE_RECONNECT_DELAY_MS * 2 ** reconnectAttempt, MAX_RECONNECT_DELAY_MS);
        reconnectAttempt++;
        clearReconnectTimer();
        reconnectTimer = setTimeout(() => {
          createSocket(url, dispatch, getState);
        }, delay);
      }
    };

    socket.onerror = (error) => {
      dispatch(catchError({ error }));
    };
  };

  return ({ dispatch, getState }: AppStateContext) => {
    return (next) => (action) => {
      if (connect.match(action)) {
        clearReconnectTimer();
        reconnectAttempt = 0;
        createSocket(action.payload.url, dispatch, getState);
      }
      if (disconnect.match(action) && socket !== null) {
        clearReconnectTimer();
        intentionalClose = true;
        socket.close();
        socket = null;
      }
      return next(action);
    };
  };
};
