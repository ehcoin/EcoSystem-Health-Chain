// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php


package les

import (
	"context"
	"math/big"

	"github.com/ecosystem/go-ecosystem/accounts"
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
	"github.com/ecosystem/go-ecosystem/light"
	"github.com/ecosystem/go-ecosystem/params"
	"github.com/ecosystem/go-ecosystem/rpc"
)

type LesApiBackend struct {
	ehc *LightEcosystem
	gpo *gasprice.Oracle
}

func (b *LesApiBackend) ChainConfig() *params.ChainConfig {
	return b.ehc.chainConfig
}

func (b *LesApiBackend) CurrentBlock() *types.Block {
	return types.NewBlockWithHeader(b.ehc.BlockChain().CurrentHeader())
}

func (b *LesApiBackend) SetHead(number uint64) {
	b.ehc.protocolManager.downloader.Cancel()
	b.ehc.blockchain.SetHead(number)
}

func (b *LesApiBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	if blockNr == rpc.LatestBlockNumber || blockNr == rpc.PendingBlockNumber {
		return b.ehc.blockchain.CurrentHeader(), nil
	}

	return b.ehc.blockchain.GetHeaderByNumberOdr(ctx, uint64(blockNr))
}

func (b *LesApiBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, err
	}
	return b.GetBlock(ctx, header.Hash())
}

func (b *LesApiBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	return light.NewState(ctx, header, b.ehc.odr), header, nil
}

func (b *LesApiBackend) GetBlock(ctx context.Context, blockHash common.Hash) (*types.Block, error) {
	return b.ehc.blockchain.GetBlockByHash(ctx, blockHash)
}

func (b *LesApiBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	if number := rawdb.ReadHeaderNumber(b.ehc.chainDb, hash); number != nil {
		return light.GetBlockReceipts(ctx, b.ehc.odr, hash, *number)
	}
	return nil, nil
}

func (b *LesApiBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	if number := rawdb.ReadHeaderNumber(b.ehc.chainDb, hash); number != nil {
		return light.GetBlockLogs(ctx, b.ehc.odr, hash, *number)
	}
	return nil, nil
}

func (b *LesApiBackend) GetTd(hash common.Hash) *big.Int {
	return b.ehc.blockchain.GetTdByHash(hash)
}

func (b *LesApiBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	context := core.NewEVMContext(msg, header, b.ehc.blockchain, nil)
	return vm.NewEVM(context, state, b.ehc.chainConfig, vmCfg), state.Error, nil
}

func (b *LesApiBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.ehc.txPool.Add(ctx, signedTx)
}

func (b *LesApiBackend) RemoveTx(txHash common.Hash) {
	b.ehc.txPool.RemoveTx(txHash)
}

func (b *LesApiBackend) GetPoolTransactions() (types.Transactions, error) {
	return b.ehc.txPool.GetTransactions()
}

func (b *LesApiBackend) GetPoolTransaction(txHash common.Hash) *types.Transaction {
	return b.ehc.txPool.GetTransaction(txHash)
}

func (b *LesApiBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.ehc.txPool.GetNonce(ctx, addr)
}

func (b *LesApiBackend) Stats() (pending int, queued int) {
	return b.ehc.txPool.Stats(), 0
}

func (b *LesApiBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.ehc.txPool.Content()
}

func (b *LesApiBackend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.ehc.txPool.SubscribeNewTxsEvent(ch)
}

func (b *LesApiBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.ehc.blockchain.SubscribeChainEvent(ch)
}

func (b *LesApiBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.ehc.blockchain.SubscribeChainHeadEvent(ch)
}

func (b *LesApiBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.ehc.blockchain.SubscribeChainSideEvent(ch)
}

func (b *LesApiBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.ehc.blockchain.SubscribeLogsEvent(ch)
}

func (b *LesApiBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.ehc.blockchain.SubscribeRemovedLogsEvent(ch)
}

func (b *LesApiBackend) Downloader() *downloader.Downloader {
	return b.ehc.Downloader()
}

func (b *LesApiBackend) ProtocolVersion() int {
	return b.ehc.LesVersion() + 10000
}

func (b *LesApiBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *LesApiBackend) ChainDb() ehcdb.Database {
	return b.ehc.chainDb
}

func (b *LesApiBackend) EventMux() *event.TypeMux {
	return b.ehc.eventMux
}

func (b *LesApiBackend) AccountManager() *accounts.Manager {
	return b.ehc.accountManager
}

func (b *LesApiBackend) BloomStatus() (uint64, uint64) {
	if b.ehc.bloomIndexer == nil {
		return 0, 0
	}
	sections, _, _ := b.ehc.bloomIndexer.Sections()
	return light.BloomTrieFrequency, sections
}

func (b *LesApiBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.ehc.bloomRequests)
	}
}

//YY
func (b *LesApiBackend) SignTx(signedTx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	account := accounts.Account{Address: b.ehc.config.Etherbase}
	return b.AccountManager().Wallets()[0].SignTx(account, signedTx, chainID)
}

//YY
func (b *LesApiBackend) SendBroadTx(ctx context.Context, signedTx *types.Transaction, bType bool) error {
	return nil //b.ehc.txPool.AddBroadTx(signedTx,bType)
}

//YY
func (b *LesApiBackend) FetcherNotify(hash common.Hash, number uint64) {
	//ids := ca.Ide.GetRoleGroup(common.RoleValidator)
	//for _,id := range ids{
	//peer := b.ehc.protocolManager.peers.Peer(id.String())
	//b.ehc.protocolManager.fetcher.Notify(id.String(), hash, number, time, peer.RequestOneHeader, peer.RequestBodies)
	//}
}
