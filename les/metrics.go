// Copyright (c) 2018 The Ecosystem Authors
// Distributed under the MIT software license, see the accompanying
// file COPYING or or or http://www.opensource.org/licenses/mit-license.php


package les

import (
	"github.com/ecosystem/go-ecosystem/metrics"
	"github.com/ecosystem/go-ecosystem/p2p"
)

var (
	/*	propTxnInPacketsMeter     = metrics.NewMeter("ehc/prop/txns/in/packets")
		propTxnInTrafficMeter     = metrics.NewMeter("ehc/prop/txns/in/traffic")
		propTxnOutPacketsMeter    = metrics.NewMeter("ehc/prop/txns/out/packets")
		propTxnOutTrafficMeter    = metrics.NewMeter("ehc/prop/txns/out/traffic")
		propHashInPacketsMeter    = metrics.NewMeter("ehc/prop/hashes/in/packets")
		propHashInTrafficMeter    = metrics.NewMeter("ehc/prop/hashes/in/traffic")
		propHashOutPacketsMeter   = metrics.NewMeter("ehc/prop/hashes/out/packets")
		propHashOutTrafficMeter   = metrics.NewMeter("ehc/prop/hashes/out/traffic")
		propBlockInPacketsMeter   = metrics.NewMeter("ehc/prop/blocks/in/packets")
		propBlockInTrafficMeter   = metrics.NewMeter("ehc/prop/blocks/in/traffic")
		propBlockOutPacketsMeter  = metrics.NewMeter("ehc/prop/blocks/out/packets")
		propBlockOutTrafficMeter  = metrics.NewMeter("ehc/prop/blocks/out/traffic")
		reqHashInPacketsMeter     = metrics.NewMeter("ehc/req/hashes/in/packets")
		reqHashInTrafficMeter     = metrics.NewMeter("ehc/req/hashes/in/traffic")
		reqHashOutPacketsMeter    = metrics.NewMeter("ehc/req/hashes/out/packets")
		reqHashOutTrafficMeter    = metrics.NewMeter("ehc/req/hashes/out/traffic")
		reqBlockInPacketsMeter    = metrics.NewMeter("ehc/req/blocks/in/packets")
		reqBlockInTrafficMeter    = metrics.NewMeter("ehc/req/blocks/in/traffic")
		reqBlockOutPacketsMeter   = metrics.NewMeter("ehc/req/blocks/out/packets")
		reqBlockOutTrafficMeter   = metrics.NewMeter("ehc/req/blocks/out/traffic")
		reqHeaderInPacketsMeter   = metrics.NewMeter("ehc/req/headers/in/packets")
		reqHeaderInTrafficMeter   = metrics.NewMeter("ehc/req/headers/in/traffic")
		reqHeaderOutPacketsMeter  = metrics.NewMeter("ehc/req/headers/out/packets")
		reqHeaderOutTrafficMeter  = metrics.NewMeter("ehc/req/headers/out/traffic")
		reqBodyInPacketsMeter     = metrics.NewMeter("ehc/req/bodies/in/packets")
		reqBodyInTrafficMeter     = metrics.NewMeter("ehc/req/bodies/in/traffic")
		reqBodyOutPacketsMeter    = metrics.NewMeter("ehc/req/bodies/out/packets")
		reqBodyOutTrafficMeter    = metrics.NewMeter("ehc/req/bodies/out/traffic")
		reqStateInPacketsMeter    = metrics.NewMeter("ehc/req/states/in/packets")
		reqStateInTrafficMeter    = metrics.NewMeter("ehc/req/states/in/traffic")
		reqStateOutPacketsMeter   = metrics.NewMeter("ehc/req/states/out/packets")
		reqStateOutTrafficMeter   = metrics.NewMeter("ehc/req/states/out/traffic")
		reqReceiptInPacketsMeter  = metrics.NewMeter("ehc/req/receipts/in/packets")
		reqReceiptInTrafficMeter  = metrics.NewMeter("ehc/req/receipts/in/traffic")
		reqReceiptOutPacketsMeter = metrics.NewMeter("ehc/req/receipts/out/packets")
		reqReceiptOutTrafficMeter = metrics.NewMeter("ehc/req/receipts/out/traffic")*/
	miscInPacketsMeter  = metrics.NewRegisteredMeter("les/misc/in/packets", nil)
	miscInTrafficMeter  = metrics.NewRegisteredMeter("les/misc/in/traffic", nil)
	miscOutPacketsMeter = metrics.NewRegisteredMeter("les/misc/out/packets", nil)
	miscOutTrafficMeter = metrics.NewRegisteredMeter("les/misc/out/traffic", nil)
)

// meteredMsgReadWriter is a wrapper around a p2p.MsgReadWriter, capable of
// accumulating the above defined metrics based on the data stream contents.
type meteredMsgReadWriter struct {
	p2p.MsgReadWriter     // Wrapped message stream to meter
	version           int // Protocol version to select correct meters
}

// newMeteredMsgWriter wraps a p2p MsgReadWriter with metering support. If the
// metrics system is disabled, this function returns the original object.
func newMeteredMsgWriter(rw p2p.MsgReadWriter) p2p.MsgReadWriter {
	if !metrics.Enabled {
		return rw
	}
	return &meteredMsgReadWriter{MsgReadWriter: rw}
}

// Init sets the protocol version used by the stream to know which meters to
// increment in case of overlapping message ids between protocol versions.
func (rw *meteredMsgReadWriter) Init(version int) {
	rw.version = version
}

func (rw *meteredMsgReadWriter) ReadMsg() (p2p.Msg, error) {
	// Read the message and short circuit in case of an error
	msg, err := rw.MsgReadWriter.ReadMsg()
	if err != nil {
		return msg, err
	}
	// Account for the data traffic
	packets, traffic := miscInPacketsMeter, miscInTrafficMeter
	packets.Mark(1)
	traffic.Mark(int64(msg.Size))

	return msg, err
}

func (rw *meteredMsgReadWriter) WriteMsg(msg p2p.Msg) error {
	// Account for the data traffic
	packets, traffic := miscOutPacketsMeter, miscOutTrafficMeter
	packets.Mark(1)
	traffic.Mark(int64(msg.Size))

	// Send the packet to the p2p layer
	return rw.MsgReadWriter.WriteMsg(msg)
}
