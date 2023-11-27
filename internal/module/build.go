package module

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/opentofu/registry-stable/internal"

	"golang.org/x/mod/semver"
)

func (m Module) UpdateMetadataFile() error {
	m.Logger.Info("Beginning version bump process for module", slog.String("module", m.Namespace+"/"+m.Name+"/"+m.TargetSystem))

	meta, err := m.BuildMetadata()
	if err != nil {
		return err
	}

	return m.WriteMetadata(*meta)
}

// BuildMetadata builds the Metadata for the module by collating the tags from
// the module repository.
func (m Module) BuildMetadata() (*Metadata, error) {
	tags, err := m.getSemverTags()
	if err != nil {
		return nil, err
	}

	versions := make([]Version, len(tags))
	for i, t := range tags {
		versions[i] = Version{Version: t}
	}

	return &Metadata{Versions: versions}, nil
}

func (m Module) getSemverTags() ([]string, error) {
	tags, err := m.Github.GetTags(m.RepositoryURL())
	if err != nil {
		return nil, err
	}

	var semverTags = make([]string, 0)
	for _, tag := range tags {
		tagWithPrefix := fmt.Sprintf("v%s", internal.TrimTagPrefix(tag))
		if semver.IsValid(tagWithPrefix) {
			semverTags = append(semverTags, tag)
		}
	}

	semverSortFunc := func(a, b string) int {
		return -semver.Compare(fmt.Sprintf(a), fmt.Sprintf(b))
	}
	slices.SortFunc(semverTags, semverSortFunc)

	return semverTags, nil
}
