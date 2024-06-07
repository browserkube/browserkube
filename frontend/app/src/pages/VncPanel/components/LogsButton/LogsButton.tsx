import React from 'react';
import { IconButton, Tooltip } from '@mui/material';
import TerminalSharpIcon from '@mui/icons-material/TerminalSharp';
import { useAppDispatch } from '@hooks/useAppDispatch';
import { updateVncSessionLogStatus } from '@redux/VncSession/VncSessionSlice';
import styles from './LogsButton.module.scss';

export const LogsButton = () => {
  const dispatch = useAppDispatch();

  const handleOnClick = (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    e.preventDefault();
    dispatch(updateVncSessionLogStatus({ swipeableOpened: true }));
  };

  return (
    <Tooltip key="btn-terminate232" title="Logs">
      <IconButton key="btn-terminate11" onClick={handleOnClick} className={styles.btnLogs}>
        <TerminalSharpIcon />
      </IconButton>
    </Tooltip>
  );
};
