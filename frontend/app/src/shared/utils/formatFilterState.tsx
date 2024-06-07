import { type ChipsSliceState } from '@redux/chips/chipsSlice';

export const makeFilter = (data: ChipsSliceState): Record<string, string[]> => {
  const filterObject: Record<string, string[]> = {};

  for (const [key, chip] of Object.entries(data)) {
    const chipValues = Object.entries(chip.values)
      .filter(([_, value]) => value)
      .map(([label]) => label);

    if (chipValues.length > 0) {
      filterObject[key] = chipValues;
    }
  }

  return filterObject;
};
