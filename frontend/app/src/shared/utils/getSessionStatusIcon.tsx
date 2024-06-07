/**
 *
 * @param status
 * @returns icon as JSX.element according session status
 */

import { SessionActive } from '@shared/icons/SessionActive';
import { SessionPending } from '@shared/icons/SessionPending';
import { SessionTerminated } from '@shared/icons/SessionTerminated';

export const getSessionStatusIcon = (status: string) => {
  switch (status) {
    case 'Terminated': {
      return <SessionTerminated />;
    }
    case 'running': {
      return <SessionActive />;
    }
    default:
      return <SessionPending />;
  }
};
