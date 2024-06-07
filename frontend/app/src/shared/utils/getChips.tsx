/**
 *
 * @param chipsArray [key:string]: string
 * @param id for key prop
 * @param closeMode if we need close btn for current chip
 * @returns chip component
 */

import { IconButton } from '@mui/material';
import { LinuxIcon } from '@shared/chips/linux';
import { ChromeIcon } from '@shared/chips/chrome';
import { EdgeIcon } from '@shared/chips/edge';
import { FirefoxIcon } from '@shared/chips/firefox';
import { SessionManual } from '@shared/chips/sessionManual';
import { SessionAuto } from '@shared/chips/sessionAuto';
import { WebdriverIcon } from '@shared/chips/webdriver';
import { PlaywrightIcon } from '@shared/chips/playwright';
import { ScreenResolutionIcon } from '@shared/chips/screenResolution';
import styles from '@pages/LiveSessions/LiveSession.module.scss';
import { CloseChipIcon } from '@shared/chips/closeChipIcon';

const iconLibrary = {
  // platformName
  platformName: <LinuxIcon />,

  // browsers
  chrome: <ChromeIcon />,
  edge: <EdgeIcon />,
  firefox: <FirefoxIcon />,

  // session type
  auto: <SessionAuto />,
  manual: <SessionManual />,

  playwright: <PlaywrightIcon />,
  webdriver: <WebdriverIcon />,

  screenResolution: <ScreenResolutionIcon />,
};

export type ChipType = Record<string, string>;

const textStyle = {
  fontSize: '13px',
  fontWeight: '500',
  color: '#3F3F3F',
  width: 'max-content',
} as const;

const selectedGroup = {
  borderRadius: '4px',
  border: '2px solid #00B0D1',
};

export const getChip = (
  chipsArray: ChipType[],
  id: string,
  selectedChip?: string,
  onRemove?: (index: number) => void
): React.ReactNode => {
  const SELF_DELETING_CHIP = !!onRemove;
  return chipsArray.map((chip: ChipType, index: number) => {
    const [key] = Object.entries(chip);
    const focusedChip = selectedChip === key[0] ? selectedGroup : undefined;
    return (
      <div key={`${id}_${index}`} className={`${styles.chip_body} ${key[0]}`} style={focusedChip}>
        {iconLibrary[key[0] as keyof typeof iconLibrary]}
        <div style={textStyle}>{key[1]}</div>
        {SELF_DELETING_CHIP && (
          <IconButton
            onClick={() => {
              onRemove(index);
            }}>
            {CloseChipIcon()}
          </IconButton>
        )}
      </div>
    );
  });
};
