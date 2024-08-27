// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"regexp"
	"strings"

	"github.com/opentofu/libregistry/vcs"
	"golang.org/x/mod/semver"
)

var versionRe = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+(|-[a-zA-Z0-9-.]+)$`)

const maxVersionLength = 255

// VersionNumber describes the semver version number. Note that in contrast to provider versions module versions
// do not have a compulsory "v" prefix. Call ToVCSVersion() before you call Normalize() in order to get the correct
// VCS version.
type VersionNumber string

// Normalize adds a "v" prefix if none is present. Note, however, that in contrast to provider versions module versions
// do not have a compulsory "v" prefix. Call ToVCSVersion() before you call Normalize() in order to get the correct
// VCS version.
func (v VersionNumber) Normalize() VersionNumber {
	return VersionNumber("v" + strings.TrimPrefix(string(v), "v"))
}

func (v VersionNumber) Compare(other VersionNumber) int {
	return semver.Compare(string(v.Normalize()), string(other.Normalize()))
}

func (v VersionNumber) Validate() error {
	normalizedV := v.Normalize()
	if len(normalizedV) > maxVersionLength {
		return &InvalidVersionNumber{v}
	}
	if !versionRe.MatchString(string(normalizedV)) {
		return &InvalidVersionNumber{v}
	}
	return nil
}

// ToVCSVersion creates a vcs.VersionNumber from the VersionNumber. Note that in contrast to provider versions module
// versions do not have a compulsory "v" prefix. Call ToVCSVersion() before you call Normalize() in order to get the
// correct VCS version.
func (v VersionNumber) ToVCSVersion() vcs.VersionNumber {
	return vcs.VersionNumber(v)
}

type InvalidVersionNumber struct {
	VersionNumber VersionNumber
}

func (i InvalidVersionNumber) Error() string {
	return "Invalid version: " + string(i.VersionNumber)
}
