import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import { WebSocketStatus } from '@shared/types/events';

interface WebSocketState {
  status: WebSocketStatus;
  url: string;
  error: Event | null;
}

const initialState: WebSocketState = {
  status: WebSocketStatus.DISCONNECTED,
  url: '',
  error: null,
};

const webSocketSlice = createSlice({
  name: 'WebSocket',
  initialState,
  reducers: {
    connect(state, action: PayloadAction<{ url: string }>) {
      state.url = action.payload.url;
      state.status = WebSocketStatus.ESTABLISHING;
      return state;
    },
    disconnect(state) {
      state.status = WebSocketStatus.DISCONNECTED;
      return state;
    },
    connected(state) {
      state.status = WebSocketStatus.CONNECTED;
      return state;
    },
    catchError(state, action: PayloadAction<{ error: Event }>) {
      state.error = action.payload.error;
      return state;
    },
  },
});

export const { reducer: webSocketReducer, actions } = webSocketSlice;
export const { connect, disconnect, connected, catchError } = actions;
