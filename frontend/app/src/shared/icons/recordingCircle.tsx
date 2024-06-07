import styles from '@pages/LiveSessions/LiveSession.module.scss';

export const RecordingIcon = ({ isAnimation }: { isAnimation: boolean }) => {
  return (
    <svg
      className={isAnimation ? styles['pulse-animation'] : ''}
      xmlns="http://www.w3.org/2000/svg"
      width="20"
      height="20"
      viewBox="0 0 20 20"
      fill="none">
      <circle cx="10" cy="10" r="8" fill="#D32F2F" />
    </svg>
  );
};
