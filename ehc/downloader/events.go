// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php


package downloader

type DoneEvent struct{}
type StartEvent struct{}
type FailedEvent struct{ Err error }
