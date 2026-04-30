import { useEffect, useState } from 'react';
import { formatTime } from '@shared/utils/getVideoTimeFormated';
import { lang } from '@app/constants';

const { startTime, duration } = lang.liveSession;

interface CountdownTimerProps {
  createdAt: number;
  sessionId: string;
}

export const CountdownTimer = ({ createdAt, sessionId }: CountdownTimerProps) => {
  const [remainingTime, setRemainingTime] = useState(duration);

  useEffect(() => {
    const updateTimer = () => {
      const elapsed = Date.now() - createdAt;
      const newTime = Math.max(duration - elapsed, 0);

      if (newTime === 0) {
        clearInterval(timerInterval);
      }

      setRemainingTime(newTime);
    };

    const timerInterval = setInterval(updateTimer, 1000);
    return () => {
      clearInterval(timerInterval);
    };
  }, [createdAt]);

  return <>{sessionId ? formatTime(remainingTime) : startTime}</>;
};

export const useIsTimerActive = (createdAt: number) => {
  const [active, setActive] = useState(true);

  useEffect(() => {
    const check = () => {
      const elapsed = Date.now() - createdAt;
      if (duration - elapsed <= 0) {
        setActive(false);
        clearInterval(id);
      }
    };
    const id = setInterval(check, 1000);
    return () => clearInterval(id);
  }, [createdAt]);

  return active;
};
