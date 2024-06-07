import { type SessionLine } from '@shared/types/sessions';

type Filter = Record<string, string[]>;

export const getSessionsFiltered = (filter: Filter, data: SessionLine[]): SessionLine[] => {
  return data.filter((session) => {
    return Object.entries(filter).every(([key, values]) => {
      switch (key) {
        case 'auto':
          return values.some((value) => !session.manual && value === 'Auto');
        case 'manual':
          return values.some((value) => session.manual && value === 'Manual');
        case session.browser:
          return values.includes(session.browserVersion);
        case 'screenResolution':
          return values.includes(session.screenResolution);
        default:
          return false;
      }
    });
  });
};
