// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package fuse

import (
	"syscall"
)

const (
	ENODATA = Errno(syscall.ENODATA)
)

const (
	errNoXattr = ENODATA
)

func init() {
	errnoNames[errNoXattr] = "ENODATA"
}
