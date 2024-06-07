const puppeteer = require("puppeteer");
const browserkubeTools = require("./browserkube_helpers.js");

 // Get a field's value
 (async () => {
    var resp = await browserkubeTools.startBrowser();
    console.log(resp);
    const sessionID = resp.value.sessionId

    try {
        const browser = await puppeteer.connect({
            browserWSEndpoint: resp["value"]["capabilities"]["se:cdp"]
        });
        const page = await browser.newPage()
        await page.goto('https://news.ycombinator.com/news')
        const name = await page.$eval('.hnname > a', el => el.innerText)
        console.log(name);
        await browser.close();
        await browserkubeTools.stopBrowser(sessionID);
    } catch (e) {
        console.error(e);
    }
})();


// Mouse click
(async () => {
    var resp = await browserkubeTools.startBrowser();
    console.log(resp);
    const sessionID = resp.value.sessionId
    try {
        const browser = await puppeteer.connect({
            browserWSEndpoint: resp["value"]["capabilities"]["se:cdp"]
        });
        const page = await browser.newPage()

        // set the viewport so we know the dimensions of the screen
        await page.setViewport({ width: 800, height: 600 })

        // go to a page setup for mouse event tracking
        await page.goto('http://unixpapa.com/js/testmouse.html')

        // click an area
        await page.mouse.click(132, 103, { button: 'left' })
        await browser.close();
        await browserkubeTools.stopBrowser(sessionID);
    } catch (e) {
        console.log(e);
    }
})();

// Create a pdf
(async () => {
    var resp = await browserkubeTools.startBrowser();
    console.log(resp);
    const sessionID = resp.value.sessionId
    try {
        const browser = await puppeteer.connect({
            browserWSEndpoint: resp["value"]["capabilities"]["se:cdp"]
        });
        const page = await browser.newPage()
        // 1. Create PDF from URL
        await page.goto('https://github.com/GoogleChrome/puppeteer/blob/master/docs/api.md#pdf')
        await page.pdf({ path: 'api.pdf', format: 'A4' })

        // 2. Create PDF from static HTML
        const htmlContent = `<body>
  <h1>An example static HTML to PDF</h1>
  </body>`
        await page.setContent(htmlContent)
        await page.pdf({ path: 'html.pdf', format: 'A4' })
        await browser.close();
        await browserkubeTools.stopBrowser(sessionID);
    } catch (e) {
        console.log(e)
    }
})();