import { type BrowserData } from '@shared/types/createSession';

export const getResolutionOptions = (browserData: BrowserData | null, selectedVersion: string) => {
  if (!browserData) return [];
  const foundBrowserVerion = browserData.versions.find((version) => version.value === selectedVersion);
  if (!foundBrowserVerion) return [];
  return foundBrowserVerion.resolutions.map((resolution) => ({
    label: resolution,
    value: resolution,
  }));
};
