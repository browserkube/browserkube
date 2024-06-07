import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import { type Modal, type UIData } from '@shared/types/UI';

const START_MODAL_ZINDEX = 1000;

const initialState: UIData = {
  isSideBarOpen: false,
  modals: [],
};

const UISlice = createSlice({
  name: 'UIState',
  initialState,
  reducers: {
    openSideBar(state) {
      return { ...state, isSideBarOpen: true };
    },
    closeSideBar(state) {
      return { ...state, isSideBarOpen: false };
    },
    openModal(state, action: PayloadAction<Pick<Modal, 'id' | 'component'>>) {
      const { modals } = state;
      const { id, component } = action.payload;
      const lastModalZIndex = modals.length ? modals[modals.length - 1].zIndex : START_MODAL_ZINDEX;
      return { ...state, modals: [...modals, { id, component, zIndex: lastModalZIndex + 1 }] };
    },
    closeModal(state, action: PayloadAction<string>) {
      const { modals } = state;
      return { ...state, modals: modals.filter(({ id }) => id !== action.payload) };
    },
    closeAllModals(state) {
      return { ...state, modals: [] };
    },
  },
});

export const { reducer: UIReducer, actions } = UISlice;
export const { openSideBar, closeSideBar, openModal, closeModal, closeAllModals } = actions;
