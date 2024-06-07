import React from 'react';
import { useNavigate } from 'react-router-dom';
import { IconButton, Tooltip } from '@mui/material';
import ExitToAppOutlinedIcon from '@mui/icons-material/ExitToAppOutlined';
import { useAppDispatch } from '@hooks/useAppDispatch';
import { deleteVncSession } from '@redux/VncSession/VncSessionSlice';
import styles from './ExitButton.module.scss';

export const ExitButton = () => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();

  const handleOnClick = (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    e.preventDefault();

    dispatch(deleteVncSession());
    navigate('/', { replace: true });
  };

  return (
    <Tooltip title="Close Session">
      <IconButton onClick={handleOnClick} className={styles.btnClose}>
        <ExitToAppOutlinedIcon />
      </IconButton>
    </Tooltip>
  );
};
