// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package reelection

import (
	"github.com/ecosystem/go-ecosystem/common"

	"github.com/ecosystem/go-ecosystem/params/man"
	"github.com/ecosystem/go-ecosystem/log"
)

func (self *ReElection)boolTopStatus(height uint64,types common.RoleType)bool{
	if _,_,err:=self.readElectData(types,height);err!=nil{
		return false
	}
	return true
}
func (self *ReElection)checkTopGenStatus(height uint64)error{

	if ok:=self.boolTopStatus(common.GetNextReElectionNumber(height),common.RoleMiner);ok==false{
		log.Warn(Module,"矿工拓扑图需要重新算 高度",height)
		if height==0{
			return nil
		}
		ReElectionHeight:=common.GetNextReElectionNumber(height)
		 if err:=self.ToGenMinerTop(ReElectionHeight - man.MinerTopologyGenerateUpTime);err!=nil{
		 	return err
		 }

	}

	if ok:=self.boolTopStatus(common.GetNextReElectionNumber(height),common.RoleValidator);ok==false{
		log.Warn(Module,"验证者拓扑图需要重新算 高度",height)
		if height==0{
			return nil
		}
		ReElectionHeight:=common.GetNextReElectionNumber(height)
		if err:=self.ToGenValidatorTop(ReElectionHeight-man.VerifyTopologyGenerateUpTime);err!=nil{
			return err
		}
	}
	return nil
}