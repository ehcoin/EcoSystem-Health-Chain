// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package ole

import "unsafe"

type IConnectionPoint struct {
	IUnknown
}

type IConnectionPointVtbl struct {
	IUnknownVtbl
	GetConnectionInterface      uintptr
	GetConnectionPointContainer uintptr
	Advise                      uintptr
	Unadvise                    uintptr
	EnumConnections             uintptr
}

func (v *IConnectionPoint) VTable() *IConnectionPointVtbl {
	return (*IConnectionPointVtbl)(unsafe.Pointer(v.RawVTable))
}
