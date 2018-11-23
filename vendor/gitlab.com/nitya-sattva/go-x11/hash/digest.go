// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package hash

type Digest interface {
	// See hash.Hash
	Hash

	// Close the digest by writing the last bits and storing the hash
	// in dst. This prepares the digest for reuse, calls Hash.Reset.
	Close(dst []byte, bits uint8, bcnt uint8) error
}
