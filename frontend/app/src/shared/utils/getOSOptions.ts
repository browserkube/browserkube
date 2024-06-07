import { type FormData } from '@shared/types/createSession';

export const getOSOptions = (formData: FormData) => {
  return formData.map((os) => ({
    label: os.label ?? os.value,
    value: os.value,
  }));
};
