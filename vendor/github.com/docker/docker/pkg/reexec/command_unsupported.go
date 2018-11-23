// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
// +build !linux,!windows,!freebsd,!solaris,!darwin

package reexec

import (
	"os/exec"
)

// Command is unsupported on operating systems apart from Linux, Windows, Solaris and Darwin.
func Command(args ...string) *exec.Cmd {
	return nil
}
