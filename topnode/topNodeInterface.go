// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package topnode

import (
	"github.com/ecosystem/go-ecosystem/accounts/signhelper"
	"github.com/ecosystem/go-ecosystem/ca"
	"github.com/ecosystem/go-ecosystem/common"
	"github.com/ecosystem/go-ecosystem/event"
	"github.com/ecosystem/go-ecosystem/hd"
	"github.com/ecosystem/go-ecosystem/mc"
)

type OnlineState uint8

const (
	Online OnlineState = iota + 1
	Offline
)

func (o OnlineState) String() string {
	switch o {
	case Online:
		return "在线"
	case Offline:
		return "下线"
	default:
		return "未知状态"
	}
}

type NodeOnLineInfo struct {
	Address     common.Address
	OnlineState [30]uint8
}

type TopNodeStateInterface interface {
	GetTopNodeOnlineState() []NodeOnLineInfo
}

type ValidatorAccountInterface interface {
	SignWithValidate(hash []byte, validate bool) (sig common.Signature, err error)
	IsSelfAddress(addr common.Address) bool
}

type MessageSendInterface interface {
	SendNodeMsg(subCode mc.EventCode, msg interface{}, Roles common.RoleType, address []common.Address)
}

type MessageCenterInterface interface {
	SubscribeEvent(aim mc.EventCode, ch interface{}) (event.Subscription, error)
	PublishEvent(aim mc.EventCode, data interface{}) error
}

////////////////////////////////////////////////////////////////////
type TopNodeInstance struct {
	signHelper *signhelper.SignHelper
	hd         *hd.HD
}

func NewTopNodeInstance(sh *signhelper.SignHelper, hd *hd.HD) *TopNodeInstance {
	return &TopNodeInstance{
		signHelper: sh,
		hd:         hd,
	}
}

func (self *TopNodeInstance) GetTopNodeOnlineState() []NodeOnLineInfo {
	return nil
}

func (self *TopNodeInstance) SignWithValidate(hash []byte, validate bool) (sig common.Signature, err error) {
	return self.signHelper.SignHashWithValidate(hash, validate)
}

func (self *TopNodeInstance) IsSelfAddress(addr common.Address) bool {
	return ca.GetAddress() == addr
}

func (self *TopNodeInstance) SendNodeMsg(subCode mc.EventCode, msg interface{}, Roles common.RoleType, address []common.Address) {
	self.hd.SendNodeMsg(subCode, msg, Roles, address)
}

func (self *TopNodeInstance) SubscribeEvent(aim mc.EventCode, ch interface{}) (event.Subscription, error) {
	return mc.SubscribeEvent(aim, ch)
}

func (self *TopNodeInstance) PublishEvent(aim mc.EventCode, data interface{}) error {
	return mc.PublishEvent(aim, data)
}
