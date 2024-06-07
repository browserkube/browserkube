import { type Modal, MODAL_TYPE } from '@shared/types/UI';
import { CreateSessionModal } from '../../entities/session/ui';

export const ModalComponent = ({ id, component, zIndex }: Modal) => {
  if (component === MODAL_TYPE.CREATE_SESSION) {
    return <CreateSessionModal id={id} zIndex={zIndex} />;
  }
  return null;
};
