package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	pm0 "github.com/TrixiS/pm0/internal/cli"
	"github.com/TrixiS/pm0/internal/cli/command_context"
)

const releaseURL string = "https://api.github.com/repos/TrixiS/pm0/releases/latest"

type releaseAsset struct {
	BrowserDownloadURL string `json:"browser_download_url"`
}

type release struct {
	Assets []releaseAsset `json:"assets"`
}

func Update(ctx *command_context.CommandContext) error {
	res, err := http.Get(releaseURL)

	if err != nil {
		return fmt.Errorf("failed to get release %s: %w", releaseURL, err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get release %s, code %d", releaseURL, res.StatusCode)
	}

	bodyBytes, err := io.ReadAll(res.Body)

	if err != nil {
		return fmt.Errorf("failed to read release %s: %w", releaseURL, err)
	}

	r := release{}

	if err := json.Unmarshal(bodyBytes, &r); err != nil {
		return fmt.Errorf("failed to parse release body %s: %w", releaseURL, err)
	}

	if len(r.Assets) == 0 {
		pm0.Printf("no release assets found %s", releaseURL)
		return nil
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(r.Assets))

	for _, a := range r.Assets {
		go func(asset releaseAsset) {
			defer wg.Done()

			res, err := http.Get(asset.BrowserDownloadURL)

			if err != nil {
				pm0.Printf("failed to get asset %s: %v", asset.BrowserDownloadURL, err)
				return
			}

			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				pm0.Printf(
					"failed to get asset %s, code %d",
					asset.BrowserDownloadURL,
					res.StatusCode,
				)

				return
			}

			contentDispositionHeader := res.Header.Get("Content-Disposition")
			filename := getAssetFilename(contentDispositionHeader)

			if filename == "" {
				pm0.Printf(
					"failed to get asset filename from content disposition header %s",
					contentDispositionHeader,
				)

				return
			}

			assetFile, err := os.OpenFile("./"+filename, os.O_CREATE|os.O_RDWR, 0o777)

			if err != nil {
				pm0.Printf("failed to open asset file %s: %v", filename, err)
				return
			}

			defer assetFile.Close()

			if _, err := io.Copy(assetFile, res.Body); err != nil {
				pm0.Printf("failed to read asset %s: %w", asset.BrowserDownloadURL, err)
				return
			}

			pm0.Printf("updated %s", assetFile.Name())
		}(a)
	}

	wg.Wait()
	return nil
}

func getAssetFilename(contentDisposition string) string {
	parts := strings.SplitN(contentDisposition, ";", 3)

	if len(parts) < 2 {
		return ""
	}

	filenamePart := parts[1]
	filenameParts := strings.SplitN(filenamePart, "=", 3)

	if len(filenameParts) < 2 {
		return ""
	}

	return filenameParts[1]
}
