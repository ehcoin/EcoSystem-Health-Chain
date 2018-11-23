// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package reelection

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ecosystem/go-ecosystem/log"
	"github.com/ecosystem/go-ecosystem/random"

	"github.com/ecosystem/go-ecosystem/common"
	"github.com/ecosystem/go-ecosystem/mc"

	"github.com/ecosystem/go-ecosystem/ehc"
)

//func Post() {
//	blockNum := 20
//	for {
//
//		err := mc.PostEvent("CA_RoleUpdated", mc.RoleUpdatedMsg{Role: common.RoleValidator, BlockNum: uint64(blockNum)})
//		blockNum++
//		//fmt.Println("CA_RoleUpdated", mc.RoleUpdatedMsg{Role: common.RoleValidator, BlockNum: uint64(blockNum)})
//		log.Info("err", err)
//		time.Sleep(5 * time.Second)
//
//	}
//}
//
//func TestReElect(t *testing.T) {
//
//	electseed, err := random.NewElectionSeed()
//
//	log.Info("electseed", electseed)
//	log.Info("seed err", err)
//
//	var ehc *ehc.Ecosystem
//	reElect, err := New(ehc)
//	log.Info("err", err)
//
//	go Post()
//
//	time.Sleep(10000 * time.Second)
//	time.Sleep(3 * time.Second)
//	ans1, ans2, ans3 := reElect.readElectData(common.RoleMiner, 240)
//	fmt.Println("READ ELECT", ans1, ans2, ans3)
//	fmt.Println("READ ELECT", 240)
//
//	fmt.Println(reElect)
//}

func TestT(t *testing.T) {
	ans := big.NewInt(100)
	ans1 := common.BigToHash(ans)
	fmt.Println(ans1)

}
func TestCase(t *testing.T) {
	ans1, ans2 := GetAllElectedByHeight(big.NewInt(100), common.RoleMiner)
	fmt.Println(ans1, ans2)
}
