import React from 'react';
import { IconButton, Tooltip } from '@mui/material';
import { useSelector } from 'react-redux';
import { useAppDispatch } from '@hooks/useAppDispatch';
import { getVncSessionScreenStatus } from '@redux/VncSession/VncSessionSelectors';
import { updateVncSessionScreenStatus } from '@redux/VncSession/VncSessionSlice';
import { LockerIcon } from '@shared/icons/lockerIcon';
import { UnLockerIcon } from '@shared/icons/unLockerIcon';

export const LockButton = () => {
  const dispatch = useAppDispatch();
  const isLocked = useSelector((state) => getVncSessionScreenStatus(state));

  const handleOnClick = (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    e.preventDefault();
    dispatch(updateVncSessionScreenStatus({ VNCScreenStatus: !isLocked }));
  };
  return (
    <Tooltip title={isLocked ? 'Unlock VNC' : 'Lock VNC'}>
      <IconButton onClick={handleOnClick}>{isLocked ? <LockerIcon /> : <UnLockerIcon />}</IconButton>
    </Tooltip>
  );
};
