// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php


package ehcclient

import "github.com/ecosystem/go-ecosystem"

// Verify that Client implements the ecosystem interfaces.
var (
	_ = ecosystem.ChainReader(&Client{})
	_ = ecosystem.TransactionReader(&Client{})
	_ = ecosystem.ChainStateReader(&Client{})
	_ = ecosystem.ChainSyncReader(&Client{})
	_ = ecosystem.ContractCaller(&Client{})
	_ = ecosystem.GasEstimator(&Client{})
	_ = ecosystem.GasPricer(&Client{})
	_ = ecosystem.LogFilterer(&Client{})
	_ = ecosystem.PendingStateReader(&Client{})
	// _ = ecosystem.PendingStateEventer(&Client{})
	_ = ecosystem.PendingContractCaller(&Client{})
)
