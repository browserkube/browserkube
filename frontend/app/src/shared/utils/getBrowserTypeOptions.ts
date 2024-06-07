import { type FormData } from '@shared/types/createSession';

export const getBrowserTypeOptions = (formData: FormData, selectedOS: string) => {
  const OSData = formData.find((os) => os.value === selectedOS);
  if (!OSData)
    return {
      browserOptions: [],
      foundOS: null,
    };
  return {
    browserOptions: OSData.browsers.map((browser) => ({
      label: browser.label ?? browser.value,
      value: browser.value,
    })),
    foundOS: OSData,
  };
};
