// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
//+build linux

package main

import (
	"fmt"
	"os"
	"syscall"
)

func initPanicFile() {
	file, err := os.OpenFile(panicFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Create panic file err", err)
	}
	globalFile = file
	err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
	if err != nil {

		fmt.Println("dup2 failed", err)
	}
}
