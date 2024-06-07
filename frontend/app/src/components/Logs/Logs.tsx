import React, { useRef } from 'react';
import cx from 'classnames';
import { useTerminal } from '@components/Logs/useTerminal';
import { useWebSocketForTerminal } from '@components/Logs/useWebSocketForTerminal';
import styles from './Logs.module.scss';

export interface LogsProps {
  logsUrl: string;
  isOpened: boolean;
}

export const Logs = ({ logsUrl, isOpened }: LogsProps) => {
  const terminalContainer = useRef<HTMLDivElement | null>(null);
  const [xterm] = useTerminal({ container: terminalContainer });
  useWebSocketForTerminal({ wsUrl: logsUrl, xterm });

  const handleSearchInput = (event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const data = event.currentTarget.value;

    if (data.length > 0) {
      xterm?.addons.fitAddon.fit();
      xterm?.addons.searchAddon.findNext(event.currentTarget.value, {
        wholeWord: false,
        caseSensitive: false,
        regex: true,
        incremental: false,
        decorations: {
          matchBackground: '#37c971',
          activeMatchBackground: '#37c971',
          matchOverviewRuler: '#37c971',
          activeMatchColorOverviewRuler: '#37c971',
        },
      });
    } else {
      xterm?.addons.searchAddon.clearDecorations();
    }
  };

  return (
    <div className={styles.logsTerminal}>
      <input
        type="text"
        placeholder="Search logs"
        className={styles.logsTerminal__search}
        onChange={handleSearchInput}
      />
      <div className={cx(styles.logsTerminal__container, { [styles.opened]: isOpened })} ref={terminalContainer} />
    </div>
  );
};
