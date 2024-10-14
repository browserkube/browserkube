package test

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/reportportal/goRP/v5/pkg/gorp"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tebeka/selenium"
)

// Chrome Store Extension ID for "metamask" addon
const (
	extensionID   = "nkbihfbeogaeaoehlefnkodbefgpgknn"
	chromeVersion = "124.0"
)

func defaultCaps() selenium.Capabilities {
	return selenium.Capabilities{
		"browserName":    "chrome",
		"browserVersion": fmt.Sprintf("%s-selenoid", chromeVersion),
		"browserkube:options": map[string]interface{}{
			"enableVNC": true,
		},
	}
}

func TestWebDriver(t *testing.T) {
	suite.Run(t, new(WebDriverTestsSuite))
}

type WebDriverTestsSuite struct {
	suite.Suite
	wdURL    string
	rpClient *gorp.Client
	launchID string
}

func (suite *WebDriverTestsSuite) TearDownTest() {
	if suite.rpClient == nil {
		return
	}
	_, err := suite.rpClient.FinishLaunch(suite.launchID, &gorp.FinishExecutionRQ{
		EndTime: gorp.NewTimestamp(time.Now()),
	},
	)
	require.NoError(suite.T(), err)
}

func (suite *WebDriverTestsSuite) SetupTest() {
	selenium.SetDebug(true)
	file, err := os.Open(".env")
	if err == nil {
		lErr := loadEnvFile(file)
		require.NoError(suite.T(), lErr, "Unable to load env file")
	}
	suite.wdURL = os.Getenv("WD_URL")
	require.NotEmpty(suite.T(), suite.wdURL, "Webdriver URL isn't provided. Use 'WD_URL' env variable")

	rpEndpoint := os.Getenv("RP_ENDPOINT")
	rpProject := os.Getenv("RP_PROJECT")
	rpUuid := os.Getenv("RP_UUID")

	if rpEndpoint != "" && rpProject != "" && rpUuid != "" {
		suite.rpClient = gorp.NewClient(rpEndpoint, rpProject, rpUuid)
		u := uuid.New()
		launch, err := suite.rpClient.StartLaunch(&gorp.StartLaunchRQ{
			StartRQ: gorp.StartRQ{
				Name:      "Browserkube Test",
				UUID:      &u,
				StartTime: gorp.NewTimestamp(time.Now()),
			},
			Mode: gorp.LaunchModes.Default,
		})
		require.NoError(suite.T(), err)
		suite.launchID = launch.ID
	}

	caps := defaultCaps()
	caps["browserkube:options"].(map[string]interface{})["reportportal"] = map[string]string{
		"launchId": suite.launchID,
		"project":  rpProject,
	}
}

func (suite *WebDriverTestsSuite) TestBasicSelenium() {
	caps := defaultCaps()
	caps["browserVersion"] = fmt.Sprintf("%s-selenium", chromeVersion)
	caps["browserkube:options"].(map[string]interface{})["name"] = fmt.Sprintf("test selenium %s", caps["browserVersion"])

	suite.testBasic(caps)
}

func (suite *WebDriverTestsSuite) TestBasicSelenoid() {
	caps := defaultCaps()
	caps["browserVersion"] = fmt.Sprintf("%s-selenoid", chromeVersion)
	caps["browserkube:options"].(map[string]interface{})["name"] = fmt.Sprintf("test selenoid %s", caps["browserVersion"])

	suite.testBasic(caps)
}

func (suite *WebDriverTestsSuite) TestBasicWithVideo() {
	caps := defaultCaps()
	caps["browserVersion"] = fmt.Sprintf("%s-selenoid", chromeVersion)

	bkOpts := caps["browserkube:options"].(map[string]interface{})
	bkOpts["enableVideo"] = true
	bkOpts["name"] = fmt.Sprintf("Test with Video %s", caps["browserVersion"])

	sessionID := suite.testBasic(caps)
	baseURL, err := url.Parse(suite.wdURL)
	require.NoError(suite.T(), err)
	baseURL.Path = fmt.Sprintf("/browserkube/session/files/%s/video.mp4", sessionID)

	client := http.DefaultClient
	resp, err := client.Get(baseURL.String())
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	require.Equal(suite.T(), "video/mp4", resp.Header.Get("Content-Type"))
	require.NotEmpty(suite.T(), resp.Header.Get("Content-Length"))
}

func (suite *WebDriverTestsSuite) TestDownload() {
	wd, err := selenium.NewRemote(defaultCaps(), suite.wdURL)
	require.NoError(suite.T(), err)
	defer wd.Quit()

	err = wd.Get("https://www.browserstack.com/test-on-the-right-mobile-devices")
	require.NoError(suite.T(), err)

	cookinEl, err := wd.FindElement(selenium.ByID, "accept-cookie-notification")
	require.NoError(suite.T(), err)
	require.NoError(suite.T(), cookinEl.Click())

	downloadEl, err := wd.FindElement(selenium.ByXPATH, "//a[text()='CSV']")
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), downloadEl)
	require.NoError(suite.T(), downloadEl.Click())

	_, err = wd.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{downloadEl})
	require.NoError(suite.T(), err)

	_, err = wd.ExecuteScript("window.scrollBy(0,-100)", []interface{}{})
	require.NoError(suite.T(), err)
	time.Sleep(1 * time.Second)

	fileResp, err := http.Get(fmt.Sprintf(`%s/session/%s/browserkube/downloads/%s`, suite.wdURL, wd.SessionID(), "BrowserStack - List of devices to test on.csv"))

	require.NoError(suite.T(), err)

	defer fileResp.Body.Close()
	require.True(suite.T(), fileResp.StatusCode < http.StatusBadRequest, "Status Code: %s", fileResp.StatusCode)

	content, err := io.ReadAll(fileResp.Body)
	require.NoError(suite.T(), err)

	fmt.Println(string(content))
}

// Prerequisites: Start an image with option plugin and goog:chromeOptions with load-extension arg inside capabilities
// For reference: /docs/userguide.md
func (suite *WebDriverTestsSuite) TestPlugin() {
	caps := defaultCaps()
	caps["browserkube:options"].(map[string]interface{})["extensions"] = []map[string]interface{}{
		{
			"extensionId": "nkbihfbeogaeaoehlefnkodbefgpgknn",
		},
	}
	caps["goog:chromeOptions"] = map[string]interface{}{
		"args": []string{"load-extension=/opt/google/chrome/extensions/" + extensionID},
	}

	// check caps definition
	wd, err := selenium.NewRemote(caps, suite.wdURL)
	require.NoError(suite.T(), err)
	defer wd.Quit()

	err = wd.Get("chrome://extensions/?id=" + extensionID)
	require.NoError(suite.T(), err)
	<-time.After(30 * time.Second)
}

func loadEnvFile(f *os.File) error {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		if err := os.Setenv(parts[0], parts[1]); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (suite *WebDriverTestsSuite) testBasic(seleniumCaps selenium.Capabilities) (sessionID string) {
	suite.T().Helper()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // <--- Problem
	}
	http.DefaultClient.Transport = tr
	wd, err := selenium.NewRemote(seleniumCaps, suite.wdURL)
	require.NoError(suite.T(), err)
	defer wd.Quit()

	// Navigate to the simple playground interface.
	err = wd.Get("https://go.dev/play/?simple=1")
	require.NoError(suite.T(), err)

	_, err = wd.Screenshot()
	require.NoError(suite.T(), err)

	// Get a reference to the text box containing code.
	elem, err := wd.FindElement(selenium.ByCSSSelector, "#code")
	require.NoError(suite.T(), err)

	// Remove the boilerplate code already in the text box.
	require.NoError(suite.T(), elem.Clear())

	// Enter some new code in text box.
	err = elem.SendKeys(`
		package main
		import "fmt"
		func main() {
			fmt.Println("Hello WebDriver!\n")
		}
	`)
	require.NoError(suite.T(), err)

	_, err = wd.Screenshot()
	require.NoError(suite.T(), err)

	// Click the run button.
	btn, err := wd.FindElement(selenium.ByCSSSelector, "#run")
	require.NoError(suite.T(), err)

	require.NoError(suite.T(), btn.Click())

	// Wait for the program to finish running and get the output.
	outputDiv, err := wd.FindElement(selenium.ByCSSSelector, "pre.Playground-output")
	require.NoError(suite.T(), err)

	for {
		output, err := outputDiv.Text()
		require.NoError(suite.T(), err)

		if output != "Waiting for remote server..." {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	_, err = wd.Screenshot()
	require.NoError(suite.T(), err)

	// Wait for the program to finish running and get the output.
	outputPre, err := outputDiv.FindElement(selenium.ByCSSSelector, "span.stdout")
	require.NoError(suite.T(), err)

	var output string
	output, err = outputPre.Text()
	require.NoError(suite.T(), err)
	time.Sleep(time.Second * 20)

	fmt.Printf("%s", strings.Replace(output, "\n\n", "\n", -1))

	// Example Output:
	// Hello WebDriver!
	//
	// Program exited.

	return wd.SessionID()
}
