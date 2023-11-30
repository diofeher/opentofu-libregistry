package v1api

import (
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/opentofu/registry-stable/internal/files"
	"github.com/opentofu/registry-stable/internal/gpg"
	"github.com/opentofu/registry-stable/internal/provider"
)

// ProviderGenerator is responsible for generating the response for the provider version listing API endpoints.
type ProviderGenerator struct {
	provider.Provider
	provider.Metadata

	KeyLocation string
	Destination string
	log         *slog.Logger
}

// NewProviderGenerator creates a new ProviderGenerator which will generate the response for the provider version listing API endpoints and write it to the given destination.
func NewProviderGenerator(p provider.Provider, destination string, gpgKeyLocation string) (ProviderGenerator, error) {
	metadata, err := p.ReadMetadata()
	if err != nil {
		return ProviderGenerator{}, err
	}

	return ProviderGenerator{
		Provider: p,
		Metadata: metadata,

		KeyLocation: gpgKeyLocation,
		Destination: destination,
		log:         p.Logger,
	}, err
}

// VersionListingPath returns the path to the provider version listing file.
func (p ProviderGenerator) VersionListingPath() string {
	return filepath.Join(p.Destination, "v1", "providers", p.Provider.Namespace, p.Provider.ProviderName, "versions")
}

// VersionDownloadPath returns the path to the provider version download file.
func (p ProviderGenerator) VersionDownloadPath(ver provider.Version, details ProviderVersionDetails) string {
	return filepath.Join(p.Destination, "v1", "providers", p.Provider.Namespace, p.Provider.ProviderName, ver.Version, "download", details.OS, details.Arch)
}

// VersionListing will take the provider metadata and generate the responses for the provider version listing API endpoints.
func (p ProviderGenerator) VersionListing() ProviderVersionListingResponse {
	versions := make([]ProviderVersionResponseItem, len(p.Metadata.Versions))

	for versionIdx, ver := range p.Metadata.Versions {
		verResp := ProviderVersionResponseItem{
			Version:   ver.Version,
			Protocols: ver.Protocols,
			Platforms: make([]Platform, len(ver.Targets)),
		}

		for targetIdx, target := range ver.Targets {
			verResp.Platforms[targetIdx] = Platform{
				OS:   target.OS,
				Arch: target.Arch,
			}
		}
		versions[versionIdx] = verResp
	}

	return ProviderVersionListingResponse{versions}
}

// VersionDetails will take the provider metadata and generate the responses for the provider version download API endpoints.
func (p ProviderGenerator) VersionDetails() (map[string]ProviderVersionDetails, error) {
	versionDetails := make(map[string]ProviderVersionDetails)

	keyCollection := gpg.KeyCollection{
		Namespace: p.Provider.EffectiveNamespace(),
		Directory: p.KeyLocation,
	}

	keys, err := keyCollection.ListKeys()
	if err != nil {
		p.log.Error("Failed to list keys", slog.Any("err", err))
		return nil, err
	}

	for _, ver := range p.Metadata.Versions {
		for _, target := range ver.Targets {
			details := ProviderVersionDetails{
				Protocols:           ver.Protocols,
				OS:                  target.OS,
				Arch:                target.Arch,
				Filename:            target.Filename,
				DownloadURL:         target.DownloadURL,
				SHASumsURL:          ver.SHASumsURL,
				SHASumsSignatureURL: ver.SHASumsSignatureURL,
				SHASum:              target.SHASum,
				SigningKeys: SigningKeys{
					GPGPublicKeys: keys,
				},
			}
			versionDetails[p.VersionDownloadPath(ver, details)] = details
		}
	}
	return versionDetails, nil
}

// Generate generates the responses for the provider version listing API endpoints.
func (p ProviderGenerator) Generate() error {
	p.log.Info("Generating")

	details, err := p.VersionDetails()
	if err != nil {
		return err
	}

	for location, details := range details {
		err := files.SafeWriteObjectToJSONFile(location, details)
		if err != nil {
			return fmt.Errorf("failed to write metadata version download file: %w", err)
		}
	}

	err = files.SafeWriteObjectToJSONFile(p.VersionListingPath(), p.VersionListing())
	if err != nil {
		return err
	}

	p.log.Info("Generated")

	return nil
}
