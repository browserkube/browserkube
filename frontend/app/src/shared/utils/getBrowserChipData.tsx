import { type Browser } from '@shared/types/browsers';
import { DEFAULT_VALUES } from '@redux/chips/chipsSlice';

const deleteDuplicates = (array: string[]) => {
  const res = new Set<string>();
  array.forEach((element) => res.add(element));

  return Array.from(res);
};

export const getBrowserData = (data: Browser[]): Map<string, string[]> => {
  const res = new Map<string, string[]>();
  DEFAULT_VALUES.forEach((value) => res.set(value, []));

  data.forEach((element) => {
    const existingResolution = res.get('screenResolution') ?? [];
    const mergedResoluttions = [...existingResolution, ...element.resolutions];
    const uniqueResolutions = deleteDuplicates(mergedResoluttions);

    const existingBrowser = res.get(element.name) ?? [];
    const updatedBrowser = [...existingBrowser, ...[element.version]];
    res.set(element.name, updatedBrowser);
    res.set('screenResolution', uniqueResolutions);
  });

  return res;
};
