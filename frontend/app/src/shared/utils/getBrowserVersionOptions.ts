import { type OSData } from '@shared/types/createSession';

export const getBrowserVersionOptions = (OSData: OSData | null, selectedBrowser: string) => {
  const emptyData = {
    browserVersionOptions: [],
    foundBrowserVersion: null,
  };
  if (!OSData) return emptyData;
  const foundBrowserVersion = OSData?.browsers.find(
    (browser) => browser.value === selectedBrowser || browser.label === selectedBrowser
  );
  if (!foundBrowserVersion) return emptyData;
  return {
    browserVersionOptions: foundBrowserVersion.versions.map((version) => ({
      label: version.label ?? version.value,
      value: version.value,
    })),
    foundBrowserVersion,
  };
};
