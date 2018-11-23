// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package otto

func (runtime *_runtime) newNumberObject(value Value) *_object {
	return runtime.newPrimitiveObject("Number", value.numberValue())
}
