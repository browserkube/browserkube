import IconButton from '@mui/material/IconButton';

import React, { useEffect, useMemo, useState } from 'react';
import { useSelector } from 'react-redux';
import { createPortal } from 'react-dom';
import { lang } from '@app/constants';
import { ExpandIcon } from '@shared/icons/expandIcon';
import { ScreenshotIcon } from '@shared/icons/screenshotIcon';
import { RecordingIcon } from '@shared/icons/recordingCircle';
import { getActiveSessionId } from '@redux/sessionDetails/selectors';
import { formatTime } from '@shared/utils/getVideoTimeFormated';
import { getSessions } from '@redux/sessions/sessionsSelectors';
import { LockButton } from '@pages/VncPanel/components/LockButton/LockButton';
import { useVncPanel } from '@pages/VncPanel/useVncPanel';
import { ActiveVncPanel } from '@pages/VncPanel/components/ActiveVncPanel/ActiveVncPanel';
import { handleScreenshot } from '@shared/utils/createScreenshot';
import styles from '@pages/LiveSessions/LiveSession.module.scss';
import { ExpandMode } from './ExpandMode';

const containerColor = {
  backgroundColor: '#F7F7F8',
  height: '100%',
  width: '100%',
} as const;

const screenShotContainer = {
  padding: '12px 32px',
  display: 'flex',
  justifyContent: 'center',
  height: '100%',
  alignItems: 'center',
} as const;

const { startTime, duration } = lang.liveSession;

export const LiveSession = React.memo(function LiveSession() {
  const activeSessionId = useSelector(getActiveSessionId);
  const { isInitialized } = useVncPanel(activeSessionId);
  const sessions = useSelector(getSessions);

  const SESSION_CREATED_AT = sessions.byId[activeSessionId]?.createdAt ?? 0;
  const IS_SESSION_RUNNING = sessions?.byId[activeSessionId]?.state === 'running';
  const IS_SESSION_AUTO = !sessions?.byId[activeSessionId]?.manual;

  const [remainingTime, setRemainingTime] = useState(duration);
  const [isExpand, toggleIsExpand] = useState<boolean>(false);

  const handleExpandMode = () => {
    toggleIsExpand(!isExpand);
  };

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

  const vncBlock = useMemo(() => {
    if (!isInitialized) {
      return <div className={styles.vncPanelLoader} />;
    }

    return <div style={screenShotContainer}>{IS_SESSION_RUNNING && <ActiveVncPanel />}</div>;
  }, [activeSessionId, isInitialized, isExpand]);

  return (
    <div style={containerColor}>
      <div className={styles.record_container}>
        <div className={styles.record_bar}>
          <RecordingIcon isAnimation={remainingTime !== 0} />
          <div>Recording 00:00:00</div>
        </div>
        <div>Time Left: {activeSessionId ? formatTime(remainingTime) : startTime}</div>
        <div>
          <LockButton />
          <IconButton
            disabled={IS_SESSION_AUTO}
            onClick={() => {
              void handleScreenshot(activeSessionId);
            }}>
            <ScreenshotIcon />
          </IconButton>
          <IconButton disabled={!activeSessionId} onClick={handleExpandMode}>
            <ExpandIcon />
          </IconButton>
          {isExpand &&
            createPortal(
              <>
                <ExpandMode
                  onClose={() => {
                    toggleIsExpand(false);
                  }}
                />
              </>,
              document.body
            )}
        </div>
      </div>
      {IS_SESSION_RUNNING && vncBlock}
    </div>
  );
});
