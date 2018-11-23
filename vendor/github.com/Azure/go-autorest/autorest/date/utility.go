// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package date

import (
	"strings"
	"time"
)

// ParseTime to parse Time string to specified format.
func ParseTime(format string, t string) (d time.Time, err error) {
	return time.Parse(format, strings.ToUpper(t))
}
