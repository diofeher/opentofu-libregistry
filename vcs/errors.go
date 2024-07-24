// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package vcs

type RequestFailedError struct {
	Cause error
}

func (r RequestFailedError) Error() string {
	return "VCS request failed: " + r.Cause.Error()
}

func (r RequestFailedError) Unwrap() error {
	return r.Cause
}
