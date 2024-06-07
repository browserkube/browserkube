import { createSlice, type PayloadAction } from '@reduxjs/toolkit';
import { REDUCER_STATUS, type ReducerType } from '@shared/types/reducerType';
import { type Browser } from '@shared/types/browsers';
import { fetchBrowsers } from '@redux/browsers/browsersThunk';

type BrowsersStore = ReducerType<{
  browsers: Browser[];
}>;

const initialState: BrowsersStore = {
  data: {
    browsers: [],
  },
  status: REDUCER_STATUS.IDLE,
  error: null,
};

const Browsers = createSlice({
  name: 'browsers',
  initialState,
  reducers: {},
  extraReducers(builder) {
    builder
      .addCase(fetchBrowsers.pending, (state) => {
        state.status = REDUCER_STATUS.PENDING;
        state.error = null;
      })
      .addCase(fetchBrowsers.fulfilled, (state, action: PayloadAction<Browser[]>) => {
        state.data.browsers = action.payload;
        state.status = REDUCER_STATUS.FULFILLED;
      })
      .addCase(fetchBrowsers.rejected, (state) => {
        state.status = REDUCER_STATUS.REJECTED;
      });
  },
});

export const { reducer: browsersReducer } = Browsers;
