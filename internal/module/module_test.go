package module

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Module(t *testing.T) {
	m := Module{
		Namespace:    "spacename",
		Name:         "name",
		TargetSystem: "target",
	}
	v := Version{
		Version: "v1.0.1",
	}
	assert.Equal(t, "https://github.com/spacename/terraform-target-name", m.RepositoryURL())
	assert.Equal(t, "git::https://github.com/spacename/terraform-target-name?ref=v1.0.1", m.VersionDownloadURL(v))
	assert.Equal(t, "dir/s/spacename/name/target.json", m.MetadataPath("./dir"))
	assert.Equal(t, "dir/v1/modules/spacename/name/target/versions", m.VersionListingPath("./dir"))
	assert.Equal(t, "dir/v1/modules/spacename/name/target/1.0.1/download", m.VersionDownloadPath("./dir", v))
}
