import { useState, useEffect } from 'react';
import { IconButton } from '@mui/material';
import { useSelector } from 'react-redux';
import styles from '@app/styles/details.module.scss';
import { API_URL, api } from '@shared/api';
import { getActiveSessionId, getSessionDetails, getSessionDetailsLogs } from '@redux/sessionDetails/selectors';
import { useAppDispatch } from '@hooks/useAppDispatch';
import { fetchSessionDetailsLogs } from '@redux/sessionDetails/thunk';

import { type AttachmentsTabProps } from '@shared/types/UI';
import { ArrowDown } from '@shared/icons/arrowDown';
import { ArrowUp } from '@shared/icons/arrowUp';
import { DownloadIcon } from '@shared/icons/downloadIcon';
import { getSessionFileUrl } from '@shared/utils/getSessionFileUrl';
import { ScreenshotsComponent } from './ScreenshotsComponent';
import { CommandVideos } from './CommandsVideos';

export const textStyle = {
  fontSize: '14px',
  fontWeight: '500',
} as const;

interface ScreenshotResponce {
  screenshots: string[];
}

const getScreenshots = async (URL: string): Promise<string[]> => {
  try {
    const { screenshots }: ScreenshotResponce = await api.get({ url: URL });
    return screenshots;
  } catch (e) {
    console.log('error while fetching screenshot DATA', e);
    return [];
  }
};

export const AttachmentsTab = (props: AttachmentsTabProps) => {
  const activeSessionId = useSelector(getActiveSessionId);
  const dispatch = useAppDispatch();
  const sessionDetails = useSelector(getSessionDetails);
  const logs = useSelector(getSessionDetailsLogs);

  const logsFileUri = sessionDetails.logsRefAddr;
  const url = getSessionFileUrl(activeSessionId, logsFileUri);

  const [toggleScreenShots, setToggleScreenShots] = useState(true);
  const [screenShot, setScreenshot] = useState<string[]>([]);
  const SCREENSHOT_AVAILIABLE = screenShot && screenShot.length > 0;

  useEffect(() => {
    if (activeSessionId) {
      const fetchScreenshotData = async () => {
        const URL = `${API_URL.SESSIONS}${activeSessionId}${API_URL.SCREENSHOTS}`;
        try {
          const res: string[] = await getScreenshots(URL);
          setScreenshot(res);
        } catch (e) {
          console.log('Error, while making the screenshot', e);
        }
      };
      void fetchScreenshotData();
    }
  }, [activeSessionId]);

  const handleDetails = () => {
    setToggleScreenShots(!toggleScreenShots);
  };

  useEffect(() => {
    if (logs.shouldRefetch && logsFileUri) {
      void dispatch(fetchSessionDetailsLogs(url));
    }
  }, [dispatch, logs.shouldRefetch, logsFileUri, url]);

  return (
    <div className={styles.container}>
      <div className={styles.screenshot_section}>
        <div className={styles.details_section}>
          <IconButton onClick={handleDetails}>{(toggleScreenShots && <ArrowDown />) || <ArrowUp />}</IconButton>
          <div style={textStyle}>Screenshots</div>
        </div>
        <div className={styles.download} style={SCREENSHOT_AVAILIABLE ? {} : { opacity: '50%' }}>
          <div>Download all</div>
          <IconButton disabled={!SCREENSHOT_AVAILIABLE} onClick={() => null}>
            <DownloadIcon />
          </IconButton>
        </div>
      </div>
      {toggleScreenShots && SCREENSHOT_AVAILIABLE && <ScreenshotsComponent screeshotArr={screenShot} />}
      <CommandVideos {...props} />
    </div>
  );
};
