// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
// +build !windows

package ole

func (v *IInspectable) GetIids() ([]*GUID, error) {
	return []*GUID{}, NewError(E_NOTIMPL)
}

func (v *IInspectable) GetRuntimeClassName() (string, error) {
	return "", NewError(E_NOTIMPL)
}

func (v *IInspectable) GetTrustLevel() (uint32, error) {
	return uint32(0), NewError(E_NOTIMPL)
}
