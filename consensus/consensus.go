// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php


// Package consensus implements different Ecosystem consensus engines.
package consensus

import (
	"math/big"

	"github.com/ecosystem/go-ecosystem/common"
	"github.com/ecosystem/go-ecosystem/core/state"
	"github.com/ecosystem/go-ecosystem/core/types"
	"github.com/ecosystem/go-ecosystem/mc"
	"github.com/ecosystem/go-ecosystem/params"
	"github.com/ecosystem/go-ecosystem/rpc"
)

// ChainReader defines a small collection of methods needed to access the local
// blockchain during header and/or uncle verification.
type ChainReader interface {
	// Config retrieves the blockchain's chain configuration.
	Config() *params.ChainConfig

	// CurrentHeader retrieves the current header from the local chain.
	CurrentHeader() *types.Header

	// GetHeader retrieves a block header from the database by hash and number.
	GetHeader(hash common.Hash, number uint64) *types.Header

	// GetHeaderByNumber retrieves a block header from the database by number.
	GetHeaderByNumber(number uint64) *types.Header

	// GetHeaderByHash retrieves a block header from the database by its hash.
	GetHeaderByHash(hash common.Hash) *types.Header

	// GetBlock retrieves a block from the database by hash and number.
	GetBlock(hash common.Hash, number uint64) *types.Block
}

// Engine is an algorithm agnostic consensus engine.
type Engine interface {
	// Author retrieves the Ecosystem address of the account that minted the given
	// block, which may be different from the header's coinbase if a consensus
	// engine is based on signatures.
	Author(header *types.Header) (common.Address, error)

	// VerifyHeader checks whether a header conforms to the consensus rules of a
	// given engine. Verifying the seal may be done optionally here, or explicitly
	// via the VerifySeal method.
	VerifyHeader(chain ChainReader, header *types.Header, seal bool) error

	// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
	// concurrently. The method returns a quit channel to abort the operations and
	// a results channel to retrieve the async verifications (the order is that of
	// the input slice).
	VerifyHeaders(chain ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error)

	// VerifyUncles verifies that the given block's uncles conform to the consensus
	// rules of a given engine.
	VerifyUncles(chain ChainReader, block *types.Block) error

	// VerifySeal checks whether the crypto seal on a header is valid according to
	// the consensus rules of the given engine.
	VerifySeal(chain ChainReader, header *types.Header) error

	// Prepare initializes the consensus fields of a block header according to the
	// rules of a particular engine. The changes are executed inline.
	Prepare(chain ChainReader, header *types.Header) error

	// Finalize runs any post-transaction state modifications (e.g. block rewards)
	// and assembles the final block.
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	Finalize(chain ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error)

	// Seal generates a new block for the given input block with the local miner's
	// seal place on top.
	Seal(chain ChainReader, header *types.Header, stop <-chan struct{}, isBroadcastNode bool) (*types.Header, error)

	// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
	// that a new block should have.
	CalcDifficulty(chain ChainReader, time uint64, parent *types.Header) *big.Int

	// APIs returns the RPC APIs this consensus engine provides.
	APIs(chain ChainReader) []rpc.API
}

// PoW is a consensus engine based on proof-of-work.
type PoW interface {
	Engine

	// Hashrate returns the current mining hashrate of a PoW consensus engine.
	Hashrate() float64
}

type ValidatorReader interface {
	// GetHeaderByNumber retrieves a block header from the database by number.
	GetCurrentNumber() uint64
	GetValidatorByNumber(number uint64) (*mc.TopologyGraph, error)
}

type DPOSEngine interface {
	VerifyBlock(reader ValidatorReader, header *types.Header) error

	VerifyBlocks(reader ValidatorReader, headers []*types.Header) error

	//verify hash in current block
	VerifyHash(reader ValidatorReader, signHash common.Hash, signs []common.Signature) ([]common.Signature, error)

	//verify hash in given number block
	VerifyHashWithNumber(reader ValidatorReader, signHash common.Hash, signs []common.Signature, number uint64) ([]common.Signature, error)

	//VerifyHashWithStocks(signHash common.Hash, signs []common.Signature, stocks map[common.Address]uint16) ([]common.Signature, error)

	VerifyHashWithVerifiedSigns(reader ValidatorReader, signs []*common.VerifiedSign) ([]common.Signature, error)

	VerifyHashWithVerifiedSignsAndNumber(reader ValidatorReader, signs []*common.VerifiedSign, number uint64) ([]common.Signature, error)
}
