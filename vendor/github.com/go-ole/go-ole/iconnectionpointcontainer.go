// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package ole

import "unsafe"

type IConnectionPointContainer struct {
	IUnknown
}

type IConnectionPointContainerVtbl struct {
	IUnknownVtbl
	EnumConnectionPoints uintptr
	FindConnectionPoint  uintptr
}

func (v *IConnectionPointContainer) VTable() *IConnectionPointContainerVtbl {
	return (*IConnectionPointContainerVtbl)(unsafe.Pointer(v.RawVTable))
}
