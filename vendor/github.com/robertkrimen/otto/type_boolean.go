// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package otto

import (
	"strconv"
)

func (runtime *_runtime) newBooleanObject(value Value) *_object {
	return runtime.newPrimitiveObject("Boolean", toValue_bool(value.bool()))
}

func booleanToString(value bool) string {
	return strconv.FormatBool(value)
}
