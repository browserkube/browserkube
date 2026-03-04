import { useEffect, useRef } from 'react';
import { type TerminalOutput } from '@components/Logs/useTerminal';

const RECONNECT_DELAY_MS = 2000;
const MAX_RECONNECT_ATTEMPTS = 5;

interface UseWebSocketOptions {
  wsUrl: string;
  xterm: TerminalOutput | null;
}

export const useWebSocketForTerminal = ({ wsUrl, xterm }: UseWebSocketOptions): void => {
  const decoder = useRef(new TextDecoder('utf8'));
  const socketRef = useRef<WebSocket | null>(null);
  const attemptsRef = useRef(0);
  const cancelledRef = useRef(false);

  useEffect(() => {
    if (!xterm) return;

    cancelledRef.current = false;
    attemptsRef.current = 0;

    function connect() {
      if (cancelledRef.current) return;

      const ws = new WebSocket(wsUrl);
      ws.binaryType = 'arraybuffer';
      socketRef.current = ws;

      ws.addEventListener('open', () => {
        attemptsRef.current = 0;
        xterm!.terminal.clear();
        xterm!.terminal.writeln('Session Logs...');
      });

      ws.addEventListener('message', (e: MessageEvent) => {
        xterm!.terminal.writeln(decoder.current.decode(e.data));
      });

      ws.addEventListener('close', (e: CloseEvent) => {
        if (cancelledRef.current) return;
        if (attemptsRef.current < MAX_RECONNECT_ATTEMPTS) {
          attemptsRef.current++;
          xterm!.terminal.writeln(`\r\nConnection lost. Reconnecting (${attemptsRef.current}/${MAX_RECONNECT_ATTEMPTS})...`);
          setTimeout(connect, RECONNECT_DELAY_MS);
        } else {
          xterm!.terminal.writeln('\r\nCould not reconnect to session logs.');
        }
      });
    }

    connect();

    return () => {
      cancelledRef.current = true;
      if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
        socketRef.current.close();
      }
    };
  }, [xterm, wsUrl]);
};
