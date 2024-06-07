import { useSelector } from 'react-redux';
import { BASE_URL } from '@shared/lib';
import { Vnc } from '@components/Vnc/Vnc';
import { getVncSessionId, getVncSessionScreenStatus } from '@redux/VncSession/VncSessionSelectors';

export const VncControl = ({ isExpand }: { isExpand: boolean }) => {
  const isLocked = useSelector((state) => getVncSessionScreenStatus(state));
  const sessionId = useSelector((state) => getVncSessionId(state));
  const isHttps = window.location.protocol === 'https:';
  const prefix = isHttps ? 'wss' : 'ws';
  const vncUrl = `${prefix}://${window.location.hostname}${BASE_URL}/vnc/${sessionId}`;

  return <Vnc key={`locked-${String(isLocked)} ${sessionId}`} locked={isLocked} vncUrl={vncUrl} isExpand={isExpand} />;
};
