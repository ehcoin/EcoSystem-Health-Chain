// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package ole

import "unsafe"

type IEnumVARIANT struct {
	IUnknown
}

type IEnumVARIANTVtbl struct {
	IUnknownVtbl
	Next  uintptr
	Skip  uintptr
	Reset uintptr
	Clone uintptr
}

func (v *IEnumVARIANT) VTable() *IEnumVARIANTVtbl {
	return (*IEnumVARIANTVtbl)(unsafe.Pointer(v.RawVTable))
}
