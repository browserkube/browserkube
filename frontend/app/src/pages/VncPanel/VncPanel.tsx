import { useSelector } from 'react-redux';
import { getActiveSessionId } from '@redux/sessionDetails/selectors';
import { useVncPanel } from './useVncPanel';
import { TerminatedVncPanel } from './components/TerminatedVncPanel/TerminatedVncPanel';
import { ActiveVncPanel } from './components/ActiveVncPanel/ActiveVncPanel';
import styles from './VncPanel.module.scss';

export const VncPanel = () => {
  const activeSessionId = useSelector(getActiveSessionId);
  const { sessionId, isInitialized } = useVncPanel(activeSessionId);

  if (!isInitialized) {
    return <div className={styles.vncPanelLoader} />;
  }

  return sessionId ? <ActiveVncPanel /> : <TerminatedVncPanel />;
};
