// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package fuse

func localVolume(conf *mountConfig) error {
	conf.options["local"] = ""
	return nil
}

func volumeName(name string) MountOption {
	return func(conf *mountConfig) error {
		conf.options["volname"] = name
		return nil
	}
}

func daemonTimeout(name string) MountOption {
	return func(conf *mountConfig) error {
		conf.options["daemon_timeout"] = name
		return nil
	}
}

func noAppleXattr(conf *mountConfig) error {
	conf.options["noapplexattr"] = ""
	return nil
}

func noAppleDouble(conf *mountConfig) error {
	conf.options["noappledouble"] = ""
	return nil
}

func exclCreate(conf *mountConfig) error {
	conf.options["excl_create"] = ""
	return nil
}
