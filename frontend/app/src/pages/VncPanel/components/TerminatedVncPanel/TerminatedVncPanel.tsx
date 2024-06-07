import { Stack } from '@mui/material';
import Typography from '@mui/material/Typography';
import { ExitButton } from '../ExitButton/ExitButton';
import styles from './TerminatedVncPanel.module.scss';

export const TerminatedVncPanel = () => {
  return (
    <div className={styles.vncPanelTerminated}>
      <Stack spacing={2}>
        <ExitButton />
      </Stack>
      <Typography variant="h3" component="h1">
        Session is terminated.
      </Typography>
    </div>
  );
};
