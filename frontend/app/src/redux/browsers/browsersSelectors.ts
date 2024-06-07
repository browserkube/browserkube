import { type AppState } from '@redux/store';
import { type Browser } from '@shared/types/browsers';

export const getBrowsers = (state: AppState): Browser[] => state.browsers.data.browsers;
