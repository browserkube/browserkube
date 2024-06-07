import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import { getStringFormated } from '@shared/utils/getStringFormated';
import { type Browser } from '@shared/types/browsers';

export interface ChipProps {
  values: Record<string, boolean>;
  counter: number;
  label: string;
}

export type ChipsSliceState = Record<string, ChipProps>;

const initialState: ChipsSliceState = {
  chrome: {
    values: {},
    counter: 0,
    label: 'Chrome',
  },
  firefox: {
    values: {},
    counter: 0,
    label: 'Firefox',
  },
  edge: {
    values: {},
    counter: 0,
    label: 'Edge',
  },
  screenResolution: {
    values: {},
    counter: 0,
    label: '',
  },
  auto: {
    values: {},
    label: 'Auto',
    counter: 0,
  },
  manual: {
    values: {},
    label: 'Manual',
    counter: 0,
  },
};

export const DEFAULT_VALUES = ['chrome', 'edge', 'firefox', 'screenResolution'];
const HIDE_EMPTY_CHIP = { ' ': false };
const SHOW_EMPTY_CHIP = { ' ': true };

const createInitialChipsState = (state: ChipsSliceState, data: Browser[]): ChipsSliceState => {
  const screenResolutionSet = new Set<string>();
  DEFAULT_VALUES.forEach((value) => {
    state[value].values = {};
    state[value].counter = 0;
  });

  data.forEach((browser) => {
    const { name, version, resolutions } = browser;
    if (DEFAULT_VALUES.includes(name)) {
      resolutions.forEach((resolution) => {
        screenResolutionSet.add(resolution);
      });
      state[name].values[version] = false;
      state[name].counter = 0;
    }
  });
  const screenResolutionArr = Array.from(screenResolutionSet);
  screenResolutionArr.forEach((resolution: string) => {
    state.screenResolution.values[resolution] = false;
  });

  return state;
};

const chipsSlice = createSlice({
  name: 'chips',
  initialState,
  reducers: {
    selectChip: (state, action: PayloadAction<{ chipName: string; currentChip: string; chipValue: boolean }>) => {
      const { chipName, currentChip, chipValue } = action.payload;
      state[currentChip].values = { ...state[currentChip].values, [chipName]: !chipValue, ...HIDE_EMPTY_CHIP };
      let amount = state[currentChip].counter;
      state[currentChip].counter = chipValue ? (amount -= 1) : (amount += 1);
    },
    selectAll: (state, action: PayloadAction<{ currentChip: string; chipValue: boolean }>) => {
      const { currentChip, chipValue } = action.payload;
      state[currentChip].values = Object.fromEntries(
        Object.keys(state[currentChip].values).map((chipName) => [chipName, chipValue])
      );
      // hide empty chip
      state[currentChip].values = { ...state[currentChip].values, ...HIDE_EMPTY_CHIP };
      const counter = chipValue ? Object.keys(state[currentChip].values).length : 0;
      state[currentChip].counter = counter;
    },
    removeChip: (state, action: PayloadAction<{ chipName: string; currentChip: string }>) => {
      const { chipName, currentChip } = action.payload;
      state[currentChip].values = { ...state[currentChip].values, [chipName]: false };
      state[currentChip].counter -= 1;
    },
    addChip: (state, action: PayloadAction<{ currentChip: string }>) => {
      const { currentChip } = action.payload;
      if (Object.values(state[currentChip].values).includes(true)) {
        return state;
      }
      state[currentChip].values = { ...state[currentChip].values, ...SHOW_EMPTY_CHIP };
      if (!DEFAULT_VALUES.includes(currentChip)) {
        state[currentChip].values = { [getStringFormated(currentChip)]: true };
      }
    },
    clearChips: (state, action: PayloadAction<{ currentChip: string }>) => {
      const { currentChip } = action.payload;
      state[currentChip].values = Object.fromEntries(
        Object.keys(state[currentChip].values).map((chipName) => [chipName, false])
      );
      state[currentChip].counter = 0;
    },
    clearAllChips: (state) => {
      Object.keys(state).forEach((currentChip: string) => {
        state[currentChip].counter = 0;
        state[currentChip].values = Object.fromEntries(
          Object.keys(state[currentChip].values).map((chipName) => [chipName, false])
        );
      });
    },
    addBrowserChips: (state, action: PayloadAction<{ data: Browser[] }>) => {
      state = createInitialChipsState(state, action.payload.data);
    },
  },
});

export const { selectChip, selectAll, removeChip, addChip, clearChips, addBrowserChips, clearAllChips } =
  chipsSlice.actions;

export const { reducer: chipsReducer } = chipsSlice;
