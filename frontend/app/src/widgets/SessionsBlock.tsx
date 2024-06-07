import { type ChangeEvent, useEffect, useMemo, useRef, useState, useCallback } from 'react';
import { IconButton } from '@mui/material';
import { useSelector } from 'react-redux';
import { createPortal } from 'react-dom';
import { Tabs } from 'app/theme/TabKit/Tabs';
import { Tab } from 'app/theme/TabKit/Tab';
import { Divider, DividerType } from 'widgets/Divider';
import { getSessions } from '@redux/sessions/sessionsSelectors';
import { getMaxTimeout } from '@redux/sessionStatus/sessionStatusSelectors';

import { getTerminatedSessions } from '@redux/terminatedSessions/selectors';
import { useAppDispatch } from '@hooks/useAppDispatch';
import { fetchTerminatedSessions } from '@redux/terminatedSessions/thunk';
import { fetchSessions } from '@redux/sessions/sessionsThunk';
import { getGridRows } from '@shared/utils/getGridRows';
import { getBrowserChip } from '@shared/utils/getBrowserChip';
import { type SessionLine, type SessionLinesArray } from '@shared/types/sessions';
import { getSessionStatusIcon } from '@shared/utils/getSessionStatusIcon';
import { getOSChip } from '@shared/utils/getOsChip';
import { getSessionMode } from '@shared/utils/getSessionMode';
import { DropDown, DropDownKit } from 'app/theme/DropButtonKit/DropButtonKit';
import { SessionManual } from '@shared/chips/sessionManual';
import { SessionAuto } from '@shared/chips/sessionAuto';
import { getSessionsFiltered } from '@shared/utils/getSessionsFiltered';
import { MODAL_TYPE } from '@shared/types/UI';
import { clearAllChips, clearChips } from '@redux/chips/chipsSlice';
import { openModal } from '@redux/UI/UISlice';
import { makeFilter } from '@shared/utils/formatFilterState';
import { getChip } from '@shared/utils/getChips';
import { getChipArray } from '@shared/utils/getFilterAreaChips';
import { selectChips } from '@redux/chips/selectors';
import { ArchiveIcon } from '@shared/icons/archiveIcon';
import { FilterIcon } from '@shared/icons/filterIcon';
import { SearchIcon } from '@shared/icons/searchIcon';

import '@app/theme/fonts.css';
import { ArrowDown } from '@shared/icons/arrowDown';
import { AddFilterIcon } from '@shared/icons/addFilterIcon';
import { CloseFilterAreaIcon } from '@shared/icons/closeFilterAreaIcon';
import { MagnifierIcon } from '@shared/icons/magnifierIcon';
import { MagnifierIconActive } from '@shared/icons/magnifierIconActive';
import styles from '@pages/LiveSessions/LiveSession.module.scss';
import { useSession } from '../pages/LiveSessions/SessionContext';
import { InstanceUsage } from './InstanceUsage';
import { FilterBlock } from './FilterBlock';
import { ArchiveModal } from './ArchieveModal';

// filter styles

const leftBlock = {
  display: 'flex',
  gap: '8px',
  flex: '1',
  alignItems: 'center',
} as const;

const filterChips = {
  display: 'flex',
  gap: '8px',
  flexWrap: 'wrap',
} as const;

//

const headerBtns = {
  display: 'flex',
  gap: '8px',
  alignItems: 'center',
} as const;

const textStyle = {
  fontSize: '20px',
  lineHeight: '31px',
} as const;

const tabsStyle = {
  padding: '22px 16px 0',
  backgroundColor: 'white',
} as const;

const tabStyle = {
  marginRight: '24px',
  paddingBottom: '14px',
} as const;

const lastTabStyle = {
  paddingBottom: '14px',
} as const;

const middleHeaderIcons = {
  display: 'flex',
  gap: '20px',
} as const;

const tabContainerStyle = {
  width: '100%',
} as const;

const checkedSessionStyle = {
  borderRadius: '4px',
  border: '2px solid #00B0D1',
} as const;

export interface FilterBlockPosition {
  top: number;
  left: number;
}

export const SessionsBlock = () => {
  const dispatch = useAppDispatch();
  const { session, setSession } = useSession();
  const terminatedSessions = useSelector(getTerminatedSessions);
  const sessions = useSelector(getSessions);
  const timeout = useSelector(getMaxTimeout);
  const rowsActive = useMemo(() => getGridRows(sessions, timeout), [sessions, timeout]);
  const sessionData = [...rowsActive, ...terminatedSessions];
  const chipState = useSelector(selectChips);

  const [tab, setTab] = useState('allSessions');
  const [currentChip, setCurrentChip] = useState<string>('');

  // filterArea
  const [isFilterArea, setIsFilterArea] = useState<boolean>(false);
  const [toggleFilter, isToggleFilter] = useState<boolean>(false);
  const [toggleArchieve, isToggleArchieve] = useState<boolean>(false);

  // search state
  const [isSearch, toggleIsSearch] = useState<boolean>(false);
  const [search, setSeacrh] = useState<string>('');
  const inputRef = useRef<HTMLInputElement | null>(null);

  const filterListRef = useRef<HTMLDivElement | null>(null);

  // extra filter dropdown
  const addButtonRef = useRef<HTMLDivElement | null>(null);
  const [isExtraFilterBlock, setIsExtraFilterBlock] = useState<boolean>(false);
  const [filterBlockPosition, setFilterBlockPosition] = useState<FilterBlockPosition>({
    top: 0,
    left: 0,
  });

  const calculateFilterBlockPosition = () => {
    if (addButtonRef.current) {
      const addButtonRect = addButtonRef.current.getBoundingClientRect();
      const top = addButtonRect.bottom + window.scrollY - 90;
      const left = addButtonRect.left + window.screenX - 25;
      setFilterBlockPosition({ top, left });
    }
  };

  const handleToggleExtraFilterBlock = () => {
    isToggleFilter(false);
    setIsExtraFilterBlock(!isExtraFilterBlock);

    if (!isExtraFilterBlock) {
      calculateFilterBlockPosition();
    }
  };
  //

  const selectSession = async (session: SessionLine) => {
    if (session.state !== 'Terminated' && sessions) {
      const sessionToSet = sessions.byId[session.id];
      setSession(sessionToSet);
    } else {
      setSession(session);
    }
  };

  const handleRemoveChip = (index: number) => {
    const chipArray = getChipArray(chipState);
    if (chipArray.length === 0) {
      return;
    }
    const chipName = Object.keys(chipArray[index])[0];
    dispatch(clearChips({ currentChip: chipName }));
  };

  const selectTab = (e: React.SyntheticEvent, value: string) => {
    setTab(value);
  };

  const handleFilter = () => {
    isToggleFilter(!toggleFilter);
  };

  const handleSeacrh = () => {
    toggleIsSearch(!isSearch);
  };

  const handleSearchFilter = (e: ChangeEvent<HTMLInputElement>) => {
    setSeacrh(e.target.value);
  };

  const handleCloseFilterArea = () => {
    setIsFilterArea(false);
    dispatch(clearAllChips());
  };

  const openSessionModal = (mode: string) => {
    const key = `${Date.now()}`;
    const modalType = mode === 'manual' ? MODAL_TYPE.CREATE_SESSION : MODAL_TYPE.GENERATE_CODE_SNIPPET;
    dispatch(openModal({ id: key, component: modalType }));
  };

  const handleArchieve = () => {
    isToggleArchieve(!toggleArchieve);
  };

  useEffect(() => {
    const getSessions = async () => {
      try {
        await dispatch(fetchTerminatedSessions());
      } catch (e) {
        console.log(e);
      }
    };
    void getSessions();
    void dispatch(fetchSessions());
  }, [dispatch]);

  const filteredSessionsList = useCallback(
    (arr: SessionLine[]) => {
      if (!search) {
        return arr;
      }
      return arr.filter(({ name }) => {
        return name.includes(search);
      });
    },
    [search]
  );

  // useMemo for array
  const CreateListOfSessions = (props: SessionLinesArray) => {
    let filteredSessions = props.sessionArr;
    if (getChipArray(chipState).length > 0 && filteredSessions.length > 0) {
      const mappedFilter = makeFilter(chipState);
      filteredSessions = getSessionsFiltered(mappedFilter, filteredSessions);
    }

    const result = filteredSessionsList(filteredSessions);
    return (
      <div className={styles.sessions_list_container}>
        {result.length > 0 &&
          result.map((sessionLine, index: number) => {
            return (
              <div key={`${sessionLine.id}_${index}`}>
                <div
                  style={session?.id === sessionLine.id ? checkedSessionStyle : {}}
                  className={styles.session_line}
                  onClick={() => {
                    void selectSession(sessionLine);
                  }}>
                  <div className={styles.line_status}>{getSessionStatusIcon(sessionLine.state)}</div>
                  <div className={styles.line_name}>{sessionLine.name}</div>
                  <div className={styles.chips_line_container}>
                    <div className={styles.line_chips}>{getOSChip(sessionLine.platformName)}</div>
                    <div className={styles.line_chips}>{getBrowserChip(sessionLine.browser)}</div>
                    <div className={styles.line_chips}>{getSessionMode(sessionLine.manual)}</div>
                  </div>
                </div>
              </div>
            );
          })}
      </div>
    );
  };

  useEffect(() => {
    if (inputRef.current && isSearch) {
      inputRef.current.focus();
    }
  }, [isSearch]);

  return (
    <div className={styles.sessions_container}>
      <div className={styles.header}>
        <div style={textStyle}>Sessions</div>
        <div style={headerBtns}>
          <InstanceUsage />
          <IconButton onClick={handleArchieve}>
            <ArchiveIcon />
          </IconButton>
          {toggleArchieve &&
            createPortal(
              <ArchiveModal
                onClose={() => {
                  isToggleArchieve(false);
                }}
              />,
              document.body
            )}
        </div>
      </div>
      <div className={styles.header_middle}>
        <div style={middleHeaderIcons}>
          <IconButton onClick={handleFilter}>
            <FilterIcon />
          </IconButton>
          <IconButton onClick={handleSeacrh}>{isSearch ? <MagnifierIconActive /> : <SearchIcon />}</IconButton>
        </div>
        <DropDownKit icon={ArrowDown({ color: 'white' })} placeholder={'Create Session'}>
          <DropDown
            key={'manual_dropdown'}
            value={'Manual'}
            onCustomClick={() => {
              openSessionModal('manual');
            }}
            icon={<SessionManual />}
          />
          <DropDown
            key={'auto_dropdown'}
            value={'Automatic'}
            onCustomClick={() => {
              openSessionModal('auto');
            }}
            icon={<SessionAuto />}
          />
        </DropDownKit>
      </div>
      {isFilterArea && <Divider type={DividerType.HORIZON} />}
      <div ref={filterListRef} className={isFilterArea ? styles.filterArea : ''}>
        {isFilterArea && (
          <>
            <div className={styles.tooltipCursor} />
            <div className={styles.filterInputArea}>
              <div style={leftBlock}>
                <div style={filterChips}>
                  {getChip(getChipArray(chipState), 'filter_area_chips', currentChip, handleRemoveChip)}
                  <div className={styles.add_filter_btn} ref={addButtonRef} onClick={handleToggleExtraFilterBlock}>
                    <AddFilterIcon />
                  </div>
                </div>
              </div>
              <IconButton onClick={handleCloseFilterArea} className={styles.close_filter_area_btn}>
                <CloseFilterAreaIcon />
              </IconButton>
            </div>
          </>
        )}
        <FilterBlock
          props={{
            currentChip,
            setCurrentChip,
            toggleFilter,
            isToggleFilter,
            filterListRef,
            isFilterArea,
            setIsFilterArea,
            filterBlockPosition,
            setIsExtraFilterBlock,
            isExtraFilterBlock,
          }}
        />
      </div>
      {isSearch && (
        <div className={styles.inner_container}>
          <div className={styles.input_container}>
            <MagnifierIcon />
            <input
              type="text"
              ref={inputRef}
              className={styles.search_input}
              placeholder="Filled Text"
              value={search}
              onInput={handleSearchFilter}></input>
            <div className={styles.close_btn} onClick={handleSeacrh}>
              <CloseFilterAreaIcon />
            </div>
          </div>
        </div>
      )}
      <div style={tabContainerStyle}>
        <Tabs value={tab} onChange={selectTab} style={tabsStyle}>
          <Tab value={'allSessions'} label="All Sessions" style={tabStyle} />
          <Tab value={'active'} label="Active" style={tabStyle} />
          <Tab value={'terminated'} label="Terminated" style={lastTabStyle} />
        </Tabs>
      </div>
      <Divider type={DividerType.HORIZON} />
      {tab === 'allSessions' && <CreateListOfSessions sessionArr={sessionData} />}
      {tab === 'active' && <CreateListOfSessions sessionArr={rowsActive} />}
      {tab === 'terminated' && <CreateListOfSessions sessionArr={terminatedSessions} />}
    </div>
  );
};
