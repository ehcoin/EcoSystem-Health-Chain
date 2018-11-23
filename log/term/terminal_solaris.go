// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package term

import "golang.org/x/sys/unix"

// IsTty returns true if the given file descriptor is a terminal.
func IsTty(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TCGETA)
	return err == nil
}
