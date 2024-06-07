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
import { formatTime } from '@shared/utils/getVideoTimeFormated';
import { getSessions } from '@redux/sessions/sessionsSelectors';
import { RecordingIcon } from '@shared/icons/recordingCircle';
import { ScreenshotIcon } from '@shared/icons/screenshotIcon';
import { LeaveIcon } from '@shared/icons/leaveIcon';
import { PinIcon } from '@shared/icons/pinIcon';
import { PinSelectedIcon } from '@shared/icons/pinSelectedIcon';
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
  liveSession: { startTime, duration },
} = lang;

export const ExpandMode = ({ onClose }: ExpandModeProps) => {
  const activeSessionId = useSelector(getActiveSessionId);
  const sessions = useSelector(getSessions);

  const [isPin, toggleIsPin] = useState<boolean>(false);
  const [containerStyles, setContainerStyles] = useState<string>(start);
  const [remainingTime, setRemainingTime] = useState(duration);

  const SESSION_CREATED_AT = sessions.byId[activeSessionId]?.createdAt ?? 0;

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

  useEffect(() => {
    const updateTimer = () => {
      const elapsed = Date.now() - SESSION_CREATED_AT;
      const newRemainingTime = Math.max(duration - elapsed, 0);

      if (newRemainingTime === 0) {
        console.info('Timer has reached zero, session will be closed soon');
        clearInterval(timerInterval);
      }

      setRemainingTime(newRemainingTime);
    };

    const timerInterval = setInterval(updateTimer, 1000);

    return () => {
      clearInterval(timerInterval);
    };
  }, [SESSION_CREATED_AT]);

  return (
    <div className={styles.expand_container}>
      <div className={styles[containerStyles]}>
        <div className={styles.left_block}>
          <RecordingIcon isAnimation={true} />
          <div>Recording 00:00:00</div>
        </div>
        <div>Time Left: {activeSessionId ? formatTime(remainingTime) : startTime}</div>
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
