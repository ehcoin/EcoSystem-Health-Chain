// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package duktape

// Must returns existing *Context or throw panic.
// It is highly recommended to use Must all the time.
func (d *Context) Must() *Context {
	if d.duk_context == nil {
		panic("[duktape] Context does not exists!\nYou cannot call any contexts methods after `DestroyHeap()` was called.")
	}
	return d
}
