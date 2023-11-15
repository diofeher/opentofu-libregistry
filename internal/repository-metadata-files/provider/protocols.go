package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"registry-stable/internal/github"
)

type Manifest struct {
	Metadata ManifestMetadata `json:"metadata"`
}
type ManifestMetadata struct {
	ProtocolVersions []string `json:"protocol_versions"`
}

func GetProtocols(ctx context.Context, manifestDownloadUrl string) ([]string, error) {
	if manifestDownloadUrl == "" {
		slog.Warn("Could not find manifest file, using default protocols")
		return []string{"5.0"}, nil
	}

	assetContents, err := github.DownloadAssetContents(ctx, manifestDownloadUrl)
	if err != nil {
		return nil, err
	}

	contents, contentsErr := io.ReadAll(assetContents)
	if contentsErr != nil {
		return nil, fmt.Errorf("failed to read asset contents: %w", contentsErr)
	}

	manifest, err := parseManifestContents(contents)
	if err != nil {
		return nil, err
	}

	return manifest.Metadata.ProtocolVersions, nil
}

func parseManifestContents(contents []byte) (*Manifest, error) {
	var manifest *Manifest
	err := json.Unmarshal(contents, &manifest)
	if err != nil {
		slog.Error("Failed to parse manifest contents")
		return nil, err
	}

	return manifest, nil
}
