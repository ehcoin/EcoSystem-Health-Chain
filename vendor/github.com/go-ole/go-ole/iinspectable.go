// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package ole

import "unsafe"

type IInspectable struct {
	IUnknown
}

type IInspectableVtbl struct {
	IUnknownVtbl
	GetIIds             uintptr
	GetRuntimeClassName uintptr
	GetTrustLevel       uintptr
}

func (v *IInspectable) VTable() *IInspectableVtbl {
	return (*IInspectableVtbl)(unsafe.Pointer(v.RawVTable))
}
