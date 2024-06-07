import { type ActiveSessionGridRowsType, type Sessions } from '@shared/types/sessions';

const toDateFormat = (timestamp: number): string => {
  const date = new Date(timestamp);
  return `${date.getHours()}h:${date.getMinutes()}m:${date.getSeconds()}s`;
};

export const getGridRows = (state: Sessions, timeout: number): ActiveSessionGridRowsType[] => {
  return Object.values(state.byId).map((session) => {
    return {
      id: session.id,
      availability: `${toDateFormat(session.createdAt)} - ${toDateFormat(session.createdAt + timeout)}`,
      browser: session.browser,
      browserVersion: session.browserVersion,
      image: session.image,
      manual: session.manual,
      name: session.name,
      platformName: session.platformName,
      state: session.state,
      screenResolution: session.screenResolution,
      type: session.type,
    };
  });
};
