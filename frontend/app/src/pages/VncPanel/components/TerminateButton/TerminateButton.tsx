import React from 'react';
import { useNavigate } from 'react-router-dom';
import { IconButton, Tooltip } from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import { useSelector } from 'react-redux';
import { useAppDispatch } from '@hooks/useAppDispatch';
import { deleteVncSession } from '@redux/VncSession/VncSessionSlice';
import { deleteWdSession } from '@redux/sessions/sessionsThunk';
import { getVncSessionId } from '@redux/VncSession/VncSessionSelectors';
import styles from './TerminateButton.module.scss';

export const TerminateButton = () => {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const sessionId = useSelector((state) => getVncSessionId(state));

  const handleOnClick = (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    e.preventDefault();

    dispatch(deleteVncSession());
    void dispatch(deleteWdSession({ sessionId }));
    navigate('/', { replace: true });
  };

  return (
    <Tooltip title="Terminate Session">
      <IconButton key="btn-terminate" onClick={handleOnClick} className={styles.btnTerminate}>
        <DeleteIcon />
      </IconButton>
    </Tooltip>
  );
};
