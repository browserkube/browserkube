import { useSelector } from 'react-redux';
import { BASE_URL } from '@shared/lib';
import { Logs } from '@components/Logs/Logs';
import { getVncSessionId, getVncSessionLogStatus } from '@redux/VncSession/VncSessionSelectors';

export const SwipeableTermninalLogs = () => {
  const sessionId = useSelector((state) => getVncSessionId(state));
  const logsUrl = `wss://${window.location.hostname}${BASE_URL}/logs/${sessionId}`;
  const isOpened = useSelector((state) => getVncSessionLogStatus(state));
  return <Logs logsUrl={logsUrl} isOpened={isOpened} />;
};
