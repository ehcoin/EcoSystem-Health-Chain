// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
// Copyright (c) 2012 VMware, Inc.

package gosigar

import (
	"unsafe"
)

func bytePtrToString(ptr *int8) string {
	bytes := (*[10000]byte)(unsafe.Pointer(ptr))

	n := 0
	for bytes[n] != 0 {
		n++
	}

	return string(bytes[0:n])
}

func chop(buf []byte) []byte {
	return buf[0 : len(buf)-1]
}
