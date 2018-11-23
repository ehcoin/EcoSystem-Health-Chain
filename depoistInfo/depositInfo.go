// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package depoistInfo

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ecosystem/go-ecosystem/common"
	"github.com/ecosystem/go-ecosystem/common/hexutil"
	"github.com/ecosystem/go-ecosystem/core/state"
	"github.com/ecosystem/go-ecosystem/core/types"
	"github.com/ecosystem/go-ecosystem/core/vm"
	"github.com/ecosystem/go-ecosystem/rpc"
)

type ehcBackend interface {
	StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error)
}

type DepositInfo struct {
	MatrixDeposit vm.MatrixDeposit
	Contract      *vm.Contract
	ehcApi        ehcBackend
	p             vm.PrecompiledContract
}

var depositInfo *DepositInfo

func NewDepositInfo(ehcApi ehcBackend) {
	pre := vm.PrecompiledContractsByzantium[common.BytesToAddress([]byte{10})]
	depositInfo = &DepositInfo{ehcApi: ehcApi, p: pre}
	//getDepositListTest()
}

func AddOnlineTime(stateDB vm.StateDB, address common.Address, ot *big.Int) error {
	return depositInfo.MatrixDeposit.AddOnlineTime(depositInfo.Contract, stateDB, address, ot)
}

func GetOnlineTime(stateDB vm.StateDB, address common.Address) (*big.Int, error) {
	return depositInfo.MatrixDeposit.GetOnlineTime(depositInfo.Contract, stateDB, address), nil
}

func SetOnlineTime(stateDB vm.StateDB, address common.Address, ot *big.Int) error {
	return depositInfo.MatrixDeposit.SetOnlineTime(depositInfo.Contract, stateDB, address, ot)
}

func GetDepositList(tm *big.Int, getDeposit common.RoleType) ([]vm.DepositDetail, error) {
	db, err := getDepositInfo(tm)
	if err != nil {
		return nil, err
	}
	contract := vm.NewContract(vm.AccountRef(common.HexToAddress("1337")), vm.AccountRef(common.BytesToAddress([]byte{10})), big.NewInt(0), 60000)
	var depositList []vm.DepositDetail
	switch getDeposit {
	case common.RoleValidator:
		depositList = depositInfo.MatrixDeposit.GetValidatorDepositList(contract, db)
	case common.RoleMiner:
		depositList = depositInfo.MatrixDeposit.GetMinerDepositList(contract, db)
	}
	return depositList, nil
}

func GetDepositAndWithDrawList(tm *big.Int) ([]vm.DepositDetail, error) {
	db, err := getDepositInfo(tm)
	if err != nil {
		return nil, err
	}
	contract := vm.NewContract(vm.AccountRef(common.HexToAddress("1337")), vm.AccountRef(common.BytesToAddress([]byte{10})), big.NewInt(0), 60000)
	var depositList []vm.DepositDetail
	depositList = depositInfo.MatrixDeposit.GetAllDepositList(contract, db, true)
	return depositList, nil
}
func GetAllDeposit(tm *big.Int) ([]vm.DepositDetail, error) {
	db, err := getDepositInfo(tm)
	if err != nil {
		return nil, err
	}
	contract := vm.NewContract(vm.AccountRef(common.HexToAddress("1337")), vm.AccountRef(common.BytesToAddress([]byte{10})), big.NewInt(0), 60000)
	var depositList []vm.DepositDetail
	depositList = depositInfo.MatrixDeposit.GetAllDepositList(contract, db, false)
	return depositList, nil
}

func getDepositInfo(tm *big.Int) (db vm.StateDB, err error) {
	depositInfo.Contract = vm.NewContract(vm.AccountRef(common.HexToAddress("1337")), vm.AccountRef(common.BytesToAddress([]byte{10})), big.NewInt(0), 0)
	var c context.Context
	var h rpc.BlockNumber
	encode := hexutil.EncodeBig(tm)
	err = h.UnmarshalJSON([]byte(encode))
	if err != nil {
		return nil, err
	}
	db, _, err = depositInfo.ehcApi.StateAndHeaderByNumber(c, h)
	return db, err
}

func getDepositListTest() {
	db, err := getDepositInfo(big.NewInt(0))
	if err != nil {
		return
	}
	contract := vm.NewContract(vm.AccountRef(common.HexToAddress("1337")), vm.AccountRef(common.BytesToAddress([]byte{10})), big.NewInt(0), 60000)
	address := depositInfo.MatrixDeposit.GetValidatorDepositList(contract, db)
	fmt.Printf("get depositList:%v %d\n", address, len(address))
	address = depositInfo.MatrixDeposit.GetMinerDepositList(contract, db)
	fmt.Printf("get miner:%v %d\n", address, len(address))
}
