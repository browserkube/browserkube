import styles from '@app/styles/details.module.scss';
import { getStringFormated } from '@shared/utils/getStringFormated';

import { type Session } from '@shared/types/sessions';

const detailsTextStyle = {
  fontSize: '13px',
  lineHeight: '20px',
  padding: '0 0 24px 32px',
} as const;

const titleStyle = {
  marginRight: '8px',
} as const;

export const DetailsInfo = ({ session }: { session: Session }) => {
  const { name, platformName, browser, browserVersion, type, createdAt, manual } = session;

  return (
    <div style={detailsTextStyle}>
      <div className={styles.detail_line}>
        <div style={titleStyle}>Name:</div>
        <div>{name}</div>
      </div>
      <div className={styles.detail_line}>
        <div style={titleStyle}>Platform:</div>
        <div>{getStringFormated(platformName)}</div>
      </div>
      <div className={styles.detail_line}>
        <div style={titleStyle}>Browser:</div>
        <div>{`${getStringFormated(browser)} ${browserVersion}`}</div>
      </div>
      <div className={styles.detail_line}>
        <div style={titleStyle}>Session Type:</div>
        <div>{manual ? 'Manual' : 'Auto'}</div>
      </div>
      <div className={styles.detail_line}>
        <div style={titleStyle}>Protocol Type:</div>
        <div>{getStringFormated(type)}</div>
      </div>
      <div className={styles.detail_line}>
        <div style={titleStyle}>Started at:</div>
        <div>{new Date(createdAt).toLocaleString().replaceAll('.', '/')}</div>
      </div>
      <div className={styles.detail_line}>
        <div style={titleStyle}>Completed at:</div>
        <div>no info yet</div>
      </div>
      <div className={styles.detail_line}>
        <div style={titleStyle}>Duration at:</div>
        <div>no info yet</div>
      </div>
    </div>
  );
};
