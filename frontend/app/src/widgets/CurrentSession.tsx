import { useEffect, useMemo, useState } from 'react';
import { Button } from '@reportportal/ui-kit';
import styles from '@pages/LiveSessions/LiveSession.module.scss';
import { Divider, DividerType } from 'widgets/Divider';
import { deleteWdSession } from '@redux/sessions/sessionsThunk';
import { useAppDispatch } from '@hooks/useAppDispatch';
import { getChip } from '@shared/utils/getChips';
import { getStringFormated } from '@shared/utils/getStringFormated';
import { Tabs } from 'app/theme/TabKit/Tabs';
import { Tab } from 'app/theme/TabKit/Tab';
import { useSession } from '../pages/LiveSessions/SessionContext';
import { LiveSession } from './LiveSession';
import { SessionDetails } from './SessionDetails';
import { AttachmentsTab } from './Attachments';

export const tabStyle = {
  padding: '22px 0 14px',
  marginRight: '24px',
} as const;

export const lastTabStyle = {
  padding: '22px 0 14px',
} as const;

const currentSessionTitle = {
  fontSize: '20px',
  padding: '0 32px',
} as const;

const deleteBtn = {
  padding: '7px 16px',
  width: 'max-content',
} as const;

export const CurrentSession = () => {
  const dispatch = useAppDispatch();
  const { session } = useSession();
  const [tab, setTab] = useState('liveSession');
  const SELECTED_SESSION_IS_TERMINATED = session && session.state === 'Terminated';

  const updateSetTabfromChild = (newState: string) => {
    setTab(newState);
  };

  const handleChange = (event: React.SyntheticEvent, tabName: string) => {
    setTab(tabName);
  };

  const handleTerminate = () => {
    void dispatch(deleteWdSession({ sessionId: session.id }));
  };

  const chipsAvailiableArr = useMemo(() => {
    if (!session) {
      return;
    }
    const { type, browser, browserVersion, platformName, manual, screenResolution } = session;
    const chipArray: Array<Record<string, string>> = [];

    if (screenResolution) {
      chipArray.push({ screenResolution });
    }

    if (platformName) {
      chipArray.push({ platformName: getStringFormated(platformName) });
    }

    chipArray.push({ manual: manual ? 'Manual' : 'Auto' });
    chipArray.push({ [type.toLowerCase()]: getStringFormated(type) });
    chipArray.push({ [browser]: browserVersion });
    return chipArray;
  }, [session]);

  const CHIPS_AVAILIABLE = chipsAvailiableArr && chipsAvailiableArr.length > 0;

  useEffect(() => {
    if (SELECTED_SESSION_IS_TERMINATED) {
      setTab('details');
    } else {
      setTab('liveSession');
    }
  }, [session?.state]);

  return (
    <div className={styles.currrentSession_container}>
      <div className={styles.activeSession_container}>
        <div style={currentSessionTitle}>{session?.name || 'Choose session'}</div>
        <div className={styles.panel_container}>
          <div className={styles.left_panel}>
            <Tabs value={tab} onChange={handleChange}>
              {!SELECTED_SESSION_IS_TERMINATED && <Tab value={'liveSession'} label="Live Session" style={tabStyle} />}
              <Tab disabled={!session} value={'details'} label="Details" style={tabStyle} />
              <Tab value={'attachments'} label="Attachments" style={lastTabStyle} />
            </Tabs>
          </div>
          <div className={styles.right_panel}>
            <div className={styles.chips_session_container}>
              {CHIPS_AVAILIABLE && getChip(chipsAvailiableArr, session.id)}
            </div>
            <Button
              style={deleteBtn}
              disabled={session?.state === 'Terminated'}
              variant={'danger'}
              onClick={handleTerminate}>
              {session?.state !== 'Terminated' ? 'Terminate' : 'Delete Session'}
            </Button>
          </div>
        </div>
      </div>
      <Divider type={DividerType.HORIZON} />
      {tab === 'liveSession' && <LiveSession />}
      {tab === 'details' && <SessionDetails />}
      {tab === 'attachments' && (
        <AttachmentsTab setTab={updateSetTabfromChild} isTerminated={SELECTED_SESSION_IS_TERMINATED} />
      )}
    </div>
  );
};
