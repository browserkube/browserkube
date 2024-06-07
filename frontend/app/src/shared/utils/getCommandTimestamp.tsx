import { type Commands } from '@shared/types/sessions';
import { formatTime } from './getVideoTimeFormated';

export const getCommandTimestamp = (timestamp: string, data: Commands[]) => {
  const startTime = new Date(data[0].timestamp).getTime();
  const timeDiffer = new Date(timestamp).getTime() - startTime;
  const commandTime = formatTime(timeDiffer);
  const secondsToJump = timeDiffer / 1000;

  return { secondsToJump, commandTime };
};
