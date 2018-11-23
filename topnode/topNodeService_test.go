// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package topnode

import (
	"github.com/ecosystem/go-ecosystem/accounts/keystore"
	"github.com/ecosystem/go-ecosystem/common"
	"github.com/ecosystem/go-ecosystem/consensus/mtxdpos"
	"github.com/ecosystem/go-ecosystem/core/types"
	"github.com/ecosystem/go-ecosystem/crypto"
	"github.com/ecosystem/go-ecosystem/event"
	"github.com/ecosystem/go-ecosystem/log"
	"github.com/ecosystem/go-ecosystem/mc"
	"github.com/pborman/uuid"
	"reflect"
	"testing"
	"time"
)

var (
	testServs  []testNodeService
	fullstate  = [30]uint8{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	offState   = [30]uint8{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0}
	dposStocks = make(map[common.Address]uint16)
	nodeInfo   = make([]NodeOnLineInfo, 11)
)

type testDPOSEngine struct {
	dops *mtxdpos.MtxDPOS
}

func (tsdpos *testDPOSEngine) VerifyBlock(header *types.Header) error {
	return tsdpos.dops.VerifyBlock(header)
}

//verify hash in current block
func (tsdpos *testDPOSEngine) VerifyHash(signHash common.Hash, signs []common.Signature) ([]common.Signature, error) {
	return tsdpos.dops.VerifyHash(signHash, signs)
}

//verify hash in given number block
func (tsdpos *testDPOSEngine) VerifyHashWithNumber(signHash common.Hash, signs []common.Signature, number uint64) ([]common.Signature, error) {
	return tsdpos.dops.VerifyHashWithStocks(signHash, signs, dposStocks)
}

//VerifyHashWithStocks(signHash common.Hash, signs []common.Signature, stocks map[common.Address]uint16) ([]common.Signature, error)

func (tsdpos *testDPOSEngine) VerifyHashWithVerifiedSigns(signs []*common.VerifiedSign) ([]common.Signature, error) {
	return tsdpos.dops.VerifyHashWithVerifiedSigns(signs)
}

func (tsdpos *testDPOSEngine) VerifyHashWithVerifiedSignsAndNumber(signs []*common.VerifiedSign, number uint64) ([]common.Signature, error) {
	return tsdpos.dops.VerifyHashWithVerifiedSignsAndNumber(signs, number)
}

type Center struct {
	FeedMap map[mc.EventCode]*event.Feed
}

func newCenter() *Center {
	msgCenter := &Center{FeedMap: make(map[mc.EventCode]*event.Feed)}
	for i := 0; i < int(mc.LastEventCode); i++ {
		msgCenter.FeedMap[mc.EventCode(i)] = new(event.Feed)
	}
	return msgCenter
}
func (cen *Center) SubscribeEvent(aim mc.EventCode, ch interface{}) (event.Subscription, error) {
	feed, ok := cen.FeedMap[aim]
	if !ok {
		return nil, mc.SubErrorNoThisEvent
	}
	return feed.Subscribe(ch), nil
}

func (cen *Center) PublishEvent(aim mc.EventCode, data interface{}) error {
	feed, ok := cen.FeedMap[aim]
	if !ok {
		return mc.PostErrorNoThisEvent
	}
	feed.Send(data)
	return nil
}

type testNodeState struct {
	self keystore.Key
}

func newTestNodeState() *testNodeState {
	key, _ := crypto.GenerateKey()
	id := uuid.NewRandom()
	keystore := keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(key.PublicKey),
		PrivateKey: key,
	}
	return &testNodeState{keystore}
}
func (ts *testNodeState) GetTopNodeOnlineState() []NodeOnLineInfo {

	return nodeInfo
}
func (ts *testNodeState) SendNodeMsg(msg interface{}, dstRole int, address []common.Address) {
	switch msg.(type) {
	case *mc.OnlineConsensusReqs:
		data := msg.(*mc.OnlineConsensusReqs)
		for i := 0; i < len(data.ReqList); i++ {
			data.ReqList[i].Leader = ts.self.Address
		}
		for _, serv := range testServs {
			serv.msgChan <- msg
		}

		//		serv.TN.msgCenter.PublishEvent(mc.HD_TopNodeConsensusReq,data.(*mc.OnlineConsensusReqs))
	case *mc.HD_OnlineConsensusVotes:
		data := msg.(*mc.HD_OnlineConsensusVotes)
		for i := 0; i < len(data.Vote); i++ {
			data.Vote[i].FromAccount = ts.self.Address
		}
		for _, serv := range testServs {
			serv.msgChan <- msg
		}

		//		testServs[1].msgChan <-msg
		//		serv.TN.msgCenter.PublishEvent(mc.HD_TopNodeConsensusVote,data.(*mc.HD_OnlineConsensusVotes))
	default:
		for _, serv := range testServs {
			serv.msgChan <- msg
		}
		//		log.Error("Type Error","type",reflect.TypeOf(data))
	}
	//	for _,serv := range testServs{
	//		serv.msgChan <-msg
	//	}
}

func (ts *testNodeState) SignWithValidate(hash []byte, validate bool) (sig []byte, err error) {
	return crypto.SignWithValidate(hash, validate, ts.self.PrivateKey)
}
func (ts *testNodeState) IsSelfAddress(addr common.Address) bool {
	return ts.self.Address == addr
}

type testNodeService struct {
	TN       *TopNodeService
	msgChan  chan interface{}
	testInfo *testNodeState
}

func (serv *testNodeService) getMessageLoop() {
	for {
		select {
		case data := <-serv.msgChan:
			switch data.(type) {
			case *mc.LeaderChangeNotify:
				serv.TN.msgCenter.PublishEvent(mc.Leader_LeaderChangeNotify, data.(*mc.LeaderChangeNotify))
			case *mc.OnlineConsensusReqs:
				serv.TN.msgCenter.PublishEvent(mc.HD_TopNodeConsensusReq, data.(*mc.OnlineConsensusReqs))
			case *mc.HD_OnlineConsensusVotes:
				serv.TN.msgCenter.PublishEvent(mc.HD_TopNodeConsensusVote, data.(*mc.HD_OnlineConsensusVotes))
			default:
				log.Error("Type Error", "type", reflect.TypeOf(data))
			}
		}
	}
}
func newTestNodeService(testInfo *testNodeState) *TopNodeService {
	testServ := NewTopNodeService()
	testServ.topNodeState = testInfo
	testServ.validatorSign = testInfo
	testServ.msgSender = testInfo
	testServ.msgCenter = newCenter()
	testServ.cd = &testDPOSEngine{mtxdpos.NewMtxDPOS(nil)}

	testServ.Start()

	return testServ

}
func newTestServer() {

	testServs = make([]testNodeService, 11)
	nodes := make([]common.Address, 11)
	for i := 0; i < 11; i++ {
		testServs[i].msgChan = make(chan interface{}, 10)
		testServs[i].testInfo = newTestNodeState()
		testServs[i].TN = newTestNodeService(testServs[i].testInfo)
		nodes[i] = testServs[i].testInfo.self.Address
		dposStocks[nodes[i]] = 1
		go testServs[i].getMessageLoop()
	}
	for i := 0; i < 11; i++ {
		testServs[i].TN.stateMap.setElectNodes(nodes)
	}
	for i := 0; i < 11; i++ {
		nodeInfo[i].Address = testServs[i].testInfo.self.Address
		if i == 9 {
			nodeInfo[i].OnlineState = offState
		} else {
			nodeInfo[i].OnlineState = fullstate
		}
	}

}
func setLeader(index int, number uint64, turn uint8) {
	serv := testServs[index]
	leader := mc.LeaderChangeNotify{
		ConsensusState: true,
		Leader:         serv.testInfo.self.Address,
		Number:         number,
		ReelectTurn:    turn,
	}
	serv.TN.msgSender.SendNodeMsg(&leader, 10, nil)
}
func TestNewTopNodeService(t *testing.T) {
	log.InitLog(1)
	newTestServer()
	go setLeader(0, 1, 0)
	time.Sleep(time.Second * 5)
	for i := 0; i < 11; i++ {
		t.Log(testServs[i].TN.stateMap.finishedProposal.DPosVoteS[0].Proposal)
	}

}
func TestNewTopNodeServiceRound(t *testing.T) {
	log.InitLog(1)
	newTestServer()
	go func() {
		setLeader(0, 1, 0)
		time.Sleep(time.Second)
		nodeInfo[7].OnlineState = offState
		setLeader(1, 2, 0)
	}()
	time.Sleep(time.Second * 5)
	for i := 0; i < 11; i++ {
		t.Log(testServs[i].TN.stateMap.finishedProposal.DPosVoteS[0].Proposal)
		t.Log(testServs[i].TN.stateMap.finishedProposal.DPosVoteS[1].Proposal)
	}

}
