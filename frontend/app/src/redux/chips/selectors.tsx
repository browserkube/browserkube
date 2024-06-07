import { type ChipsSliceState } from './chipsSlice';

export const selectChips = (state: { chips: ChipsSliceState }) => state.chips;
