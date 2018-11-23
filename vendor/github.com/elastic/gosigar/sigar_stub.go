// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
// +build !darwin,!freebsd,!linux,!openbsd,!windows

package gosigar

import (
	"runtime"
)

func (c *Cpu) Get() error {
	return ErrNotImplemented{runtime.GOOS}
}

func (l *LoadAverage) Get() error {
	return ErrNotImplemented{runtime.GOOS}
}

func (m *Mem) Get() error {
	return ErrNotImplemented{runtime.GOOS}
}

func (s *Swap) Get() error {
	return ErrNotImplemented{runtime.GOOS}
}

func (s *HugeTLBPages) Get() error {
	return ErrNotImplemented{runtime.GOOS}
}

func (f *FDUsage) Get() error {
	return ErrNotImplemented{runtime.GOOS}
}

func (p *ProcTime) Get(int) error {
	return ErrNotImplemented{runtime.GOOS}
}

func (self *FileSystemUsage) Get(path string) error {
	return ErrNotImplemented{runtime.GOOS}
}

func (self *CpuList) Get() error {
	return ErrNotImplemented{runtime.GOOS}
}

func (p *ProcState) Get(int) error {
	return ErrNotImplemented{runtime.GOOS}
}

func (p *ProcExe) Get(int) error {
	return ErrNotImplemented{runtime.GOOS}
}

func (p *ProcMem) Get(int) error {
	return ErrNotImplemented{runtime.GOOS}
}

func (p *ProcFDUsage) Get(int) error {
	return ErrNotImplemented{runtime.GOOS}
}

func (p *ProcEnv) Get(int) error {
	return ErrNotImplemented{runtime.GOOS}
}

func (p *ProcList) Get() error {
	return ErrNotImplemented{runtime.GOOS}
}

func (p *ProcArgs) Get(int) error {
	return ErrNotImplemented{runtime.GOOS}
}

func (self *Rusage) Get(int) error {
	return ErrNotImplemented{runtime.GOOS}
}
