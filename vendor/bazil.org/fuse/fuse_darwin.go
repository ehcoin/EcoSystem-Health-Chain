// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package fuse

// Maximum file write size we are prepared to receive from the kernel.
//
// This value has to be >=16MB or OSXFUSE (3.4.0 observed) will
// forcibly close the /dev/fuse file descriptor on a Setxattr with a
// 16MB value. See TestSetxattr16MB and
// https://github.com/bazil/fuse/issues/42
const maxWrite = 16 * 1024 * 1024
