// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
// +build windows

package ole

import (
	"syscall"
	"unsafe"
)

func (enum *IEnumVARIANT) Clone() (cloned *IEnumVARIANT, err error) {
	hr, _, _ := syscall.Syscall(
		enum.VTable().Clone,
		2,
		uintptr(unsafe.Pointer(enum)),
		uintptr(unsafe.Pointer(&cloned)),
		0)
	if hr != 0 {
		err = NewError(hr)
	}
	return
}

func (enum *IEnumVARIANT) Reset() (err error) {
	hr, _, _ := syscall.Syscall(
		enum.VTable().Reset,
		1,
		uintptr(unsafe.Pointer(enum)),
		0,
		0)
	if hr != 0 {
		err = NewError(hr)
	}
	return
}

func (enum *IEnumVARIANT) Skip(celt uint) (err error) {
	hr, _, _ := syscall.Syscall(
		enum.VTable().Skip,
		2,
		uintptr(unsafe.Pointer(enum)),
		uintptr(celt),
		0)
	if hr != 0 {
		err = NewError(hr)
	}
	return
}

func (enum *IEnumVARIANT) Next(celt uint) (array VARIANT, length uint, err error) {
	hr, _, _ := syscall.Syscall6(
		enum.VTable().Next,
		4,
		uintptr(unsafe.Pointer(enum)),
		uintptr(celt),
		uintptr(unsafe.Pointer(&array)),
		uintptr(unsafe.Pointer(&length)),
		0,
		0)
	if hr != 0 {
		err = NewError(hr)
	}
	return
}
