// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
// +build windows

package ole

import (
	"syscall"
)

var (
	modcombase     = syscall.NewLazyDLL("combase.dll")
	modkernel32, _ = syscall.LoadDLL("kernel32.dll")
	modole32, _    = syscall.LoadDLL("ole32.dll")
	modoleaut32, _ = syscall.LoadDLL("oleaut32.dll")
	modmsvcrt, _   = syscall.LoadDLL("msvcrt.dll")
	moduser32, _   = syscall.LoadDLL("user32.dll")
)
