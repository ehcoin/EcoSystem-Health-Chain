// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
// +build !windows

package ole

import "unsafe"

func (v *IConnectionPoint) GetConnectionInterface(piid **GUID) int32 {
	return int32(0)
}

func (v *IConnectionPoint) Advise(unknown *IUnknown) (uint32, error) {
	return uint32(0), NewError(E_NOTIMPL)
}

func (v *IConnectionPoint) Unadvise(cookie uint32) error {
	return NewError(E_NOTIMPL)
}

func (v *IConnectionPoint) EnumConnections(p *unsafe.Pointer) (err error) {
	return NewError(E_NOTIMPL)
}
