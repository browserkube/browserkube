import { IconButton } from '@mui/material';
import { useEffect, useMemo, useState } from 'react';
import { useSelector } from 'react-redux';
import { lang } from '@app/constants';
import { getQuotesLimit, getStats, getStatus } from '@redux/sessionStatus/sessionStatusSelectors';
import { INSTANCE_COLOR } from '@shared/types/indicator';
import { useAppDispatch } from '@hooks/useAppDispatch';
import { REDUCER_STATUS } from '@shared/types/reducerType';
import { fetchSessionStatus } from '@redux/sessionStatus/sessionStatusThunk';
import { WorkloadIcon } from '@shared/icons/workloadIcon';
import styles from '@app/styles/shared.module.scss';

const { instanceIndicator } = lang.common;

const provideQuoteStats = (quotesLimit: number, all: number) => {
  const capacity = quotesLimit - all;
  const quote = all && quotesLimit > 0 ? all / quotesLimit : 0;
  return { limit: quotesLimit, capacity: capacity > 0 ? capacity : 0, quote: quote >= 1 ? 1 : quote };
};

export const InstanceUsage = () => {
  const status = useSelector(getStatus);
  const dispatch = useAppDispatch();
  useEffect(() => {
    if (status === REDUCER_STATUS.IDLE) {
      void dispatch(fetchSessionStatus());
    }
  }, [dispatch, status]);

  const qlimit = useSelector(getQuotesLimit);
  const [tooltip, isTooltip] = useState(false);
  const { all } = useSelector(getStats);
  const { quote, limit } = useMemo(() => provideQuoteStats(qlimit, all), [all, qlimit]);

  const instanceColor = useMemo(() => {
    if (limit === 0 || quote === 0) {
      return INSTANCE_COLOR.FREE;
    }
    if (quote > 0 && quote < 0.5) {
      return INSTANCE_COLOR.FREE;
    }
    if (quote >= 0.5 && quote < 0.99) {
      return INSTANCE_COLOR.MEDIUM;
    }
    return INSTANCE_COLOR.OVERLOAD;
  }, [quote, limit]);

  return (
    <div
      className={styles.tooltipContainer}
      onMouseEnter={() => {
        isTooltip(true);
      }}
      onMouseLeave={() => {
        isTooltip(false);
      }}>
      <IconButton>
        <WorkloadIcon color={instanceColor} />
      </IconButton>
      {tooltip && (
        <div className={styles.tooltipStyle}>
          {instanceIndicator} {all}/{limit}({quote * 100}%)
          <div className={styles.tooltipCursor} />
        </div>
      )}
    </div>
  );
};
