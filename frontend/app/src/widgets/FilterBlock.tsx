import cx from 'classnames';
import { IconButton } from '@mui/material';
import { useSelector } from 'react-redux';
import {
  type Dispatch,
  type SetStateAction,
  useEffect,
  useMemo,
  useRef,
  useState,
  type MutableRefObject,
  type ChangeEvent,
} from 'react';
import { Checkbox } from '@reportportal/ui-kit';
import styles from '@app/styles/filterArea.module.scss';
import { ButtonDrop, DropValue } from 'app/theme/ButtonDrop/ButtonDropKit';
import { Divider, DividerType } from 'widgets/Divider';
import { getBrowsers } from '@redux/browsers/browsersSelectors';
import { type ChipType, getChip } from '@shared/utils/getChips';
import { selectChips } from '@redux/chips/selectors';
import { useAppDispatch } from '@hooks/useAppDispatch';
import { addBrowserChips, addChip, clearChips, removeChip, selectAll, selectChip } from '@redux/chips/chipsSlice';
import { CloseFilterAreaIcon } from '@shared/icons/closeFilterAreaIcon';
import { MagnifierIcon } from '@shared/icons/magnifierIcon';
import { type FilterBlockPosition } from './SessionsBlock';

const overallFilterStatistic = {
  color: '#8D95A1',
  fontSize: '13px',
  paddingBottom: '8px',
  marginTop: '8px',
} as const;

const filterHeader = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
} as const;

const chipsList = {
  display: 'flex',
  flexWrap: 'wrap',
  gap: '4px',
  padding: '7px 0',
} as const;

const SESSION_TYPE_CHIP = ['auto', 'manual'];

interface FilterBlockProps {
  currentChip: string;
  setCurrentChip: Dispatch<SetStateAction<string>>;
  toggleFilter: boolean;
  isToggleFilter: Dispatch<SetStateAction<boolean>>;
  filterListRef: MutableRefObject<HTMLDivElement | null>;
  isFilterArea: boolean;
  setIsFilterArea: Dispatch<SetStateAction<boolean>>;
  filterBlockPosition: FilterBlockPosition;
  isExtraFilterBlock: boolean;
  setIsExtraFilterBlock: Dispatch<SetStateAction<boolean>>;
}

export const FilterBlock = ({ props }: { props: FilterBlockProps }) => {
  const {
    currentChip,
    toggleFilter,
    isToggleFilter,
    filterListRef,
    setCurrentChip,
    setIsFilterArea,
    isFilterArea,
    filterBlockPosition,
    isExtraFilterBlock,
    setIsExtraFilterBlock,
  } = props;
  const dispatch = useAppDispatch();
  const filterData = useSelector((s) => getBrowsers(s));
  const [selectedAll, setSelectedAll] = useState(false);
  const [isFIlterList, setIsFilterList] = useState(false);
  const inputRef = useRef<HTMLInputElement | null>(null);
  const chipState = useSelector(selectChips);
  const isSessionType = SESSION_TYPE_CHIP.includes(currentChip);
  const [filter, setFilter] = useState<string>('');

  const selectionRef = useRef<HTMLDivElement | null>(null);

  const filteredCheckboxes = useMemo(() => {
    if (chipState && currentChip) {
      if (!filter) {
        return Object.entries(chipState[currentChip].values);
      }
      return Object.entries(chipState[currentChip].values).filter(([chipName]) => {
        return chipName.includes(filter);
      });
    }
  }, [filter, chipState, currentChip]);

  const extraFilterBlock = {
    display: isExtraFilterBlock ? 'block' : 'none',
    position: 'absolute',
    top: `${filterBlockPosition?.top ?? 0}px`,
    left: `${filterBlockPosition?.left ?? 0}px`,
    width: 'max-content',
  } as const;

  const handleClickOutside = (event: MouseEvent) => {
    if (selectionRef.current) {
      const targetElement = event.target as Node;

      if (filterListRef?.current && !filterListRef?.current.contains(targetElement)) {
        setSelectedAll(false);
        setIsFilterList(false);
        isToggleFilter(false);
      }
    }
  };

  const handleRemoveChip = (index: number) => {
    const chipName = Object.values(formatedChipArray[index])[0];
    dispatch(removeChip({ chipName, currentChip }));
  };

  const chooseFilterValue = (result: Record<string, string>) => {
    const chosenChip = Object.keys(result)[0];
    setCurrentChip(chosenChip);
    dispatch(addChip({ currentChip: chosenChip }));
    isToggleFilter(!toggleFilter);
    setIsFilterList(true);
    setIsFilterArea(true);
  };

  const chooseExtraFilterValue = (result: Record<string, string>) => {
    const chosenChip = Object.keys(result)[0];
    setCurrentChip(chosenChip);
    dispatch(addChip({ currentChip: chosenChip }));
    isToggleFilter(false);
    setIsExtraFilterBlock(!isExtraFilterBlock);
    setIsFilterList(true);
    setIsFilterArea(true);
  };

  const selectCheckbox = (chipName: string, currentChip: string, chipValue: boolean) => {
    dispatch(selectChip({ chipName, currentChip, chipValue }));
  };

  const clearAllChips = () => {
    dispatch(clearChips({ currentChip }));
    setSelectedAll(!selectAll);
  };

  const handleFilter = (e: ChangeEvent<HTMLInputElement>) => {
    setFilter(e.target.value);
  };

  const selectAllCheckboxes = () => {
    dispatch(selectAll({ currentChip, chipValue: !selectedAll }));
    setSelectedAll(!selectedAll);
  };

  const getLabel = () => {
    const label = currentChip === 'screenResolution' ? 'Screen resolution' : chipState[currentChip].label;
    return (
      <>
        {`${label} version  `}
        <b>is</b>
      </>
    );
  };

  useEffect(() => {
    document.addEventListener('mousedown', handleClickOutside);

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  useEffect(() => {
    dispatch(addBrowserChips({ data: filterData }));
  }, [Object.values(filterData).length]);

  const formatedChipArray = useMemo(() => {
    const result: ChipType[] = [];
    if (!chipState[currentChip]) {
      return [];
    }
    Object.entries(chipState[currentChip].values).forEach(([chipName, chipValue]) => {
      if (chipValue && chipName !== ' ') {
        result.push({ [currentChip]: chipName });
      }
    });

    return result;
  }, [handleRemoveChip, selectAllCheckboxes, handleFilter]);

  const CHIPS_FILTER_EMPTY = formatedChipArray.length !== 0;

  useEffect(() => {
    if (inputRef.current) {
      inputRef.current.focus();
    }
  }, [currentChip]);

  return (
    <>
      {isExtraFilterBlock && (
        <div className="filter_block" style={isExtraFilterBlock !== undefined ? extraFilterBlock : {}}>
          <ButtonDrop title="Add Filter" onClick={chooseExtraFilterValue}>
            <DropValue title={'Browser & Version'} label={'browser'} value={['Chrome', 'Edge', 'Firefox']} />
            <DropValue title={'Session Type'} label={'sessionType'} value={['Auto', 'Manual']} />
            <DropValue title={'Screen Resolution'} label={'screenResolution'} />
          </ButtonDrop>
        </div>
      )}
      {toggleFilter && (
        <ButtonDrop title="Add Filter" onClick={chooseFilterValue}>
          <DropValue title={'Browser & Version'} label={'browser'} value={['Chrome', 'Edge', 'Firefox']} />
          <DropValue title={'Session Type'} label={'sessionType'} value={['Auto', 'Manual']} />
          <DropValue title={'Screen Resolution'} label={'screenResolution'} />
        </ButtonDrop>
      )}
      {isFIlterList && !isSessionType && isFilterArea && (
        <>
          <div className={styles.selectionContainer} ref={selectionRef}>
            <div style={filterHeader}>
              <div className={styles.filterListHeader}>{getLabel()}</div>
              <div className={styles.download} style={CHIPS_FILTER_EMPTY ? {} : { opacity: '50%' }}>
                <IconButton onClick={clearAllChips} disabled={!CHIPS_FILTER_EMPTY}>
                  <CloseFilterAreaIcon />
                </IconButton>
                <div>Clear All</div>
              </div>
            </div>
            <div style={chipsList}>{getChip(formatedChipArray, 'filter_list_chips', undefined, handleRemoveChip)}</div>
            <div className={styles.my_input}>
              <MagnifierIcon />
              <input
                type="text"
                ref={inputRef}
                value={filter}
                className={styles.default_input}
                placeholder="Search"
                onInput={handleFilter}></input>
            </div>
            <Checkbox className={cx(styles.checkboxItem)} value={selectedAll} onChange={selectAllCheckboxes}>
              All
            </Checkbox>
            <Divider type={DividerType.HORIZON} />
            <div className={styles.filterList}>
              {filteredCheckboxes
                ?.filter(([chipName]) => chipName !== ' ')
                .map(([chipName, chipValue], index) => {
                  return (
                    <Checkbox
                      className={cx(styles.checkboxItem)}
                      value={chipValue}
                      key={`${currentChip}_${chipName}_${index}`}
                      onChange={() => {
                        selectCheckbox(chipName, currentChip, chipValue);
                      }}>
                      {`${chipState[currentChip].label} ${chipName}`}
                    </Checkbox>
                  );
                })}
            </div>
            <Divider type={DividerType.HORIZON} />
            <div style={overallFilterStatistic}>
              {`${chipState[currentChip].counter} of ${
                Object.keys(chipState[currentChip].values).length - 1 ?? 0
              } selected`}
            </div>
            <div className={styles.tooltipCursor}></div>
          </div>
        </>
      )}
    </>
  );
};
