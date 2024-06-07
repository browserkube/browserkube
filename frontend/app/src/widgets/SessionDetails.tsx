import { useState, useEffect } from 'react';
import { IconButton, CircularProgress } from '@mui/material';
import { useSelector } from 'react-redux';
import styles from '@app/styles/details.module.scss';
import { lang } from '@app/constants';
import {
  getActiveSessionId,
  getSessionDetails,
  getSessionDetailsLogs,
  getSessionStatus,
} from '@redux/sessionDetails/selectors';
import { useAppDispatch } from '@hooks/useAppDispatch';
import { fetchSessionDetailsLogs } from '@redux/sessionDetails/thunk';
import { REDUCER_STATUS } from '@shared/types/reducerType';

import { ArrowDown } from '@shared/icons/arrowDown';
import { ArrowUp } from '@shared/icons/arrowUp';
import { DownloadIcon } from '@shared/icons/downloadIcon';
import { getSessionFileUrl } from '@shared/utils/getSessionFileUrl';
import { useSession } from '../pages/LiveSessions/SessionContext';
import { DetailsInfo } from './DetailsInfo';

const textStyle = {
  fontSize: '14px',
  fontWeight: '500',
} as const;

export const SessionDetails = () => {
  const activeSessionId = useSelector(getActiveSessionId);
  const [toggleDetails, setToggleDetails] = useState(true);
  const [toggleLogs, setToggleLogs] = useState(true);
  const sessionDetails = useSelector(getSessionDetails);
  const dispatch = useAppDispatch();
  const logs = useSelector(getSessionDetailsLogs);
  const logsFileUri = sessionDetails.logsRefAddr;
  const sessionStatus = useSelector(getSessionStatus);
  const url = getSessionFileUrl(activeSessionId, logsFileUri);
  const isLoading = sessionStatus === REDUCER_STATUS.PENDING;
  const areLogsNotFound = !logs.text || sessionStatus === REDUCER_STATUS.REJECTED;

  const { session } = useSession();

  const handleDetails = () => {
    setToggleDetails(!toggleDetails);
  };
  const handleLogs = () => {
    setToggleLogs(!toggleLogs);
  };

  useEffect(() => {
    if (logs.shouldRefetch && logsFileUri) {
      void dispatch(fetchSessionDetailsLogs(url));
    }
  }, [dispatch, logs.shouldRefetch, logsFileUri, url]);

  if (isLoading) {
    return <CircularProgress />;
  }

  const handleDownload = () => {
    const blob = new Blob([logs.text], { type: 'text/plain' });
    const urlLogs = window.URL.createObjectURL(blob);

    const a = document.createElement('a');
    a.href = urlLogs;
    a.download = `${session.name}.txt`;
    document.body.appendChild(a);

    a.click();

    window.URL.revokeObjectURL(urlLogs);
    document.body.removeChild(a);
  };

  return (
    <div className={styles.container}>
      <div className={styles.details_section}>
        <IconButton onClick={handleDetails}>{(toggleDetails && <ArrowDown />) || <ArrowUp />}</IconButton>
        <div style={textStyle}>Details</div>
      </div>
      {toggleDetails && session && <DetailsInfo session={session} />}
      <div className={styles.logs_section}>
        <div className={styles.details_section}>
          <IconButton onClick={handleLogs}>{(toggleLogs && <ArrowDown />) || <ArrowUp />}</IconButton>
          <div style={textStyle}>Logs</div>
        </div>
        <div className={styles.download}>
          <div>Download</div>
          <IconButton disabled={areLogsNotFound} onClick={handleDownload}>
            <DownloadIcon />
          </IconButton>
        </div>
      </div>
      {toggleLogs && (
        <div className={styles.logs_area}>
          {areLogsNotFound ? (
            <div>{lang.sessionDetails.noLogs}</div>
          ) : (
            <pre className={styles.logs_text}>{logs.text}</pre>
          )}
        </div>
      )}
    </div>
  );
};
