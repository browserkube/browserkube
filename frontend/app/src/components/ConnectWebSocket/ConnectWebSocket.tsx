import { type ReactNode, useEffect } from 'react';
import { BASE_URL } from '@shared/lib';
import { connect, disconnect } from '@redux/webSocket/webSocketSlice';
import { useAppDispatch } from '@hooks/useAppDispatch';

interface ConnectWebSocketProps {
  children: ReactNode;
}

export const ConnectWebSocket = ({ children }: ConnectWebSocketProps) => {
  const isHttps = window.location.protocol === 'https:';
  const prefix = isHttps ? 'wss' : 'ws';
  const url = `${prefix}://${window.location.hostname}${BASE_URL}/events`;
  const dispatch = useAppDispatch();

  useEffect(() => {
    dispatch(connect({ url }));
    return () => {
      dispatch(disconnect());
    };
  }, [dispatch, url]);

  return <>{children}</>;
};
