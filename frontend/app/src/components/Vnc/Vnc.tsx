import { useRef, type ElementRef, useMemo } from 'react';
import { VncScreen } from 'react-vnc';
import { useSelector } from 'react-redux';
import { getActiveSessionId } from '@redux/sessionDetails/selectors';
import { getSessions } from '@redux/sessions/sessionsSelectors';
import styles from './Vnc.module.scss';

export interface VncProps {
  vncUrl: string;
  locked: boolean;
  isExpand: boolean;
}

export const Vnc = ({ vncUrl, locked, isExpand }: VncProps) => {
  const activeSessionId = useSelector(getActiveSessionId);
  const sessions = useSelector(getSessions);
  const vncScreenRef = useRef<ElementRef<typeof VncScreen>>(null);
  const isValidUrl = () => vncUrl.startsWith('ws://') || vncUrl.startsWith('wss://');

  const vncPsw = useMemo(() => {
    if (!activeSessionId) {
      return '';
    }
    return sessions.byId[activeSessionId].vncPsw;
  }, [activeSessionId]);

  return (
    <div className={isExpand ? styles.vnc_extended : styles.vnc}>
      {isValidUrl() ? (
        <VncScreen
          viewOnly={locked}
          rfbOptions={{
            shared: false,
            credentials: {
              password: vncPsw ?? '',
            },
          }}
          url={vncUrl}
          scaleViewport
          className={styles.vncScreen}
          ref={vncScreenRef}
          loadingUI={<div className={styles.vncLoader} />}
        />
      ) : (
        <div>VNC URL is not valid.</div>
      )}
    </div>
  );
};
