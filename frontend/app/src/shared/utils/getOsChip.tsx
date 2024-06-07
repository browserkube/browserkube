import { LinuxIcon } from '@shared/chips/linux';

export const getOSChip = (platformName: string) => {
  let platformIcon;

  switch (platformName) {
    case 'linux':
      platformIcon = <LinuxIcon />;
      break;
    default:
      platformIcon = <LinuxIcon />; // remove after [BE] will be done
      break;
  }

  return platformIcon;
};
