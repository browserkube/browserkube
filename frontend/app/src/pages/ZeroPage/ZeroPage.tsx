import { useDispatch } from 'react-redux';
import { Button } from '@reportportal/ui-kit';
import { useNavigate } from 'react-router-dom';
import { useEffect } from 'react';
import { lang } from '@app/constants';
import { MODAL_TYPE } from '@shared/types/UI';
import { openModal } from '@redux/UI/UISlice';
import { BrowserKubeIcon } from '../../shared/icons/BrowserKubeIcon';

import styles from './ZeroPage.module.scss';
import '@app/theme/fonts.css';

const automationBtn = {
  backgroundColor: 'white',
  color: '#00829B',
  '&:hover': {
    backgroundColor: '#F7F7F8',
  },
} as const;

const { textLine2, textLine3 } = lang.zeroPage;

export const ZeroPage = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const handleOpen = (mode: string) => {
    const modalType = mode === 'manual' ? MODAL_TYPE.CREATE_SESSION : MODAL_TYPE.GENERATE_CODE_SNIPPET;
    dispatch(openModal({ id: 'createSession', component: modalType }));
  };

  useEffect(() => {
    const isActiveUser = !localStorage.getItem('isActiveUser');

    if (isActiveUser) {
      localStorage.setItem('isActiveUser', 'true');
    } else {
      navigate('/live-sessions/');
    }
  }, []);

  return (
    <>
      <div className={styles.container}>
        <BrowserKubeIcon />
        <div className={styles.project_title}>
          <div className={styles.text_browser}>Browser</div>
          <div className={styles.text_cube}>Kube</div>
        </div>
        <div className={styles.text_line2}>{textLine2}</div>
        <div className={styles.text_line3}>{textLine3}</div>
        <div className={styles.btn_container}>
          <Button
            style={automationBtn}
            onClick={() => {
              handleOpen('auto');
            }}>
            Generate for Automation
          </Button>
          <Button
            onClick={() => {
              handleOpen('manual');
            }}>
            Create Manual session
          </Button>
        </div>
      </div>
    </>
  );
};
