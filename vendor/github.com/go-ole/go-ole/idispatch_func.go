// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
// +build !windows

package ole

func getIDsOfName(disp *IDispatch, names []string) ([]int32, error) {
	return []int32{}, NewError(E_NOTIMPL)
}

func getTypeInfoCount(disp *IDispatch) (uint32, error) {
	return uint32(0), NewError(E_NOTIMPL)
}

func getTypeInfo(disp *IDispatch) (*ITypeInfo, error) {
	return nil, NewError(E_NOTIMPL)
}

func invoke(disp *IDispatch, dispid int32, dispatch int16, params ...interface{}) (*VARIANT, error) {
	return nil, NewError(E_NOTIMPL)
}
