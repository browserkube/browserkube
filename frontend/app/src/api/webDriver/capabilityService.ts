// TODO: tmp factory - should be replaced
export const createCapabilities = (browser: string, version: string) => {
  const extOptions = ((browser) => {
    switch (browser) {
      case 'chrome':
        return {
          'goog:chromeOptions': {
            args: ['start-maximized'],
          },
        };

      case 'firefox':
        return {
          'moz:firefoxOptions': {
            prefs: {
              'browser.fullscreen.animateUp': 0,
              'browser.fullscreen.autohide': false,
            },
          },
        };

      default:
        return {};
    }
  })(browser);
  return {
    desiredCapabilities: {
      browserName: browser,
      'browserkube:options': {
        enableVideo: true,
      },
      browserVersion: version,
      enableVNC: true,
      saveVideoEndpoint: 'file:///home/seluser/videos',
      labels: { manual: 'true' },
      sessionTimeout: '60m',
      name: 'Manual session',
      // stub param for testing custom name of session
      // 'browserkube:options': {
      //   name: `custom name ${new Date().getTime()}`,
      // },
    },
    capabilities: {
      alwaysMatch: {
        browserName: browser,
        browserVersion: version,
        'browserkube:options': {
          enableVideo: true,
        },
        'selenoid:options': {
          enableVNC: true,
          sessionTimeout: '60m',
          saveVideoEndpoint: 'file:///home/seluser/videos',
          labels: { manual: 'true' },
          screenResolution: '1920x1080x24',
        },
        ...extOptions,
      },
      firstMatch: [{}],
    },
  };
};
