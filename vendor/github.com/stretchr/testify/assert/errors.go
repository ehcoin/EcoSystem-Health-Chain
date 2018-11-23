// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package assert

import (
	"errors"
)

// AnError is an error instance useful for testing.  If the code does not care
// about error specifics, and only needs to return the error for example, this
// error should be used to make the test code more readable.
var AnError = errors.New("assert.AnError general error for testing")
