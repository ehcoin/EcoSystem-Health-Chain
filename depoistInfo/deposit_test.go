// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package depoistInfo

import (
	"fmt"
	"github.com/ecosystem/go-ecosystem/common/hexutil"
	"github.com/ecosystem/go-ecosystem/rpc"
	"math/big"
	"testing"
)

func TestBigInt(t *testing.T) {
	var tm *big.Int
	var h rpc.BlockNumber
	tm = big.NewInt(100)
	encode := hexutil.EncodeBig(tm)
	err := h.UnmarshalJSON([]byte(encode))
	fmt.Println("err", err)
	fmt.Printf("encode:%T   %v\n", encode, encode)
}
