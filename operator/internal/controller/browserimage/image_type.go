package browserimage

import (
	"fmt"
	"strings"
)

type ImageType int

const (
	ImageTypeSelenium = iota
	ImageTypeSelenoid
	ImageTypeAerokube
	ImageTypeMicrosoft
)

var homeDirMapping = map[ImageType]string{
	ImageTypeSelenium: "/home/seluser",
	ImageTypeSelenoid: "/home/user",
	ImageTypeAerokube: "/home/user",
}

var vncPassMapping = map[ImageType]string{
	ImageTypeSelenium: "browserkube",
	ImageTypeSelenoid: "selenoid",
	ImageTypeAerokube: "browserkube",
}

func (it ImageType) Homedir() string {
	return homeDirMapping[it]
}

func (it ImageType) VncPass() string {
	return vncPassMapping[it]
}

func ParseImageType(img string) (ImageType, error) {
	var idx int
	if idx = strings.LastIndexByte(img, '/'); idx < 0 {
		return -1, fmt.Errorf("unable to parse image type from: %s", img)
	}
	imgOrg := img[:idx]
	switch imgOrg {
	case "selenium":
		return ImageTypeSelenium, nil
	case "selenoid":
		return ImageTypeSelenoid, nil
	case "quay.io/browser", "cdtp", "playwright":
		return ImageTypeAerokube, nil
	case "mcr.microsoft.com":
		return ImageTypeMicrosoft, nil

	}
	return -1, fmt.Errorf("unknown image type: %s", imgOrg)
}
