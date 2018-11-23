// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
// +build appengine

package isatty

// IsTerminal returns true if the file descriptor is terminal which
// is always false on on appengine classic which is a sandboxed PaaS.
func IsTerminal(fd uintptr) bool {
	return false
}
