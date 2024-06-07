import { type ReactNode, useEffect } from 'react';
import { fetchBrowsers } from '@redux/browsers/browsersThunk';
import { useAppDispatch } from '@hooks/useAppDispatch';

interface BrowsersLoaderProps {
  children: ReactNode;
}

export const BrowsersLoader = ({ children }: BrowsersLoaderProps) => {
  const dispatch = useAppDispatch();

  useEffect(() => {
    void dispatch(fetchBrowsers());
  }, [dispatch]);

  return <>{children}</>;
};
