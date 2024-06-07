import { Divider, DividerType } from 'widgets/Divider';
import { CurrentSession } from '../../widgets/CurrentSession';
import { SessionsBlock } from '../../widgets/SessionsBlock';
import { ChosenSessionProvider } from './SessionContext';

const container = {
  display: 'flex',
  width: '1440px',
  height: '1024px',
} as const;

export const LiveSessions = () => {
  return (
    <>
      <div style={container}>
        <ChosenSessionProvider>
          <SessionsBlock />
          <Divider type={DividerType.VERTICAL} />
          <CurrentSession />
        </ChosenSessionProvider>
      </div>
    </>
  );
};
