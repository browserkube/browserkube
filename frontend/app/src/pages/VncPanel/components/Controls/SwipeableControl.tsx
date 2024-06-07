import { type SyntheticEvent } from 'react';
import { SwipeableDrawer } from '@mui/material';
import { useSelector } from 'react-redux';
import { useAppDispatch } from '@hooks/useAppDispatch';
import { getVncSession } from '@redux/VncSession/VncSessionSelectors';
import { updateVncSessionLogStatus } from '@redux/VncSession/VncSessionSlice';
import { SwipeableTermninalLogs } from './SwipeableTermninalLogs';

export const SwipeableControl = () => {
  const dispatch = useAppDispatch();
  const { swipeableOpened, swipeableAnchor } = useSelector((state) => getVncSession(state));

  const handleToggle = (expanded: boolean) => (e: SyntheticEvent) => {
    e.preventDefault();
    dispatch(updateVncSessionLogStatus({ swipeableOpened: expanded }));
  };

  return (
    <SwipeableDrawer
      anchor={swipeableAnchor}
      open={swipeableOpened}
      onOpen={handleToggle(true)}
      onClose={handleToggle(false)}>
      <SwipeableTermninalLogs />
    </SwipeableDrawer>
  );
};
