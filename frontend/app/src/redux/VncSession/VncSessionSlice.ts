import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import { type VncSessionType } from '@shared/types/VncSession';

interface VncSessionPayload {
  sessionId: string;
}

interface UpdateVncSessionScreenStatusPayload {
  VNCScreenStatus: boolean;
}

interface UpdateVncSessionLogStatusPayload {
  swipeableOpened: boolean;
}

const initialState: VncSessionType = {
  sessionId: '',
  locked: true,
  swipeableOpened: false,
  swipeableAnchor: 'bottom',
  isInitialized: false,
};

const sessionsSlice = createSlice({
  name: 'VncSession',
  initialState,
  reducers: {
    initVncSession(state): VncSessionType {
      state.isInitialized = true;
      return state;
    },
    createVncSession(state, action: PayloadAction<VncSessionPayload>): VncSessionType {
      const { sessionId } = action.payload;
      return {
        sessionId,
        locked: true,
        swipeableOpened: false,
        swipeableAnchor: 'bottom',
        isInitialized: true,
      };
    },
    deleteVncSession(state): VncSessionType {
      state.sessionId = '';
      state.isInitialized = false;
      return state;
    },
    updateVncSessionScreenStatus(state, action: PayloadAction<UpdateVncSessionScreenStatusPayload>): VncSessionType {
      const { VNCScreenStatus } = action.payload;
      state.locked = VNCScreenStatus;
      return state;
    },
    updateVncSessionLogStatus(state, action: PayloadAction<UpdateVncSessionLogStatusPayload>): VncSessionType {
      const { swipeableOpened } = action.payload;
      state.swipeableOpened = swipeableOpened;
      return state;
    },
  },
});

export const { reducer: VncSessionReducer, actions } = sessionsSlice;
export const {
  initVncSession,
  createVncSession,
  deleteVncSession,
  updateVncSessionScreenStatus,
  updateVncSessionLogStatus,
} = actions;
