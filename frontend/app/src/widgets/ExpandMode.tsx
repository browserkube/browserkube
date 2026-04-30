import { IconButton } from '@mui/material';
import { useSelector } from 'react-redux';
import { Button } from '@reportportal/ui-kit';
import { useEffect, useState } from 'react';
import styles from '@app/styles/expandMode.module.scss';
import { lang } from '@app/constants';
import { LockButton } from '@pages/VncPanel/components/LockButton/LockButton';
import { getActiveSessionId } from '@redux/sessionDetails/selectors';
import { VncControl } from '@pages/VncPanel/components/Controls/VncControl';
import { handleScreenshot } from '@shared/utils/createScreenshot';
import { getSessions } from '@redux/sessions/sessionsSelectors';
import { RecordingIcon } from '@shared/icons/recordingCircle';
import { ScreenshotIcon } from '@shared/icons/screenshotIcon';
import { LeaveIcon } from '@shared/icons/leaveIcon';
import { PinIcon } from '@shared/icons/pinIcon';
import { PinSelectedIcon } from '@shared/icons/pinSelectedIcon';
import { CountdownTimer, useIsTimerActive } from '@components/CountdownTimer/CountdownTimer';
import { type OnCloseModal } from './ArchieveModal';

interface ExpandModeProps {
  onClose: OnCloseModal;
}

export const leaveBtn = {
  color: 'white',
  backgroundColor: '#1A9CB0',
  padding: '7px 16px',
  borderRadius: '3px',
  marginLeft: '24px',
} as const;

const {
  expand: {
    mode: { start, pin, hide },
  },
} = lang;

export const ExpandMode = ({ onClose }: ExpandModeProps) => {
  const activeSessionId = useSelector(getActiveSessionId);
  const sessions = useSelector(getSessions);

  const [isPin, toggleIsPin] = useState<boolean>(false);
  const [containerStyles, setContainerStyles] = useState<string>(start);

  const SESSION_CREATED_AT = sessions.byId[activeSessionId]?.createdAt ?? 0;
  const isTimerActive = useIsTimerActive(SESSION_CREATED_AT);

  const handlePin = () => {
    toggleIsPin((prevIsPin: boolean) => {
      const newState = !prevIsPin;
      setContainerStyles(newState ? pin : hide);
      return newState;
    });
  };

  useEffect(() => {
    const timeoutID = setTimeout(() => {
      setContainerStyles(hide);
    }, 1500);

    return () => {
      clearTimeout(timeoutID);
    };
  }, []);

  return (
    <div className={styles.expand_container}>
      <div className={styles[containerStyles]}>
        <div className={styles.left_block}>
          <RecordingIcon isAnimation={isTimerActive} />
          <div>Recording 00:00:00</div>
        </div>
        <div>Time Left: <CountdownTimer createdAt={SESSION_CREATED_AT} sessionId={activeSessionId} /></div>
        <div className={styles.right_block}>
          <IconButton onClick={handlePin}>{isPin ? <PinSelectedIcon /> : <PinIcon />}</IconButton>
          <LockButton />
          <IconButton
            onClick={() => {
              void handleScreenshot(activeSessionId);
            }}>
            <ScreenshotIcon />
          </IconButton>
          <Button className={'leaveIcon'} onClick={onClose} style={leaveBtn} icon={LeaveIcon('white')} variant="text">
            Leave Full Screen
          </Button>
        </div>
      </div>
      <div className={styles.vnc_container}>
        <VncControl isExpand={true} />
      </div>
    </div>
  );
};
