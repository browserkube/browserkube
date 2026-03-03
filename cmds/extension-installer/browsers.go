package extensioninstaller

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const googleConnectionURL = `%s?response=redirect&os=linux&arch=x64&os_arch=x86_64` +
	`&nacl_arch=x86-64&prod=chromium&prodchannel=unknown&prodversion=%s&lang=en-US` +
	`&acceptformat=crx2,crx3&x=id%%3D%s%%26installsource%%3Dondemand%%26uc`

func installChromeExtension(ctx context.Context, extensionID string) error {
	fmt.Printf("Installing plugin %s for chrome\n", extensionID)
	//nolint:gosec // not a credentials
	if err := os.MkdirAll(linuxChromePluginLoc, os.ModePerm); err != nil {
		return fmt.Errorf("error while creating folder:%w", err)
	}

	updateUrl := map[string]string{"extension_update_url": chromeStoreUpdateURL}
	b, err := json.Marshal(updateUrl)
	if err != nil {
		return fmt.Errorf("error while marshaling json:%w", err)
	}

	//nolint:gosec // not a credentials
	if wErr := os.WriteFile(
		fmt.Sprintf("%s/%s.json", linuxChromePluginLoc, extensionID),
		b,
		0o644,
	); wErr != nil {
		return fmt.Errorf("error while creating file: %w", wErr)
	}
	fmt.Printf("Copying .crx file to %s for manual installation\n", linuxChromePluginLoc)

	filepath := fmt.Sprintf("/opt/extensions/%s.crx", extensionID)
	if fErr := getFile(ctx,
		filepath,
		fmt.Sprintf(googleConnectionURL, chromeStoreUpdateURL, browserVersionChrome, extensionID),
	); fErr != nil {
		return fmt.Errorf("error while downloading file:%w", fErr)
	}

	err = unzip(filepath, fmt.Sprintf("%s/%s", linuxChromePluginLoc, extensionID))
	if err != nil {
		return fmt.Errorf("error while unzipping file:%w", err)
	}
	fmt.Printf("%s copied\n", extensionID)

	if err := os.RemoveAll(fmt.Sprintf("%s/%s/_metadata", linuxChromePluginLoc, extensionID)); err != nil {
		return fmt.Errorf("error while removing package _metadata: %w", err)
	}

	return nil
}

func installFirefoxExtension(ctx context.Context, extensionID, updateURL string) error {
	fmt.Printf("Installing plugin %s for firefox\n", extensionID)
	if err := getFile(ctx, fmt.Sprintf("/opt/extensions/%s.xpi", extensionID), updateURL); err != nil {
		return err
	}
	//nolint:gosec // not a credentials
	if err := os.MkdirAll(fmt.Sprintf("%s/distribution/extensions", linuxFirefoxInstallationLoc), 0o777); err != nil {
		return fmt.Errorf("error while creating folder: %w", err)
	}
	// Copy the extension files to the installation locations
	bRead, err := os.ReadFile(fmt.Sprintf("/opt/extensions/%s.xpi", extensionID))
	if err != nil {
		return fmt.Errorf("error while opening file: %w", err)
	}

	//nolint:gosec // not a credentials
	if err := os.WriteFile(fmt.Sprintf(linuxFirefoxPluginLoc+"/%s.xpi", extensionID), bRead, 0o644); err != nil {
		return fmt.Errorf("error while writing to file: %w", err)
	}
	//nolint:gosec // not a credentials
	if err := os.WriteFile(fmt.Sprintf("%s/distribution/extensions/%s.xpi", linuxFirefoxInstallationLoc, extensionID), bRead, 0o644); err != nil {
		return fmt.Errorf("error while writing to file: %w", err)
	}
	return nil
}

func getFile(ctx context.Context, outputLoc, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error while creating request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req) //nolint:gosec //this is trusted request since url is from a trusted source
	if err != nil {
		return fmt.Errorf("error while sending request: %w", err)
	}
	defer resp.Body.Close()          //nolint:errcheck
	out, err := os.Create(outputLoc) //nolint:gosec //this is trusted path
	if err != nil {
		return fmt.Errorf("error while creating file: %w", err)
	}
	defer out.Close() //nolint:errcheck
	_, _ = io.Copy(out, resp.Body)
	return nil
}

func unzip(source, destination string) error {
	archive, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer archive.Close() //nolint:errcheck
	for _, file := range archive.File {
		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()                          //nolint:errcheck
		path := filepath.Join(destination, file.Name) //nolint:gosec //this is a trusted path since destination is from a trusted source
		_ = os.Remove(path)
		//nolint:gosec //not a threat since path is from a trusted source
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		if file.FileInfo().IsDir() {
			continue
		}
		err = os.Remove(path)
		if err != nil {
			return err
		}
		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode()) //nolint:gosec //this is a trusted path
		if err != nil {
			return err
		}
		defer writer.Close() //nolint:errcheck
		//nolint:gosec //not a threat since reader is from a trusted source
		_, err = io.Copy(writer, reader)
		if err != nil {
			return err
		}
	}
	return nil
}
