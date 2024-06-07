import { type ChipsSliceState } from '@redux/chips/chipsSlice';
import { type ChipType } from './getChips';

export const hideExtraChipValue = (chipKey: string, chipValue: string) => {
  const chipValueArr = chipValue.split(', ');

  if (chipValueArr.length <= 3) {
    return { [chipKey]: chipValue };
  } else {
    const sortedChipValueArr = chipValueArr.sort((a, b) => {
      return a.length - b.length || a.localeCompare(b);
    });

    let shortenedChipValue = sortedChipValueArr.slice(0, 3);
    let remainingValues = sortedChipValueArr.length - 3;
    if (shortenedChipValue.join().length >= 25) {
      shortenedChipValue = sortedChipValueArr.slice(0, 2);
      remainingValues = sortedChipValueArr.length - 2;
    }
    const result = `${shortenedChipValue.join(', ')} ... + ${remainingValues}`;

    return { [chipKey]: result };
  }
};

const formatFromCheckedToFilter = (currentChip: string, checkedState: Record<string, boolean>) => {
  return {
    [currentChip]: Object.keys(checkedState)
      .filter((key) => checkedState[key])
      .join(', '),
  };
};

export const getChipArray = (chipState: ChipsSliceState): ChipType[] => {
  const result: ChipType[] = [];
  const chipStateKeys = Object.keys(chipState);
  chipStateKeys.forEach((chipName: string) => {
    const chipObject = chipState[chipName];
    if (Object.values(chipObject.values).includes(true)) {
      const chipPair = formatFromCheckedToFilter(chipName, chipObject.values);
      const chipNameResult = Object.entries(chipPair).map(([chipName, chipValue]) =>
        hideExtraChipValue(chipName, chipValue)
      );
      result.push(...chipNameResult);
    }
  });

  return result;
};
