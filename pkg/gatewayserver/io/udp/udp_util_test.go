// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package udp_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/smartystreets/assertions"
	"go.thethings.network/lorawan-stack/pkg/gatewayserver/io"
	iotesting "go.thethings.network/lorawan-stack/pkg/gatewayserver/io/testing"
	encoding "go.thethings.network/lorawan-stack/pkg/gatewayserver/io/udp/encoding"
	"go.thethings.network/lorawan-stack/pkg/ttnpb"
	"go.thethings.network/lorawan-stack/pkg/types"
	"go.thethings.network/lorawan-stack/pkg/util/test"
	"go.thethings.network/lorawan-stack/pkg/util/test/assertions/should"
)

func generatePushData(eui types.EUI64, status bool, timestamps ...time.Duration) encoding.Packet {
	packet := encoding.Packet{
		GatewayEUI:      &eui,
		ProtocolVersion: encoding.Version1,
		Token:           [2]byte{0x00, 0x00},
		PacketType:      encoding.PushData,
		Data: &encoding.Data{
			RxPacket: make([]*encoding.RxPacket, len(timestamps)),
		},
	}
	if status {
		packet.Data.Stat = &encoding.Stat{
			Time: encoding.ExpandedTime(time.Now()),
		}
	}
	for i, t := range timestamps {
		up := ttnpb.NewPopulatedUplinkMessage(test.Randy, true)
		var modulation string
		var dataRate types.DataRate
		switch up.Settings.Modulation {
		case ttnpb.Modulation_LORA:
			modulation = "LORA"
			dataRate.LoRa = fmt.Sprintf("SF%dBW%d", up.Settings.SpreadingFactor, up.Settings.Bandwidth/1000)
		case ttnpb.Modulation_FSK:
			modulation = "FSK"
			dataRate.FSK = up.Settings.BitRate
		}
		packet.Data.RxPacket[i] = &encoding.RxPacket{
			Freq: float64(up.Settings.Frequency) / 1000000,
			Chan: uint8(up.Settings.ChannelIndex),
			Modu: modulation,
			DatR: encoding.DataRate{DataRate: dataRate},
			CodR: up.Settings.CodingRate,
			Size: uint16(len(up.RawPayload)),
			Data: base64.StdEncoding.EncodeToString(up.RawPayload),
			Tmst: uint32(t / time.Microsecond),
		}
	}
	return packet
}

func generatePullData(eui types.EUI64) encoding.Packet {
	return encoding.Packet{
		GatewayEUI:      &eui,
		ProtocolVersion: encoding.Version1,
		Token:           [2]byte{0x00, 0x00},
		PacketType:      encoding.PullData,
	}
}

func generateTxAck(eui types.EUI64, err encoding.TxError) encoding.Packet {
	return encoding.Packet{
		GatewayEUI:      &eui,
		ProtocolVersion: encoding.Version1,
		Token:           [2]byte{0x00, 0x00},
		PacketType:      encoding.TxAck,
		Data: &encoding.Data{
			TxPacketAck: &encoding.TxPacketAck{
				Error: err,
			},
		},
	}
}

func expectAck(t *testing.T, conn net.Conn, expect bool, packetType encoding.PacketType, token [2]byte) {
	var buf [65507]byte
	conn.SetReadDeadline(time.Now().Add(timeout))
	n, err := conn.Read(buf[:])
	if err != nil {
		if !expect {
			return
		}
		t.Fatal("Failed to read acknowledgement")
	}
	var ack encoding.Packet
	if err := ack.UnmarshalBinary(buf[0:n]); err != nil {
		t.Fatal("Failed to unmarshal acknowledgement")
	}
	if ack.PacketType != packetType {
		t.Fatalf("Packet type %v is not %v", ack.PacketType, packetType)
	}
	if !bytes.Equal(ack.Token[:], token[:]) {
		t.Fatal("Received acknowledgement with unexpected token")
	}
	if !expect {
		t.Fatal("Should not have received acknowledgement for this token")
	}
}

func expectConnection(t *testing.T, server iotesting.Server, connections *sync.Map, eui types.EUI64, expectNew bool) *io.Connection {
	a := assertions.New(t)
	var conn *io.Connection
	select {
	case conn = <-server.Connections():
		if !a.So(expectNew, should.BeTrue) {
			t.Fatal("Should not have a new connection")
		}
		actual := *conn.Gateway().GatewayIdentifiers.EUI
		if actual != eui {
			t.Fatalf("New connection for unexpected EUI %v", actual)
		}
		if _, loaded := connections.LoadOrStore(eui, conn); loaded {
			t.Fatalf("Gateway %v already connected", eui)
		}
		go func() {
			<-conn.Context().Done()
			connections.Delete(eui)
		}()
	case <-time.After(timeout):
		if !a.So(expectNew, should.BeFalse) {
			t.Fatal("New connection timeout")
		} else if v, loaded := connections.Load(eui); !loaded {
			t.Fatal("Expected existing connection")
		} else {
			conn = v.(*io.Connection)
		}
	}
	return conn
}