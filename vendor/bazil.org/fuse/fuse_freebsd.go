// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package fuse

// Maximum file write size we are prepared to receive from the kernel.
//
// This number is just a guess.
const maxWrite = 128 * 1024