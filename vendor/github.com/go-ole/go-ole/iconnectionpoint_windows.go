// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
// +build windows

package ole

import (
	"syscall"
	"unsafe"
)

func (v *IConnectionPoint) GetConnectionInterface(piid **GUID) int32 {
	// XXX: This doesn't look like it does what it's supposed to
	return release((*IUnknown)(unsafe.Pointer(v)))
}

func (v *IConnectionPoint) Advise(unknown *IUnknown) (cookie uint32, err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().Advise,
		3,
		uintptr(unsafe.Pointer(v)),
		uintptr(unsafe.Pointer(unknown)),
		uintptr(unsafe.Pointer(&cookie)))
	if hr != 0 {
		err = NewError(hr)
	}
	return
}

func (v *IConnectionPoint) Unadvise(cookie uint32) (err error) {
	hr, _, _ := syscall.Syscall(
		v.VTable().Unadvise,
		2,
		uintptr(unsafe.Pointer(v)),
		uintptr(cookie),
		0)
	if hr != 0 {
		err = NewError(hr)
	}
	return
}

func (v *IConnectionPoint) EnumConnections(p *unsafe.Pointer) error {
	return NewError(E_NOTIMPL)
}
