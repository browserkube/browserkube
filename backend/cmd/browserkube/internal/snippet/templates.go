package snippet

import (
	"embed"
	"fmt"
	"strings"
	"text/template"
)

const (
	browserChrome  = "chrome"
	browserFirefox = "firefox"
)

type TemplateOpts struct {
	BrowserName    string
	BrowserVersion string
}

var snippetTemplatesFuncMap = template.FuncMap{
	"seleniumJavaBrowserOptions": seleniumJavaBrowserOptions,
}

//go:embed templates/*
var snippetTemplatesFS embed.FS

var snippetTemplates = template.Must(template.New("snippets").
	Funcs(snippetTemplatesFuncMap).
	Delims("{{{", "}}}").ParseFS(snippetTemplatesFS, "templates/*.tmpl"))

func getTemplateName(sessionType, snippetLanguage string) string {
	return fmt.Sprintf("%s-%s.tmpl", sessionType, snippetLanguage)
}

func seleniumJavaBrowserOptions(browserName string) string {
	switch browserName {
	case browserChrome:
		return "ChromeOptions options = new ChromeOptions()"
	case browserFirefox:
		return "FirefoxOptions options = new FirefoxOptions()"
	}

	return "SpecificBrowserOptions options = new SpecificBrowserOptions()"
}

func GetSnippet(sessionType, snippetLanguage, browserName, browserVersion string) (string, error) {
	var sb strings.Builder

	err := snippetTemplates.ExecuteTemplate(&sb,
		getTemplateName(sessionType, snippetLanguage),
		TemplateOpts{
			BrowserName:    browserName,
			BrowserVersion: browserVersion,
		})
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}
