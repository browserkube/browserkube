import { test, expect , firefox} from '@playwright/test';
import dotenv from "dotenv";

test('homepage has Playwright in title and get started link linking to the intro page', async () => {
  try {

    dotenv.config();
    const wsUrl = process.env.HUB_URL
    const browser = await firefox.connect({ timeout: 0, wsEndpoint: wsUrl });

    // const browserServer = await chromium.launchServer();
    // const wsEndpoint = browserServer.wsEndpoint();
    // console.log(wsEndpoint)
    // const browser = await chromium.connect({ timeout: 0, wsEndpoint: 'ws://localhost:4444' });
    const page = await browser.newPage();

    await page.goto('https://playwright.dev/');

     // Expect a title "to contain" a substring.
    await expect.soft(page).toHaveTitle(/.*Playwright/);

    // create a locator
    const getStarted = page.locator('text=Get Started');
    
    // Expect an attribute "to be strictly equal" to the value.
    await expect.soft(getStarted).toHaveAttribute('href', '/docs/intro');

    // Expects the URL to contain intro.
    await expect.soft(page).toHaveURL(/.*intro/);

  } catch (err) {
    console.log(err);
  }

});
