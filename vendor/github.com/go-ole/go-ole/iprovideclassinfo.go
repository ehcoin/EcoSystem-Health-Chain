// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package ole

import "unsafe"

type IProvideClassInfo struct {
	IUnknown
}

type IProvideClassInfoVtbl struct {
	IUnknownVtbl
	GetClassInfo uintptr
}

func (v *IProvideClassInfo) VTable() *IProvideClassInfoVtbl {
	return (*IProvideClassInfoVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *IProvideClassInfo) GetClassInfo() (cinfo *ITypeInfo, err error) {
	cinfo, err = getClassInfo(v)
	return
}
