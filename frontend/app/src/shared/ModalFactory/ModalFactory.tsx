import { useSelector } from 'react-redux';
import { getModals } from '@redux/UI/UISelectors';
import { ModalComponent } from './ModalComponent';

export const ModalFactory = () => {
  const modals = useSelector(getModals);

  return (
    <>
      {modals.map((modal) => (
        <ModalComponent {...modal} key={modal.id} />
      ))}
    </>
  );
};
