// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
// +build !windows

package ole

func reflectQueryInterface(self interface{}, method uintptr, interfaceID *GUID, obj interface{}) (err error) {
	return NewError(E_NOTIMPL)
}

func queryInterface(unk *IUnknown, iid *GUID) (disp *IDispatch, err error) {
	return nil, NewError(E_NOTIMPL)
}

func addRef(unk *IUnknown) int32 {
	return 0
}

func release(unk *IUnknown) int32 {
	return 0
}
