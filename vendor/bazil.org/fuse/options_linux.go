// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package fuse

func localVolume(conf *mountConfig) error {
	return nil
}

func volumeName(name string) MountOption {
	return dummyOption
}

func daemonTimeout(name string) MountOption {
	return dummyOption
}

func noAppleXattr(conf *mountConfig) error {
	return nil
}

func noAppleDouble(conf *mountConfig) error {
	return nil
}

func exclCreate(conf *mountConfig) error {
	return nil
}
