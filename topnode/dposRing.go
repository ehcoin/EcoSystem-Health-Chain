// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php

package olconsensus

import (
	"sync"

	"github.com/ecosystem/go-ecosystem/common"
	"github.com/ecosystem/go-ecosystem/log"
	"github.com/ecosystem/go-ecosystem/mc"
	"github.com/ecosystem/go-ecosystem/messageState"
)

type voteInfo struct {
	hash uint64
	data *mc.HD_ConsensusVote
}
type DPosVoteState struct {
	mu               sync.RWMutex
	Hash             common.Hash
	Proposal         interface{}
	AffirmativeVotes []voteInfo
}

func (ds *DPosVoteState) hasHash(hash common.Hash) bool {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.Hash == hash
}
func (ds *DPosVoteState) addProposal(proposal interface{}) bool {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	have := ds.Proposal != nil
	ds.Proposal = proposal
	return have
}
func (ds *DPosVoteState) addVote(vote *mc.HD_ConsensusVote) (interface{}, []voteInfo) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	insVote := voteInfo{messageState.RlpFnvHash(vote), vote}
	log.Info("投票hash", "fnvHash", insVote.hash, "From", vote.From.String())
	for _, item := range ds.AffirmativeVotes {
		if item.hash == insVote.hash {
			log.Error("topnodeOnline", "添加投票,投票已经存在 vote", vote, "已经收到的票数", len(ds.AffirmativeVotes))
			return nil, nil
		}
	}
	ds.AffirmativeVotes = append(ds.AffirmativeVotes, insVote)
	log.Info("topnodeOnline", "添加投票length", len(ds.AffirmativeVotes), "proposal", ds.Proposal)
	return ds.Proposal, ds.AffirmativeVotes[:]
}
func (ds *DPosVoteState) clear(proposal common.Hash) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.Hash = proposal
	ds.Proposal = nil
	ds.AffirmativeVotes = make([]voteInfo, 0)
}
func (ds *DPosVoteState) getVotes() (interface{}, []voteInfo) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.Proposal, ds.AffirmativeVotes[:]
}

type DPosVoteRing struct {
	DPosVoteS []*DPosVoteState
	capacity  int
	mu        sync.RWMutex
	last      int
}

func NewDPosVoteRing(capacity int) *DPosVoteRing {
	ring := &DPosVoteRing{
		DPosVoteS: make([]*DPosVoteState, capacity),
		capacity:  capacity,
		mu:        sync.RWMutex{},
		last:      capacity - 1,
	}
	for i := 0; i < capacity; i++ {
		ring.DPosVoteS[i] = &DPosVoteState{}
	}
	return ring
}
func (ring *DPosVoteRing) insertLast() int {
	ring.mu.Lock()
	defer ring.mu.Unlock()
	ring.last = (ring.last + 1) % ring.capacity
	return ring.last
}
func (ring *DPosVoteRing) getLast() int {
	ring.mu.RLock()
	defer ring.mu.RUnlock()
	return ring.last
}
func (ring *DPosVoteRing) insertNewProposal(hash common.Hash) *DPosVoteState {
	last := ring.insertLast()
	ring.DPosVoteS[last].clear(hash)
	return ring.DPosVoteS[last]
}

func (ring *DPosVoteRing) getVotes(hash common.Hash) (interface{}, []voteInfo) {
	ds, _ := ring.findProposal(hash)
	return ds.getVotes()
}
func (ring *DPosVoteRing) addProposal(hash common.Hash, proposal interface{}) bool {
	ds, have := ring.findProposal(hash)
	add := ds.addProposal(proposal)
	log.Info("topnodeOnline", "have", have, "add", add, "hash", hash, "ds.proposal", ds.Proposal, "proposal", proposal)
	return !(have && add)
}
func (ring *DPosVoteRing) addVote(hash common.Hash, vote *mc.HD_ConsensusVote) (interface{}, []voteInfo) {
	ds, ret := ring.findProposal(hash)
	log.Info("topnodeOnline", "找到hash", ret)
	return ds.addVote(vote)
}
func (ring *DPosVoteRing) findProposal(hash common.Hash) (*DPosVoteState, bool) {
	for _, item := range ring.DPosVoteS {
		if item.hasHash(hash) {
			return item, true
		}
	}
	return ring.insertNewProposal(hash), false
}
