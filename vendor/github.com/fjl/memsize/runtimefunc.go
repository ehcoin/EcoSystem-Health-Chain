// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package memsize

import "unsafe"

var _ = unsafe.Pointer(nil)

//go:linkname stopTheWorld runtime.stopTheWorld
func stopTheWorld(reason string)

//go:linkname startTheWorld runtime.startTheWorld
func startTheWorld()

//go:linkname chanbuf runtime.chanbuf
func chanbuf(ch unsafe.Pointer, i uint) unsafe.Pointer
