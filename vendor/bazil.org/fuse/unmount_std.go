// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
// +build !linux

package fuse

import (
	"os"
	"syscall"
)

func unmount(dir string) error {
	err := syscall.Unmount(dir, 0)
	if err != nil {
		err = &os.PathError{Op: "unmount", Path: dir, Err: err}
		return err
	}
	return nil
}
