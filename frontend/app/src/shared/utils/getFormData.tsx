import { type Browser } from '@shared/types/browsers';
import { type FormData } from '@shared/types/createSession';
import { getStringFormated } from './getStringFormated';

export const getFormData = (data: Browser[]): FormData => {
  return data
    .filter(({ name, type }) => type === 'WEBDRIVER')
    .reduce((acc: FormData, { platformName, name, version, image, resolutions }) => {
      const platformIndex = acc.findIndex((platform) => platform.value === platformName);

      if (platformIndex === -1) {
        // Platform not found, create a new platform entry
        const newPlatform = {
          label: getStringFormated(platformName),
          value: platformName,
          browsers: [
            {
              label: getStringFormated(name),
              value: name,
              versions: [
                {
                  value: version,
                  image,
                  resolutions,
                },
              ],
            },
          ],
        };

        return [...acc, newPlatform];
      } else {
        // Platform found, check if browser exists
        const browserIndex = acc[platformIndex].browsers.findIndex((browser) => browser.value === name);

        if (browserIndex === -1) {
          // Browser not found, create a new browser entry
          const newBrowser = {
            label: getStringFormated(name),
            value: name,
            versions: [
              {
                value: version,
                image,
                resolutions,
              },
            ],
          };

          acc[platformIndex].browsers.push(newBrowser);
        } else {
          // Browser found, check if version exists
          const versionIndex = acc[platformIndex].browsers[browserIndex].versions.findIndex(
            (versionItem) => versionItem.value === version
          );

          if (versionIndex === -1) {
            // Version not found, create a new version entry
            const newVersion = {
              value: version,
              image,
              resolutions,
            };

            acc[platformIndex].browsers[browserIndex].versions.push(newVersion);
          }
        }

        return acc;
      }
    }, []);
};
