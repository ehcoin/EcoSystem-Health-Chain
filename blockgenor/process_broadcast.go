// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package blockgenor

import (
	"github.com/ecosystem/go-ecosystem/ca"
	"github.com/ecosystem/go-ecosystem/common"
	"github.com/ecosystem/go-ecosystem/crypto"
	"github.com/ecosystem/go-ecosystem/log"
	"github.com/ecosystem/go-ecosystem/matrixwork"
	"github.com/ecosystem/go-ecosystem/mc"
)

func (p *Process) AddBroadcastMinerResult(result *mc.HD_BroadcastMiningRspMsg) {
	if p.preVerifyBroadcastMinerResult(result.BlockMainData) == false {
		log.WARN(p.logExtraInfo(), "预验证广播区块挖矿结果错误", "抛弃该消息")
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// 缓存广播区块挖矿结果
	log.WARN(p.logExtraInfo(), "缓存广播区块挖矿结果成功，高度", p.number)
	p.broadcastRstCache = append(p.broadcastRstCache, result.BlockMainData)

	p.processMinerResultVerify(p.curLeader, true)
}

func (p *Process) preVerifyBroadcastMinerResult(result *mc.BlockData) bool {
	if role, _ := ca.GetAccountOriginalRole(result.Header.Leader, result.Header.Number.Uint64()); common.RoleBroadcast != role {
		log.ERROR(p.logExtraInfo(), "广播挖矿结果不是来自广播节点, role", role.String())
		return false
	}

	if 1 != len(result.Header.Signatures) {
		log.Error(p.logExtraInfo(), "广播挖矿结果非法, 签名列表数量错误", len(result.Header.Signatures))
		return false
	}
	if false == common.IsBroadcastNumber(result.Header.Number.Uint64()) {
		log.Error(p.logExtraInfo(), "广播挖矿结果非法, 不是广播区块高度", result.Header.Number.Uint64())
		return false
	}
	from, validate, err := crypto.VerifySignWithValidate(result.Header.HashNoSignsAndNonce().Bytes(), result.Header.Signatures[0].Bytes())
	if err != nil {
		log.Error(p.logExtraInfo(), "广播挖矿结果非法, 签名解析错误", err)
		return false
	}

	if from != result.Header.Leader {
		log.Error(p.logExtraInfo(), "广播挖矿结果非法, 签名不匹配，签名人", from.Hex(), "Leader", result.Header.Leader.Hex())
		return false
	}

	if false == validate {
		log.Error(p.logExtraInfo(), "广播挖矿结果非法, 签名结果为", validate)
		return false
	}

	return true
}

func (p *Process) dealMinerResultVerifyBroadcast() {
	for _, result := range p.broadcastRstCache {
		//add by hyk, 运行广播区块交易
		parent := p.blockChain().GetBlockByHash(result.Header.ParentHash)
		if parent == nil {
			log.ERROR(p.logExtraInfo(), "广播挖矿结果验证", "获取父区块错误!")
			continue
		}

		work, err := matrixwork.NewWork(p.blockChain().Config(), p.blockChain(), nil, result.Header)
		if err != nil {
			log.ERROR(p.logExtraInfo(), "广播挖矿结果验证, 创建worker错误", err)
			continue
		}

		log.INFO("*********************", "len(result.Txs)", len(result.Txs))
		for _, tx := range result.Txs {
			log.INFO("==========", "Finalize:GasPrice", tx.GasPrice(), "amount", tx.Value()) //hezi
		}
		//执行交易
		work.ProcessBroadcastTransactions(p.pm.matrix.EventMux(), result.Txs, p.pm.bc)
		_, err = p.blockChain().Engine().Finalize(p.blockChain(), result.Header, work.State, result.Txs, nil, work.Receipts)

		if err != nil {
			log.ERROR(p.logExtraInfo(), "Failed to finalize block for sealing", err)
			continue
		}

		p.blockCache.SaveReadyBlock(&mc.BlockLocalVerifyOK{
			Header:    result.Header,
			BlockHash: common.Hash{},
			Txs:       result.Txs,
			Receipts:  work.Receipts,
			State:     work.State,
		})

		readyMsg := &mc.NewBlockReadyMsg{
			Header: result.Header,
		}
		log.INFO(p.logExtraInfo(), "广播区块验证完成", "发送新区块准备完毕消息", "高度", p.number)
		mc.PublishEvent(mc.BlockGenor_NewBlockReady, readyMsg)

		p.changeState(StateBlockInsert)
		p.processBlockInsert()
		return
	}
}
