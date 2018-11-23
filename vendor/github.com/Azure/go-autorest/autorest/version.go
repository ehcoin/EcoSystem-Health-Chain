// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package autorest

import (
	"fmt"
)

const (
	major        = "7"
	minor        = "3"
	patch        = "0"
	tag          = ""
	semVerFormat = "%s.%s.%s%s"
)

var version string

// Version returns the semantic version (see http://semver.org).
func Version() string {
	if version == "" {
		version = fmt.Sprintf(semVerFormat, major, minor, patch, tag)
	}
	return version
}
