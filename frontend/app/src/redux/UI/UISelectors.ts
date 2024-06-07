import { type AppState } from '@redux/store';
import { type Modal } from '@shared/types/UI';

export const getSideBarStatus = (state: AppState): boolean => state.UI.isSideBarOpen;

export const getModals = (state: AppState): Modal[] => state.UI.modals;
