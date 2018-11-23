// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package broadcastTx

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"sync"

	"github.com/ecosystem/go-ecosystem/common"
	"github.com/ecosystem/go-ecosystem/core/types"
	"github.com/ecosystem/go-ecosystem/event"
	"github.com/ecosystem/go-ecosystem/internal/ehcapi"
	"github.com/ecosystem/go-ecosystem/log"
	"github.com/ecosystem/go-ecosystem/mc"
)

const (
	// chainHeadChanSize is the size of channel listening to ChainHeadEvent.
	sendBroadCastCHSize = 10
)

type BroadCast struct {
	ehcBackend ehcapi.Backend

	sendBroadCastCH chan mc.BroadCastEvent
	broadCastSub    event.Subscription
	wg              sync.WaitGroup
}

func NewBroadCast(apiBackEnd ehcapi.Backend) *BroadCast {

	bc := &BroadCast{
		ehcBackend:      apiBackEnd,
		sendBroadCastCH: make(chan mc.BroadCastEvent, sendBroadCastCHSize),
	}
	bc.broadCastSub, _ = mc.SubscribeEvent(mc.SendBroadCastTx, bc.sendBroadCastCH)
	bc.wg.Add(1)
	go bc.loop()
	return bc
}
func (bc *BroadCast) Start() {
	//go bc.loop()
}
func (bc *BroadCast) Stop() {
	bc.broadCastSub.Unsubscribe()
	bc.wg.Wait()
	log.Info("BroadCast Server stopped.--YY")
}

func (bc *BroadCast) loop() {
	defer bc.wg.Done()
	for {
		select {
		case ev := <-bc.sendBroadCastCH:
			bc.sendBroadCastTransaction(ev.Txtyps, ev.Height, ev.Data)
		case <-bc.broadCastSub.Err():
			return
		}
	}
}

//YY 广播交易的接口
func (bc *BroadCast) sendBroadCastTransaction(t string, h *big.Int, data []byte) error {
	currBlockHeight := bc.ehcBackend.CurrentBlock().Number()
	//TODO sunchunfeng test
	if h.Cmp(currBlockHeight) < 0 {
		log.Info("===Send BroadCastTx===", "block height less than 100")
		return errors.New("===Send BroadCastTx===,block height less than 100")
	}
	log.Info("=========YY=========", "sendBroadCastTransaction", data)
	bType := false
	if t == mc.CallTheRoll {
		bType = true
	}
	log.Info("=========YY=========11111", "sendBroadCastTransaction:hhhhhhhh", h)
	h.Quo(h, big.NewInt(int64(common.GetBroadcastInterval())))
	log.Info("=========YY=========22222", "sendBroadCastTransaction:hhhhhhhh", h)
	t += h.String()
	log.Info("=========YY=========33333", "sendBroadCastTransaction:tttttttttt", t)
	tmpData := make(map[string][]byte)
	tmpData[t] = data
	msData, _ := json.Marshal(tmpData)
	var txtype byte
	txtype = byte(1)
	tx := types.NewHeartTransaction(txtype, msData)
	var chainID *big.Int
	if config := bc.ehcBackend.ChainConfig(); config.IsEIP155(currBlockHeight) {
		chainID = config.ChainId
	}
	signed, err := bc.ehcBackend.SignTx(tx, chainID)
	if err != nil {
		log.Info("=========YY=========", "sendBroadCastTransaction:SignTx=", err)
		return err
	}
	err1 := bc.ehcBackend.SendBroadTx(context.Background(), signed, bType)
	log.Info("=========YY=========", "sendBroadCastTransaction:Return=", err1)
	return nil
}
