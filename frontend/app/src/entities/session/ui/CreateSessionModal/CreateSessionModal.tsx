import { Modal } from '@reportportal/ui-kit';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { type FormValues } from '@shared/types/createSession';
import { type ModalProps } from '@shared/types/UI';
import { useAppDispatch } from '@hooks/useAppDispatch';
import { closeModal } from '@redux/UI/UISlice';
import { CreateManualSession } from '@api/sessions/getSessions';
import { type CreateSessionResponse } from '@shared/types/webDriver';
import { saveActiveSessionId } from '@redux/sessionDetails/sessionSlice';
import { fetchSessions } from '@redux/sessions/sessionsThunk';
import { CreateSessionModalForm } from './components/CreateSessionModalForm/CreateSessionModalForm';

const INITIAL_FORM_STATE = {
  sessionName: `Session# ${Math.floor(Math.random() * 100)}`,
  platformName: 'linux',
  browserName: 'chrome',
  browserVersion: '',
  screenResolution: '',
  recordVideo: true,
};

export const CreateSessionModal = ({ id, zIndex }: ModalProps) => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const [createSessionFormValues, setCreateSessionFromValues] = useState<FormValues>(INITIAL_FORM_STATE);

  const formChangeHandler = (value: Partial<FormValues>) => {
    const dropDownKey = Object.keys(value)[0];
    if (dropDownKey === 'browserName') {
      setCreateSessionFromValues((prevState) => ({ ...prevState, browserVersion: '', screenResolution: '' }));
    }
    setCreateSessionFromValues((prevState) => ({ ...prevState, ...value }));
  };

  const closeModalHandler = () => {
    dispatch(closeModal(id));
  };

  const createNewSession = async () => {
    navigate('/live-sessions/');
    closeModalHandler();
    try {
      const { value }: CreateSessionResponse = await CreateManualSession(createSessionFormValues);
      dispatch(saveActiveSessionId({ id: value.sessionId }));
      void dispatch(fetchSessions());
    } catch (error) {
      console.log('Error in createNewSession, fill the correct data', error);
    }
  };

  return (
    <Modal
      title="Create Manual Session"
      onClose={closeModalHandler}
      okButton={{
        children: 'Create',
        // TODO:fill method in future
        // eslint-disable-next-line @typescript-eslint/no-misused-promises
        onClick: createNewSession,
      }}
      cancelButton={{
        children: 'Cancel',
        onClick: closeModalHandler,
      }}>
      <CreateSessionModalForm formValues={createSessionFormValues} onChange={formChangeHandler} />
      {/* <div className={styles.generate_link_container}>
        <Button variant="text">Generate Code Snippet for Auto Session</Button>
      </div> */}
    </Modal>
  );
};
