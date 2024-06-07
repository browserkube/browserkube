import { useEffect, useRef } from 'react';
import { type TerminalOutput } from '@components/Logs/useTerminal';

interface SocketOptions {
  wsUrl: string;
  onOpen: (e: Event) => void;
  onMessage: (e: MessageEvent) => void;
  onClose: (e: CloseEvent) => void;
}

function createSocket(options: SocketOptions) {
  const { wsUrl, onOpen, onMessage, onClose } = options;
  const ws = new WebSocket(wsUrl);
  ws.binaryType = 'arraybuffer';
  ws.addEventListener('open', onOpen);
  ws.addEventListener('message', onMessage);
  ws.addEventListener('close', onClose);
  return ws;
}

interface UseWebSocketOptions {
  wsUrl: string;
  xterm: TerminalOutput | null;
}

export const useWebSocketForTerminal = ({ wsUrl, xterm }: UseWebSocketOptions): void => {
  const decoder = useRef(new TextDecoder('utf8'));

  useEffect(() => {
    if (xterm) {
      const wsSocket = createSocket({
        wsUrl,
        onOpen,
        onMessage,
        onClose,
      });

      return () => {
        if (wsSocket.readyState === 1) {
          wsSocket.removeEventListener('open', onOpen);
          wsSocket.removeEventListener('message', onMessage);
          wsSocket.removeEventListener('close', onClose);
          wsSocket.close();
        }
      };
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [xterm]);

  function onOpen() {
    if (xterm) {
      xterm.terminal.clear();
      xterm.terminal.writeln('Session Logs...');
    }
  }

  function onMessage(e: MessageEvent) {
    if (xterm) {
      xterm.terminal.writeln(decoder.current.decode(e.data));
    }
  }

  function onClose(e: CloseEvent) {
    // TODO: add reconnection logic
  }
};
