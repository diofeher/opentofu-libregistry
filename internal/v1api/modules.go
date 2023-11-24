package v1api

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/opentofu/registry-stable/internal"
	"github.com/opentofu/registry-stable/internal/files"
	"github.com/opentofu/registry-stable/internal/module"
)

type ModuleGenerator struct {
	module.Module
	module.MetadataFile
	Destination string
	log         *slog.Logger
}

func NewModuleGenerator(m module.Module, destination string) (ModuleGenerator, error) {
	metadata, err := m.ReadMetadata()
	if err != nil {
		return ModuleGenerator{}, err
	}

	return ModuleGenerator{
		m,
		metadata,
		destination,
		m.Logger,
	}, nil
}

func (m ModuleGenerator) VersionListingPath() string {
	return filepath.Join(m.Destination, "v1", "modules", m.Namespace, m.Name, m.TargetSystem, "versions")
}

func (m ModuleGenerator) VersionDownloadPath(v module.Version) string {
	return filepath.Join(m.Destination, "v1", "modules", m.Namespace, m.Name, m.TargetSystem, internal.TrimTagPrefix(v.Version), "download")
}

func (m ModuleGenerator) VersionListing() ModuleVersionListingResponse {
	versions := make([]ModuleVersionResponseItem, len(m.Versions))
	for i, v := range m.Versions {
		versions[i] = ModuleVersionResponseItem{Version: v.Version}
	}
	return ModuleVersionListingResponse{[]ModuleVersionListingResponseItem{{versions}}}
}

func (m ModuleGenerator) VersionDownloads() map[string]ModuleVersionDownloadResponse {
	downloads := make(map[string]ModuleVersionDownloadResponse)
	for _, v := range m.Versions {
		downloads[m.VersionDownloadPath(v)] = ModuleVersionDownloadResponse{Location: m.VersionDownloadURL(v)}
	}
	return downloads
}

// Generate generates the response for the module version listing API endpoints.
// For more information see
// https://opentofu.org/docs/internals/module-registry-protocol/#list-available-versions-for-a-specific-module
// https://opentofu.org/docs/internals/module-registry-protocol/#download-source-code-for-a-specific-module-version
func (m ModuleGenerator) Generate() error {
	m.log.Info("Generating")

	for location, download := range m.VersionDownloads() {
		err := files.SafeWriteObjectToJsonFile(location, download)
		if err != nil {
			return fmt.Errorf("failed to write metadata version download file: %w", err)
		}
		m.log.Debug("Wrote metadata version download file", slog.String("path", location))
	}

	err := files.SafeWriteObjectToJsonFile(m.VersionListingPath(), m.VersionListing())
	if err != nil {
		return err
	}

	m.log.Info("Generated")

	return nil
}
