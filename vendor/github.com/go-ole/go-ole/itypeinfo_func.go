// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
// +build !windows

package ole

func (v *ITypeInfo) GetTypeAttr() (*TYPEATTR, error) {
	return nil, NewError(E_NOTIMPL)
}
