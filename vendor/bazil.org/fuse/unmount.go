// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package fuse

// Unmount tries to unmount the filesystem mounted at dir.
func Unmount(dir string) error {
	return unmount(dir)
}
