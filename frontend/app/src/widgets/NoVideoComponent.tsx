import { Button } from '@reportportal/ui-kit';
import { lang } from '@app/constants';
import { type AttachmentsTabProps } from '@shared/types/UI';
import { NoVideoIcon } from '@shared/icons/noVideoImage';
import { LeaveIcon } from '@shared/icons/leaveIcon';

const frameStyle = {
  display: 'flex',
  paddingTop: '48px',
  width: '477px',
  flexDirection: 'column',
  alignItems: 'center',
} as const;

const noContentContainer = {
  display: 'flex',
  flex: 1,
  alignItems: 'center',
  justifyContent: 'center',
  marginTop: '12px',
  paddingBottom: '32px',
} as const;

const textInfoStyle = {
  fontSize: '11px',
  textAlign: 'center',
  marginTop: '4px',
} as const;

const leaveBtn = {
  color: '#00829B',
  backgroundColor: 'white',
  marginTop: '32px',
  padding: '7px 16px',
  border: '1px solid #00829B',
} as const;

const textHeader = {
  fontFamily: 'Open Sans',
  marginTop: '48px',
  fontSize: '20px',
} as const;

const { headerText, infoText } = lang.attachments;

export const NoVideoComponent = (props: AttachmentsTabProps) => {
  const { isTerminated, setTab } = props;

  const scrollToTop = () => {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleLeave = () => {
    setTab(isTerminated ? 'details' : 'liveSession');
    scrollToTop();
  };

  return (
    <div style={noContentContainer}>
      <div style={frameStyle}>
        <NoVideoIcon />
        <div style={textHeader}>{headerText}</div>
        <div style={textInfoStyle}>{infoText}</div>
        <Button onClick={handleLeave} style={leaveBtn} icon={LeaveIcon('#00758C')} variant="text">
          {isTerminated ? 'Details' : 'Live Session'}
        </Button>
      </div>
    </div>
  );
};
