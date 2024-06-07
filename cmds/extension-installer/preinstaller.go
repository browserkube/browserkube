package extensioninstaller

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

const (
	flagExtensionID = "extensionId"
	flagUpdateURL   = "updateUrl"
	flagBrowserName = "browserName"
)

// browser plugin constants
const (
	// user home directory
	selUserDir                  = "/home/seluser"
	linuxChromePluginLoc        = "/opt/google/chrome/extensions"
	linuxFirefoxPluginLoc       = selUserDir + "/.mozilla/extensions"
	linuxFirefoxInstallationLoc = "/opt/firefox"
	chromeStoreUpdateURL        = "https://clients2.google.com/service/update2/crx"
)

// browser constants
const (
	browserNameChrome    = "chrome"
	browserVersionChrome = "108.0.5359.124"
	browserNameFirefox   = "firefox"
)

const (
	whitelistChromeEnv  = "BROWSERKUBE-BROWSER-EXTENSION-CONFIG-WHITELIST-CHROME"
	whitelistFirefoxEnv = "BROWSERKUBE-BROWSER-EXTENSION-CONFIG-WHITELIST-FIREFOX"
)

func NewApp() *cli.App {
	return &cli.App{
		Name:      "extension-installer",
		Usage:     "made for using as initContainer before deploying browser pods",
		Action:    PreInstall,
		Reader:    os.Stdin,
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     flagBrowserName,
				Aliases:  []string{"b"},
				Usage:    "browser name",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     flagExtensionID,
				Aliases:  []string{"e"},
				Usage:    "extension id from online store",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     flagUpdateURL,
				Aliases:  []string{"u"},
				Usage:    "update url for firefox browser",
				Required: true,
			},
		},
	}
}

func PreInstall(ctx *cli.Context) error {
	extensionID := ctx.StringSlice(flagExtensionID)
	updateUrl := ctx.StringSlice(flagUpdateURL)
	browserName := ctx.StringSlice(flagBrowserName)
	log.Printf("Current Config : %s, %s, %s", extensionID, updateUrl, browserName)

	if len(updateUrl) != len(browserName) ||
		len(updateUrl) != len(extensionID) {
		return fmt.Errorf("invalid format")
	}
	if err := os.MkdirAll("/opt/extensions", os.ModePerm); err != nil {
		log.Printf("error while creating extensions folder:%s", err.Error())
		return nil
	}
	whiteList := getWhitelist()
	for i := range extensionID {
		switch browserName[i] {
		case browserNameChrome:
			if extensionID[i] == "" {
				return fmt.Errorf("cannot install plugin on chrome. extensionId is missing")
			}
			if !checkSlice(whiteList[browserNameChrome], extensionID[i]) {
				log.Printf("WARNING: Extension with ID :%s is not included in whitelist", extensionID[i])
				continue
			}
			err := installChromeExtension(extensionID[i])
			if err != nil {
				log.Println(err.Error())
			}
		case browserNameFirefox:
			if extensionID[i] == "" {
				return fmt.Errorf("cannot install plugin on chrome. extensionId is missing")
			}
			if updateUrl[i] == "" {
				return fmt.Errorf("cannot install plugin on firefox. updateUrl is missing")
			}
			if !checkSlice(whiteList[browserNameFirefox], extensionID[i]) {
				log.Printf("WARNING: Extension with ID :%s is not included in whitelist", extensionID[i])
				continue
			}
			err := installFirefoxExtension(extensionID[i], updateUrl[i])
			if err != nil {
				log.Println(err.Error())
			}
		default:
			return fmt.Errorf("browser name %s could not be recognized", browserName[i])
		}
	}

	if err := checkFolders(extensionID, browserName); err != nil {
		log.Print(err.Error())
	}
	return nil
}

func checkFolders(extensionID, browserName []string) error {
	for i := 0; i < len(extensionID); i++ {
		switch browserName[i] {
		case browserNameChrome:
			_, err := os.Stat(linuxChromePluginLoc + "/" + extensionID[i])
			if os.IsNotExist(err) {
				return fmt.Errorf("check failed. Extension id %s does not exists in %s", extensionID[i], linuxChromePluginLoc)
			}
		case browserNameFirefox:
			_, err := os.Stat(linuxFirefoxInstallationLoc + "/distribution/extensions/" + extensionID[i] + ".xpi")
			if os.IsNotExist(err) {
				return fmt.Errorf("check failed. Extension id %s does not exists in %s/distribution/extensions", extensionID[i], linuxFirefoxInstallationLoc)
			}
			_, err = os.Stat(linuxFirefoxPluginLoc + "/" + extensionID[i] + ".xpi")
			if os.IsNotExist(err) {
				return fmt.Errorf("check failed. Extension id %s does not exists in %s", extensionID[i], linuxFirefoxPluginLoc)
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getWhitelist() map[string][]string {
	whitelist := make(map[string][]string)
	whitelist[browserNameChrome] = strings.Split(os.Getenv(whitelistChromeEnv), ",")
	whitelist[browserNameFirefox] = strings.Split(os.Getenv(whitelistFirefoxEnv), ",")
	return whitelist
}
