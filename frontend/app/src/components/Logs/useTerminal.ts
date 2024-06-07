import { useEffect, useState } from 'react';
import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { SearchAddon } from 'xterm-addon-search';

export interface TerminalOutput {
  terminal: Terminal;
  addons: {
    fitAddon: FitAddon;
    searchAddon: SearchAddon;
  };
}

function createTerminal(): TerminalOutput {
  const terminal = new Terminal({
    allowProposedApi: true,
    disableStdin: true,
    cursorBlink: true,
    cursorStyle: 'block',
    convertEol: true,
    fontFamily: `'Fira Mono', monospace`,
    fontSize: 16,
  });
  const fitAddon = new FitAddon();
  const searchAddon = new SearchAddon();
  terminal.loadAddon(fitAddon);
  terminal.loadAddon(searchAddon);
  return { addons: { fitAddon, searchAddon }, terminal };
}

interface UseTerminalOptions {
  container: { current: HTMLDivElement | null };
}

export const useTerminal = ({ container }: UseTerminalOptions) => {
  const [xterm, setXterm] = useState<TerminalOutput | null>(null);

  useEffect(() => {
    if (!container?.current) {
      return;
    }
    const element = container.current;
    const newXterm = createTerminal();
    setXterm(newXterm);
    newXterm.terminal.open(element);

    const ro = new ResizeObserver(() => {
      newXterm.addons.fitAddon.fit();
    });
    ro.observe(element);

    return () => {
      ro.unobserve(element);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return [xterm];
};
