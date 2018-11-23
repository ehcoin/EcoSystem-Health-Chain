// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php
package hd

import (
	"github.com/ecosystem/go-ecosystem/common"
	"github.com/ecosystem/go-ecosystem/core/types"
)

// AlgorithmMsg
type AlgorithmMsg struct {
	Account common.Address
	Data    NetData
}

//NetData
type NetData struct {
	SubCode uint32
	Msg     []byte
}

type fullBlockMsgForMarshal struct {
	Header *types.Header
	Txs    []*types.Transaction_Mx
}
