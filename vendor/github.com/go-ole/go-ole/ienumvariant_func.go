// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
// +build !windows

package ole

func (enum *IEnumVARIANT) Clone() (*IEnumVARIANT, error) {
	return nil, NewError(E_NOTIMPL)
}

func (enum *IEnumVARIANT) Reset() error {
	return NewError(E_NOTIMPL)
}

func (enum *IEnumVARIANT) Skip(celt uint) error {
	return NewError(E_NOTIMPL)
}

func (enum *IEnumVARIANT) Next(celt uint) (VARIANT, uint, error) {
	return NewVariant(VT_NULL, int64(0)), 0, NewError(E_NOTIMPL)
}
