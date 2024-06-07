package controller

import "flag"

type BrowserCtrlOpts struct {
	OperatorNamespace       string
	sidecarImage            string
	recorderImage           string
	extensionInstallerImage string
	xServerImage            string
	vncServerImage          string
	clipboardImage          string
	sidecarPort             string
	browserUserConfig       string
	browserExtensionConfig  string
	browserReadinessConfig  string
}

func InitBrowserCtrlOpts() *BrowserCtrlOpts {
	cfg := &BrowserCtrlOpts{}

	flag.StringVar(&cfg.OperatorNamespace, "operator-namespace", "browserkube", "Operator's namespace")

	flag.StringVar(&cfg.sidecarImage, "sidecar-image", "", "Image of sidecar to be used")
	flag.StringVar(&cfg.xServerImage, "x-server-image", "", "Image of x-server to be used")
	flag.StringVar(&cfg.vncServerImage, "vnc-server-image", "", "Image of vnc-server to be used")
	flag.StringVar(&cfg.clipboardImage, "clipboard-image", "", "Image of clipboard to be used")
	flag.StringVar(&cfg.recorderImage, "recorder-image", "", "Image of recorder to be used")
	flag.StringVar(&cfg.extensionInstallerImage, "extension-installer-image", "", "Image of extension-installer to be used")
	flag.StringVar(&cfg.sidecarPort, "sidecar-port", "9999", "Port of sidecar")
	flag.StringVar(&cfg.browserUserConfig, "browser-user-configmap", "browserkube-browsers-usergroup", "Browser Config Map Name")
	flag.StringVar(&cfg.browserExtensionConfig, "browser-extension-configmap", "browserkube-browser-extension-config", "Browser Config Map Name")
	flag.StringVar(&cfg.browserReadinessConfig, "browser-readinessprobe-configmap", "browserkube-browsers-readinessprobe-config", "Browser Readiness Config Map Name")

	return cfg
}
