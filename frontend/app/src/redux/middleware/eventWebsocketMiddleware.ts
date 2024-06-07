import { type Middleware } from 'redux';
import { type AppDispatch, type AppState } from '@redux/store';
import { connect, disconnect, connected, catchError } from '../webSocket/webSocketSlice';
import { EventMessageHandlers } from '../webSocket/handlers';

interface AppStateContext {
  dispatch: AppDispatch;
  getState: () => AppState;
}

export const eventWebsocketMiddleware = (): Middleware => {
  let socket: WebSocket | null = null;

  return ({ dispatch, getState }: AppStateContext) => {
    return (next) => (action) => {
      const Events = EventMessageHandlers(dispatch, getState);
      if (connect.match(action)) {
        if (socket !== null) {
          socket.close();
        }
        socket = new WebSocket(action.payload.url);

        socket.onopen = () => {
          dispatch(connected());
        };

        socket.onmessage = (message) => {
          Events.handleMessage(message);
        };

        socket.onclose = () => {
          dispatch(disconnect());
        };

        socket.onerror = (error) => {
          dispatch(catchError({ error }));
        };
      }
      return next(action);
    };
  };
};
