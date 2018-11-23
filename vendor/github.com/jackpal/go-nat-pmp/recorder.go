// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package natpmp

import "time"

type callObserver interface {
	observeCall(msg []byte, result []byte, err error)
}

// A caller that records the RPC call.
type recorder struct {
	child    caller
	observer callObserver
}

func (n *recorder) call(msg []byte, timeout time.Duration) (result []byte, err error) {
	result, err = n.child.call(msg, timeout)
	n.observer.observeCall(msg, result, err)
	return
}
