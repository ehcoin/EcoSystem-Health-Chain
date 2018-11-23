// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php


// Contains the metrics collected by the downloader.

package downloader

import (
	"github.com/ecosystem/go-ecosystem/metrics"
)

var (
	headerInMeter      = metrics.NewRegisteredMeter("ehc/downloader/headers/in", nil)
	headerReqTimer     = metrics.NewRegisteredTimer("ehc/downloader/headers/req", nil)
	headerDropMeter    = metrics.NewRegisteredMeter("ehc/downloader/headers/drop", nil)
	headerTimeoutMeter = metrics.NewRegisteredMeter("ehc/downloader/headers/timeout", nil)

	bodyInMeter      = metrics.NewRegisteredMeter("ehc/downloader/bodies/in", nil)
	bodyReqTimer     = metrics.NewRegisteredTimer("ehc/downloader/bodies/req", nil)
	bodyDropMeter    = metrics.NewRegisteredMeter("ehc/downloader/bodies/drop", nil)
	bodyTimeoutMeter = metrics.NewRegisteredMeter("ehc/downloader/bodies/timeout", nil)

	receiptInMeter      = metrics.NewRegisteredMeter("ehc/downloader/receipts/in", nil)
	receiptReqTimer     = metrics.NewRegisteredTimer("ehc/downloader/receipts/req", nil)
	receiptDropMeter    = metrics.NewRegisteredMeter("ehc/downloader/receipts/drop", nil)
	receiptTimeoutMeter = metrics.NewRegisteredMeter("ehc/downloader/receipts/timeout", nil)

	stateInMeter   = metrics.NewRegisteredMeter("ehc/downloader/states/in", nil)
	stateDropMeter = metrics.NewRegisteredMeter("ehc/downloader/states/drop", nil)
)
