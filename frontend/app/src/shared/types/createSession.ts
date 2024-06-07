export interface FormValues {
  sessionName: string;
  platformName: string;
  browserName: string;
  browserVersion: string;
  screenResolution: string;
  recordVideo: boolean;
}

export interface SelectData {
  label: string;
  value: string;
}

export interface BrowserVersion {
  label?: string;
  value: string;
  image: string;
  resolutions: string[];
}

export interface BrowserData {
  label?: string;
  value: string;
  versions: BrowserVersion[];
}

export interface OSData extends SelectData {
  browsers: BrowserData[];
}

export type FormData = OSData[];
