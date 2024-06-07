import { ChromeIcon } from '@shared/chips/chrome';
import { EdgeIcon } from '@shared/chips/edge';
import { FirefoxIcon } from '@shared/chips/firefox';

export const getBrowserChip = (browser: string) => {
  let browserIcon;
  const browserArray: string[] = ['firefox', 'edge', 'chrome'];
  const browserCheck = browserArray.find((item) => browser.includes(item));

  switch (browserCheck) {
    case 'firefox':
      browserIcon = <FirefoxIcon />;
      break;
    case 'edge':
      browserIcon = <EdgeIcon />;
      break;
    case 'chrome':
      browserIcon = <ChromeIcon />;
      break;
    default:
      browserIcon = <ChromeIcon />;
      break;
  }

  return browserIcon;
};
