import { VncControl } from '../Controls/VncControl';
import styles from './ActiveVncPanel.module.scss';

export const ActiveVncPanel = () => {
  return (
    <div className={styles.vncPanel}>
      {/* <Stack spacing={2}> */}
      {/* <ExitButton /> */}
      {/* <LogsButton /> */}
      {/* <Divider /> */}
      {/* <TerminateButton /> */}
      {/* </Stack> */}
      <VncControl isExpand={false} />
      {/* <SwipeableControl /> */}
    </div>
  );
};
