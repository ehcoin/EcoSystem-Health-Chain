// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php


package ehc

import (
	"context"
	"math/big"
	"time"

	"github.com/ecosystem/go-ecosystem/accounts"
	"github.com/ecosystem/go-ecosystem/ca"
	"github.com/ecosystem/go-ecosystem/common"
	"github.com/ecosystem/go-ecosystem/common/math"
	"github.com/ecosystem/go-ecosystem/core"
	"github.com/ecosystem/go-ecosystem/core/bloombits"
	"github.com/ecosystem/go-ecosystem/core/rawdb"
	"github.com/ecosystem/go-ecosystem/core/state"
	"github.com/ecosystem/go-ecosystem/core/types"
	"github.com/ecosystem/go-ecosystem/core/vm"
	"github.com/ecosystem/go-ecosystem/ehc/downloader"
	"github.com/ecosystem/go-ecosystem/ehc/gasprice"
	"github.com/ecosystem/go-ecosystem/ehcdb"
	"github.com/ecosystem/go-ecosystem/event"
	"github.com/ecosystem/go-ecosystem/log"
	"github.com/ecosystem/go-ecosystem/params"
	"github.com/ecosystem/go-ecosystem/rpc"
)

// EthAPIBackend implements ehcapi.Backend for full nodes
type EthAPIBackend struct {
	ehc *Ecosystem
	gpo *gasprice.Oracle
}

func (b *EthAPIBackend) ChainConfig() *params.ChainConfig {
	return b.ehc.chainConfig
}

func (b *EthAPIBackend) CurrentBlock() *types.Block {
	return b.ehc.blockchain.CurrentBlock()
}

func (b *EthAPIBackend) SetHead(number uint64) {
	b.ehc.protocolManager.downloader.Cancel()
	b.ehc.blockchain.SetHead(number)
}

func (b *EthAPIBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.ehc.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.ehc.blockchain.CurrentBlock().Header(), nil
	}
	return b.ehc.blockchain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b *EthAPIBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.ehc.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.ehc.blockchain.CurrentBlock(), nil
	}
	return b.ehc.blockchain.GetBlockByNumber(uint64(blockNr)), nil
}

func (b *EthAPIBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, state := b.ehc.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.ehc.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

func (b *EthAPIBackend) GetBlock(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.ehc.blockchain.GetBlockByHash(hash), nil
}

func (b *EthAPIBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	if number := rawdb.ReadHeaderNumber(b.ehc.chainDb, hash); number != nil {
		return rawdb.ReadReceipts(b.ehc.chainDb, hash, *number), nil
	}
	return nil, nil
}

func (b *EthAPIBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	number := rawdb.ReadHeaderNumber(b.ehc.chainDb, hash)
	if number == nil {
		return nil, nil
	}
	receipts := rawdb.ReadReceipts(b.ehc.chainDb, hash, *number)
	if receipts == nil {
		return nil, nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (b *EthAPIBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.ehc.blockchain.GetTdByHash(blockHash)
}

func (b *EthAPIBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewEVMContext(msg, header, b.ehc.BlockChain(), nil)
	return vm.NewEVM(context, state, b.ehc.chainConfig, vmCfg), vmError, nil
}

func (b *EthAPIBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.ehc.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *EthAPIBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.ehc.BlockChain().SubscribeChainEvent(ch)
}

func (b *EthAPIBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.ehc.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *EthAPIBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.ehc.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *EthAPIBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.ehc.BlockChain().SubscribeLogsEvent(ch)
}

func (b *EthAPIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.ehc.txPool.AddLocal(signedTx)
}

func (b *EthAPIBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.ehc.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *EthAPIBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.ehc.txPool.Get(hash)
}

func (b *EthAPIBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.ehc.txPool.State().GetNonce(addr), nil
}

func (b *EthAPIBackend) Stats() (pending int, queued int) {
	return b.ehc.txPool.Stats()
}

func (b *EthAPIBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.ehc.TxPool().Content()
}

func (b *EthAPIBackend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.ehc.TxPool().SubscribeNewTxsEvent(ch)
}

func (b *EthAPIBackend) Downloader() *downloader.Downloader {
	return b.ehc.Downloader()
}

func (b *EthAPIBackend) ProtocolVersion() int {
	return b.ehc.EthVersion()
}

func (b *EthAPIBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *EthAPIBackend) ChainDb() ehcdb.Database {
	return b.ehc.ChainDb()
}

func (b *EthAPIBackend) EventMux() *event.TypeMux {
	return b.ehc.EventMux()
}

func (b *EthAPIBackend) AccountManager() *accounts.Manager {
	return b.ehc.AccountManager()
}

func (b *EthAPIBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.ehc.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *EthAPIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.ehc.bloomRequests)
	}
}

//YY
func (b *EthAPIBackend) SignTx(signedTx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return b.ehc.signHelper.SignTx(signedTx, chainID)
}

//YY
func (b *EthAPIBackend) SendBroadTx(ctx context.Context, signedTx *types.Transaction, bType bool) error {
	return b.ehc.txPool.AddBroadTx(signedTx, bType)
}

//YY
func (b *EthAPIBackend) FetcherNotify(hash common.Hash, number uint64) {
	ids := ca.GetRolesByGroup(common.RoleValidator)
	log.Info("==========YY===========", "FetcherNotify()��Validator`s count", len(ids))
	for _, id := range ids {
		peer := b.ehc.protocolManager.Peers.Peer(id.String())
		log.Info("==========YY===========", "FetcherNotify()��Validator`s NodeID", id)
		log.Info("==========YY===========", "FetcherNotify()��get PeerID by Validator ID", peer.id)
		b.ehc.protocolManager.fetcher.Notify(id.String(), hash, number, time.Now(), peer.RequestOneHeader, peer.RequestBodies)
		log.Info("==========YY===========", "FetcherNotify()��send Notify completed", 111111111111111)
	}
}
