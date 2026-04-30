import IconButton from '@mui/material/IconButton';

import React, { useMemo, useState } from 'react';
import { useSelector } from 'react-redux';
import { createPortal } from 'react-dom';
import { ExpandIcon } from '@shared/icons/expandIcon';
import { ScreenshotIcon } from '@shared/icons/screenshotIcon';
import { RecordingIcon } from '@shared/icons/recordingCircle';
import { SessionPending } from '@shared/icons/SessionPending';
import { getActiveSessionId } from '@redux/sessionDetails/selectors';
import { getSessions } from '@redux/sessions/sessionsSelectors';
import { LockButton } from '@pages/VncPanel/components/LockButton/LockButton';
import { useVncPanel } from '@pages/VncPanel/useVncPanel';
import { ActiveVncPanel } from '@pages/VncPanel/components/ActiveVncPanel/ActiveVncPanel';
import { handleScreenshot } from '@shared/utils/createScreenshot';
import { CountdownTimer, useIsTimerActive } from '@components/CountdownTimer/CountdownTimer';
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

const pendingContainer = {
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'center',
  alignItems: 'center',
  height: '100%',
  gap: '16px',
  color: '#6D7F8B',
  fontSize: '14px',
} as const;

const pendingSpinner = {
  width: '40px',
  height: '40px',
} as const;

export const LiveSession = React.memo(function LiveSession() {
  const activeSessionId = useSelector(getActiveSessionId);
  const { isInitialized } = useVncPanel(activeSessionId);
  const sessions = useSelector(getSessions);

  const sessionState = sessions?.byId[activeSessionId]?.state;
  const SESSION_CREATED_AT = sessions.byId[activeSessionId]?.createdAt ?? 0;
  const IS_SESSION_RUNNING = sessionState === 'running';
  const IS_SESSION_AUTO = !sessions?.byId[activeSessionId]?.manual;
  const isTimerActive = useIsTimerActive(SESSION_CREATED_AT);

  const [isExpand, toggleIsExpand] = useState<boolean>(false);

  const handleExpandMode = () => {
    toggleIsExpand(!isExpand);
  };

  const vncBlock = useMemo(() => {
    if (IS_SESSION_RUNNING && !isInitialized) {
      return <div className={styles.vncPanelLoader} />;
    }

    if (IS_SESSION_RUNNING) {
      return <div style={screenShotContainer}><ActiveVncPanel /></div>;
    }

    return (
      <div style={pendingContainer}>
        <div style={pendingSpinner}>
          <div style={{ width: '100%', height: '100%', transform: 'scale(2.5)', transformOrigin: 'top left' }}>
            <SessionPending />
          </div>
        </div>
        <div>Session is starting…</div>
      </div>
    );
  }, [activeSessionId, isInitialized, IS_SESSION_RUNNING, isExpand]);

  return (
    <div style={containerColor}>
      {IS_SESSION_RUNNING ? (
        <>
          <div className={styles.record_container}>
            <div className={styles.record_bar}>
              <RecordingIcon isAnimation={isTimerActive} />
              <div>Recording 00:00:00</div>
            </div>
            <div>Time Left: <CountdownTimer createdAt={SESSION_CREATED_AT} sessionId={activeSessionId} /></div>
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
          {vncBlock}
        </>
      ) : (
        vncBlock
      )}
    </div>
  );
});
