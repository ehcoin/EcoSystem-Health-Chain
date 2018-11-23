// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package metrics

import (
	"net"
	"time"
)

func ExampleOpenTSDB() {
	addr, _ := net.ResolveTCPAddr("net", ":2003")
	go OpenTSDB(DefaultRegistry, 1*time.Second, "some.prefix", addr)
}

func ExampleOpenTSDBWithConfig() {
	addr, _ := net.ResolveTCPAddr("net", ":2003")
	go OpenTSDBWithConfig(OpenTSDBConfig{
		Addr:          addr,
		Registry:      DefaultRegistry,
		FlushInterval: 1 * time.Second,
		DurationUnit:  time.Millisecond,
	})
}
