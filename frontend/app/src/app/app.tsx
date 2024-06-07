import { HashRouter, Route, Routes } from 'react-router-dom';
import { createRoot } from 'react-dom/client';
import { Provider } from 'react-redux';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { toast, ToastContainer } from 'react-toastify';
import { defaultTheme } from 'app/theme/defaultTheme';
import { PATHS } from '@shared/routing';
import { VncPanel } from '@pages/VncPanel/VncPanel';
import { appStore } from '@redux/store';
import { ConnectWebSocket } from '@components/ConnectWebSocket/ConnectWebSocket';
import { BrowsersLoader } from '@components/BrowsersLoader/BrowsersLoader';
import { ModalFactory } from '@shared/ModalFactory/ModalFactory';
import { ZeroPage } from '@pages/ZeroPage/ZeroPage';
import { LiveSessions } from '@pages/LiveSessions/LiveSessions';
import 'react-toastify/ReactToastify.min.css';
import '@reportportal/ui-kit/dist/style.css';
import '@app/styles/globals.scss';

const rootElement = document.getElementById('root');

createRoot(rootElement!).render(
  <HashRouter basename="/">
    <Provider store={appStore}>
      <ThemeProvider theme={defaultTheme}>
        <CssBaseline />
        <ModalFactory />
        <ConnectWebSocket>
          <BrowsersLoader>
            <Routes>
              <Route path={PATHS.HOME} element={<ZeroPage />} />
              <Route path={PATHS.LIVE_SESSIONS} element={<LiveSessions />} />
              <Route path={PATHS.ACTIVE_SESSION_DETAILS} element={<VncPanel />} />
            </Routes>
            <ToastContainer
              position={toast.POSITION.BOTTOM_LEFT}
              autoClose={Number(process.env.REACT_APP_HIDE_EACH_TOAST_TIMEOUT)}
              hideProgressBar={false}
              newestOnTop={false}
              closeOnClick
              rtl={false}
              pauseOnFocusLoss
              pauseOnHover
              theme="colored"
            />
          </BrowsersLoader>
        </ConnectWebSocket>
      </ThemeProvider>
    </Provider>
  </HashRouter>
);
