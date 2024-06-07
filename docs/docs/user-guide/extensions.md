---
sidebar_position: 2
---

# Browser Extensions

#### Using Browser Specific Plugins(Chrome, Firefox)
Under desired capabilities of crd browser config file enter Extension ID (for Chrome and Firefox) and DownloadURl (Firefox).
```go
caps := selenium.Capabilities{
    "browserName": "chrome",
    "goog:chromeOptions:": map[string]interface{}{
        "args": []string{"load-extension=/opt/google/chrome/extensions/" + extensionID},
    },
    "browserkube:options": map[string]interface{}{
        "enableVNC": true,
        "tenant":    "test",
        "name":      "test session name",
    },
}

```
this will install plugin to default profile of the web browsers. unpacked extension directories:
```go
//For chrome
extensionDirectory := /opt/google/chrome/extensions/EXTENSION-ID.crx
//For firefox (if using selenium standalone image)
extensionDirectory = /home/seluser/.mozilla/extensions/EXTENSION-ID.xpi
// or
extensionDirectory = /opt/firefox/distribution/extensions/EXTENSION-ID.xpi
```
Note that to install the extension from local machine, you need to use installExtension function, for more info:
https://gist.github.com/nadvolod/ac8cdf55889510fcd64434d4ee1e2a60

