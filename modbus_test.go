// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gromb

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strings"
	"testing"
)

func TestRTU(t *testing.T) {
	testPack(t, ProtocolRTU)
	testParse(t, ProtocolRTU)
}

func TestAscii(t *testing.T) {
	testPack(t, ProtocolAscii)
	testParse(t, ProtocolAscii)
}

func TestTCP(t *testing.T) {
	testPack(t, ProtocolTCP)
	testParse(t, ProtocolTCP)
}

type testItem struct {
	funccode uint8    // 功能码
	devid    uint8    // 设备标识
	regaddr  uint16   // 寄存器地址
	reglen   uint16   // 寄存器数量
	arg      []uint16 // 寄存器值
	sernum   uint16   // 流水号
	rtuReq   string   // RTU请求报文
	asciiReq string   // ASCII请求报文
	tcpReq   string   // TCP请求报文
	rtuRsp   string   // RTU响应报文
	asciiRsp string   // ASCII响应报文
	tcpRsp   string   // TCP响应报文
}

func (t *testItem) packet(protocol uint8, isReq bool) []uint8 {
	switch protocol {
	case ProtocolRTU:
		if isReq {
			return strToHex(t.rtuReq)
		} else {
			return strToHex(t.rtuRsp)
		}
	case ProtocolAscii:
		if isReq {
			return strToHex(t.asciiReq)
		} else {
			return strToHex(t.asciiRsp)
		}
	case ProtocolTCP:
		if isReq {
			return strToHex(t.tcpReq)
		} else {
			return strToHex(t.tcpRsp)
		}
	}
	return nil
}

var testItems = []testItem{
	{
		funccode: FuncCodeReadCoil,
		devid:    0x01,
		regaddr:  0x0000,
		reglen:   0x0007,
		arg:      []uint16{1, 0, 1, 1, 0, 0, 1},
		sernum:   0x0001,
		rtuReq:   "01 01 00 00 00 07 7D C8",
		asciiReq: "3A 30 31 30 31 30 30 30 30 30 30 30 37 46 37 0D 0A",
		tcpReq:   "00 01 00 00 00 06 01 01 00 00 00 07 ",
		rtuRsp:   "01 01 01 4D 91 BD 11 11 11 11 11 11",
		asciiRsp: "3A 30 31 30 31 30 31 34 44 42 30 0D 0A 32 32 32 32 32 32 32 32",
		tcpRsp:   "00 01 00 00 00 04 01 01 01 4D 33 33 33 33 33 33",
	},
	{
		funccode: FuncCodeReadCoil,
		devid:    0x01,
		regaddr:  0x0131,
		reglen:   0x0002,
		arg:      []uint16{0, 1},
		sernum:   0x01,
		rtuReq:   "01 01 01 31 00 02 ED F8",
		asciiReq: "3A 30 31 30 31 30 31 33 31 30 30 30 32 43 41 0D 0A",
		tcpReq:   "00 01 00 00 00 06 01 01 01 31 00 02",
		rtuRsp:   "01 01 01 02 D0 49",
		asciiRsp: "3A 30 31 30 31 30 31 30 32 46 42 0D 0A",
		tcpRsp:   "00 01 00 00 00 04 01 01 01 02",
	},
	{
		funccode: FuncCodeReadCoil,
		devid:    0x01,
		regaddr:  0x0101,
		reglen:   0x0009,
		arg:      []uint16{0, 1, 0, 1, 0, 1, 1, 1, 0},
		sernum:   0x01,
		rtuReq:   "01 01 01 01 00 09 AC 30",
		asciiReq: "3A 30 31 30 31 30 31 30 31 30 30 30 39 46 33 0D 0A",
		tcpReq:   "00 01 00 00 00 06 01 01 01 01 00 09 ",
		rtuRsp:   "01 01 02 EA 00 F6 9C ",
		asciiRsp: "3A 30 31 30 31 30 32 45 41 30 30 31 32 0D 0A",
		tcpRsp:   "00 01 00 00 00 05 01 01 02 EA 00",
	},
	{
		funccode: FuncCodeWriteCoil,
		devid:    0x01,
		regaddr:  0x0000,
		reglen:   0x0001,
		arg:      []uint16{1},
		sernum:   0x01,
		rtuReq:   "01 05 00 00 FF 00 8C 3A ",
		asciiReq: "3A 30 31 30 35 30 30 30 30 46 46 30 30 46 42 0D 0A",
		tcpReq:   "00 01 00 00 00 06 01 05 00 00 FF 00 ",
		rtuRsp:   "01 05 00 00 FF 00 8C 3A ",
		asciiRsp: "3A 30 31 30 35 30 30 30 30 46 46 30 30 46 42 0D 0A",
		tcpRsp:   "00 01 00 00 00 06 01 05 00 00 FF 00 ",
	},
	{
		funccode: FuncCodeWriteCoil,
		devid:    0x01,
		regaddr:  0x0000,
		reglen:   0x0001,
		arg:      []uint16{0},
		sernum:   0x01,
		rtuReq:   "01 05 00 00 00 00 CD CA",
		asciiReq: "3A 30 31 30 35 30 30 30 30 30 30 30 30 46 41 0D 0A",
		tcpReq:   "00 01 00 00 00 06 01 05 00 00 00 00 ",
		rtuRsp:   "01 05 00 00 00 00 CD CA",
		asciiRsp: "3A 30 31 30 35 30 30 30 30 30 30 30 30 46 41 0D 0A",
		tcpRsp:   "00 01 00 00 00 06 01 05 00 00 00 00 ",
	},
	{
		funccode: FuncCodeWriteCoil,
		devid:    0x11,
		regaddr:  0x0100,
		reglen:   0x0001,
		arg:      []uint16{1},
		sernum:   0x01,
		rtuReq:   "11 05 01 00 FF 00 8F 56 ",
		asciiReq: "3A 31 31 30 35 30 31 30 30 46 46 30 30 45 41 0D 0A",
		tcpReq:   "00 01 00 00 00 06 11 05 01 00 FF 00 ",
		rtuRsp:   "11 05 01 00 FF 00 8F 56 ",
		asciiRsp: "3A 31 31 30 35 30 31 30 30 46 46 30 30 45 41 0D 0A",
		tcpRsp:   "00 01 00 00 00 06 11 05 01 00 FF 00 ",
	},
	{
		funccode: FuncCodeWriteCoils,
		devid:    0x01,
		regaddr:  0x0000,
		reglen:   0x0007,
		arg:      []uint16{1, 0, 1, 1, 0, 0, 1},
		sernum:   0x01,
		rtuReq:   "01 0F 00 00 00 07 01 4D 0E A3 ",
		asciiReq: "3A 30 31 30 46 30 30 30 30 30 30 30 37 30 31 34 44 39 42 0D 0A",
		tcpReq:   "00 01 00 00 00 08 01 0F 00 00 00 07 01 4D",
		rtuRsp:   "01 0F 00 00 00 07 14 09 ",
		asciiRsp: "3A 30 31 30 46 30 30 30 30 30 30 30 37 45 39 0D 0A",
		tcpRsp:   "00 01 00 00 00 06 01 0F 00 00 00 07",
	},
	{
		funccode: FuncCodeWriteCoils,
		devid:    0x01,
		regaddr:  0x0131,
		reglen:   0x0002,
		arg:      []uint16{0, 1},
		sernum:   0x01,
		rtuReq:   "01 0F 01 31 00 02 01 02 23 43",
		asciiReq: "3A 30 31 30 46 30 31 33 31 30 30 30 32 30 31 30 32 42 39 0D 0A",
		tcpReq:   "00 01 00 00 00 08 01 0F 01 31 00 02 01 02 ",
		rtuRsp:   "01 0F 01 31 00 02 84 39",
		asciiRsp: "3A 30 31 30 46 30 31 33 31 30 30 30 32 42 43 0D 0A",
		tcpRsp:   "00 01 00 00 00 06 01 0F 01 31 00 02",
	},
	{
		funccode: FuncCodeWriteCoils,
		devid:    0x01,
		regaddr:  0x0101,
		reglen:   0x0009,
		arg:      []uint16{0, 1, 0, 1, 0, 1, 1, 1, 0},
		sernum:   0x01,
		rtuReq:   "01 0F 01 01 00 09 02 EA 00 BB 0D ",
		asciiReq: "3A 30 31 30 46 30 31 30 31 30 30 30 39 30 32 45 41 30 30 46 39 0D 0A",
		tcpReq:   "00 01 00 00 00 09 01 0F 01 01 00 09 02 EA 00",
		rtuRsp:   "01 0F 01 01 00 09 C5 F1",
		asciiRsp: "3A 30 31 30 46 30 31 30 31 30 30 30 39 45 35 0D 0A",
		tcpRsp:   "00 01 00 00 00 06 01 0F 01 01 00 09",
	},
	{
		funccode: FuncCodeReadDiscrete,
		devid:    0x01,
		regaddr:  0x0000,
		reglen:   0x000A,
		arg:      []uint16{0, 1, 1, 1, 0, 1, 0, 1, 1, 1},
		sernum:   0x01,
		rtuReq:   "01 02 00 00 00 0A F8 0D ",
		asciiReq: "3A 30 31 30 32 30 30 30 30 30 30 30 41 46 33 0D 0A",
		tcpReq:   "00 01 00 00 00 06 01 02 00 00 00 0A",
		rtuRsp:   "01 02 02 AE 03 85 D9 ",
		asciiRsp: "3A 30 31 30 32 30 32 41 45 30 33 34 41 0D 0A",
		tcpRsp:   "00 01 00 00 00 05 01 02 02 AE 03",
	},
	{
		funccode: FuncCodeReadDiscrete,
		devid:    0x01,
		regaddr:  0x0000,
		reglen:   0x0007,
		arg:      []uint16{0, 1, 1, 1, 0, 1, 1},
		sernum:   0x01,
		rtuReq:   "01 02 00 00 00 07 39 C8",
		asciiReq: "3A 30 31 30 32 30 30 30 30 30 30 30 37 46 36 0D 0A",
		tcpReq:   "00 01 00 00 00 06 01 02 00 00 00 07 ",
		rtuRsp:   "01 02 01 6E 20 64",
		asciiRsp: "3A 30 31 30 32 30 31 36 45 38 45 0D 0A",
		tcpRsp:   "00 01 00 00 00 04 01 02 01 6E",
	},
	{
		funccode: FuncCodeReadHold,
		devid:    0x01,
		regaddr:  0x0000,
		reglen:   0x000A,
		arg:      []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		sernum:   0x01,
		rtuReq:   "01030000000AC5CD",
		asciiReq: "3A 30 31 30 33 30 30 30 30 30 30 30 41 46 32 0D 0A",
		tcpReq:   "00010000000601030000000A",
		rtuRsp:   "010314000100020003000400050006000700080009000A8F16",
		asciiRsp: "3A 30 31 30 33 31 34 30 30 30 31 30 30 30 32 30 30 30 33 30 30 30 34 30 30 30 35 30 30 30" +
			"36 30 30 30 37 30 30 30 38 30 30 30 39 30 30 30 41 42 31 0D 0A",
		tcpRsp: "000100000017010314000100020003000400050006000700080009000A",
	},
	{
		funccode: FuncCodeReadHold,
		devid:    0x01,
		regaddr:  0x0012,
		reglen:   0x0022,
		arg: []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7,
			8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4},
		sernum:   0x01,
		rtuReq:   "01030012002265D6",
		asciiReq: "3A 30 31 30 33 30 30 31 32 30 30 32 32 43 38 0D 0A",
		tcpReq:   "000100000006010300120022",
		rtuRsp: "010344000100020003000400050006000700080009000A000100020003000400050006000700080009000A000" +
			"100020003000400050006000700080009000A000100020003000477A7",
		asciiRsp: "3A 30 31 30 33 34 34 30 30 30 31 30 30 30 32 30 30 30 33 30 30 30 34 30 30 30 35 30 30 30 " +
			"36 30 30 30 37 30 30 30 38 30 30 30 39 30 30 30 41 30 30 30 31 30 30 30 32 30 30 30 33 30 " +
			"30 30 34 30 30 30 35 30 30 30 36 30 30 30 37 30 30 30 38 30 30 30 39 30 30 30 41 30 30 30 " +
			"31 30 30 30 32 30 30 30 33 30 30 30 34 30 30 30 35 30 30 30 36 30 30 30 37 30 30 30 38 30 " +
			"30 30 39 30 30 30 41 30 30 30 31 30 30 30 32 30 30 30 33 30 30 30 34 30 39 0D 0A",
		tcpRsp: "000100000047010344000100020003000400050006000700080009000A00010002000300040005000600070008000" +
			"9000A000100020003000400050006000700080009000A0001000200030004",
	},
	{
		funccode: FuncCodeReadHold,
		devid:    0x56,
		regaddr:  0x0087,
		reglen:   0x0001,
		arg:      []uint16{0x69},
		sernum:   0x01,
		rtuReq:   "56030087000139C4",
		asciiReq: "3A 35 36 30 33 30 30 38 37 30 30 30 31 31 46 0D 0A",
		tcpReq:   "000100000006560300870001",
		rtuRsp:   "56030200690DA6",
		asciiRsp: "3A 35 36 30 33 30 32 30 30 36 39 33 43 0D 0A",
		tcpRsp:   "0001000000055603020069",
	},
	{
		funccode: FuncCodeWriteHold,
		devid:    0x01,
		regaddr:  0x0000,
		reglen:   0x0001,
		arg:      []uint16{0x007B},
		sernum:   0x01,
		rtuReq:   "01060000007BC9E9",
		asciiReq: "3A 30 31 30 36 30 30 30 30 30 30 37 42 37 45 0D 0A",
		tcpReq:   "00010000000601060000007B",
		rtuRsp:   "01060000007BC9E9",
		asciiRsp: "3A 30 31 30 36 30 30 30 30 30 30 37 42 37 45 0D 0A",
		tcpRsp:   "00010000000601060000007B",
	},
	{
		funccode: FuncCodeWriteHold,
		devid:    0x19,
		regaddr:  0x0099,
		reglen:   0x0001,
		arg:      []uint16{0x3156},
		sernum:   0x01,
		rtuReq:   "190600993156CF93",
		asciiReq: "3A 31 39 30 36 30 30 39 39 33 31 35 36 43 31 0D 0A",
		tcpReq:   "000100000006190600993156",
		rtuRsp:   "190600993156CF93",
		asciiRsp: "3A 31 39 30 36 30 30 39 39 33 31 35 36 43 31 0D 0A",
		tcpRsp:   "000100000006190600993156",
	},
	{
		funccode: FuncCodeWriteHold,
		devid:    0x72,
		regaddr:  0x0921,
		reglen:   0x0001,
		arg:      []uint16{0x0459},
		sernum:   0x01,
		rtuReq:   "72060921045913A5",
		asciiReq: "3A 37 32 30 36 30 39 32 31 30 34 35 39 30 31 0D 0A",
		tcpReq:   "000100000006720609210459",
		rtuRsp:   "72060921045913A5",
		asciiRsp: "3A 37 32 30 36 30 39 32 31 30 34 35 39 30 31 0D 0A",
		tcpRsp:   "000100000006720609210459",
	},
	{
		funccode: FuncCodeWriteHold,
		devid:    0x47,
		regaddr:  0x1591,
		reglen:   0x0001,
		arg:      []uint16{0x9459},
		sernum:   0x01,
		rtuReq:   "4706159194597C77",
		asciiReq: "3A 34 37 30 36 31 35 39 31 39 34 35 39 32 30 0D 0A",
		tcpReq:   "000100000006470615919459",
		rtuRsp:   "4706159194597C77",
		asciiRsp: "3A 34 37 30 36 31 35 39 31 39 34 35 39 32 30 0D 0A",
		tcpRsp:   "000100000006470615919459",
	},
	{
		funccode: FuncCodeWriteHold,
		devid:    0x81,
		regaddr:  0x1285,
		reglen:   0x0001,
		arg:      []uint16{0x7351},
		sernum:   0x01,
		rtuReq:   "8106128573516657",
		asciiReq: "3A 38 31 30 36 31 32 38 35 37 33 35 31 31 45 0D 0A",
		tcpReq:   "000100000006810612857351",
		rtuRsp:   "8106128573516657",
		asciiRsp: "3A 38 31 30 36 31 32 38 35 37 33 35 31 31 45 0D 0A",
		tcpRsp:   "000100000006810612857351",
	},
	{
		funccode: FuncCodeWriteHold,
		devid:    0x01,
		regaddr:  0x0043,
		reglen:   0x0001,
		arg:      []uint16{0x0003},
		sernum:   0x01,
		rtuReq:   "010600430003381F",
		asciiReq: "3A 30 31 30 36 30 30 34 33 30 30 30 33 42 33 0D 0A",
		tcpReq:   "000100000006010600430003",
		rtuRsp:   "010600430003381F",
		asciiRsp: "3A 30 31 30 36 30 30 34 33 30 30 30 33 42 33 0D 0A",
		tcpRsp:   "000100000006010600430003",
	},
	{
		funccode: FuncCodeWriteHolds,
		devid:    0x59,
		regaddr:  0x0043,
		reglen:   0x000D,
		arg:      []uint16{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22},
		sernum:   0x01,
		rtuReq:   "59100043000D1A000A000B000C000D000E000F0010001100120013001400150016C16A",
		asciiReq: "3A 35 39 31 30 30 30 34 33 30 30 30 44 31 41 30 30 30 41 30 30 30 42 30 30 30" +
			"43 30 30 30 44 30 30 30 45 30 30 30 46 30 30 31 30 30 30 31 31 30 30 31 32 30 30 31" +
			"33 30 30 31 34 30 30 31 35 30 30 31 36 35 44 0D 0A ",
		tcpReq:   "00010000002159100043000D1A000A000B000C000D000E000F0010001100120013001400150016",
		rtuRsp:   "59100043000DFD00",
		asciiRsp: "3A 35 39 31 30 30 30 34 33 30 30 30 44 34 37 0D 0A",
		tcpRsp:   "00010000000659100043000D",
	},
	{
		funccode: FuncCodeWriteHolds,
		devid:    0x01,
		regaddr:  0x0000,
		reglen:   0x0005,
		arg:      []uint16{1, 2, 3, 4, 5},
		sernum:   0x01,
		rtuReq:   "0110000000050A00010002000300040005EA6A",
		asciiReq: "3A 30 31 31 30 30 30 30 30 30 30 30 35 30 41 30 30 30 31 30 30 30 32 30 30 30" +
			"33 30 30 30 34 30 30 30 35 44 31 0D 0A",
		tcpReq:   "0001000000110110000000050A00010002000300040005",
		rtuRsp:   "011000000005000A",
		asciiRsp: "3A 30 31 31 30 30 30 30 30 30 30 30 35 45 41 0D 0A",
		tcpRsp:   "000100000006011000000005",
	},
	{
		funccode: FuncCodeWriteHolds,
		devid:    0x19,
		regaddr:  0x0000,
		reglen:   0x000A,
		arg:      []uint16{1, 2, 3, 4, 5, 1, 2, 3, 4, 5},
		sernum:   0x01,
		rtuReq:   "19100000000A1400010002000300040005000100020003000400055A41",
		asciiReq: "3A 31 39 31 30 30 30 30 30 30 30 30 41 31 34 30 30 30 31 30 30 30 32 30 30 30" +
			"33 30 30 30 34 30 30 30 35 30 30 30 31 30 30 30 32 30 30 30 33 30 30 30 34 30 30 30 35 39 42 0D 0A",
		tcpReq:   "00010000001B19100000000A140001000200030004000500010002000300040005",
		rtuRsp:   "19100000000A43D6",
		asciiRsp: "3A 31 39 31 30 30 30 30 30 30 30 30 41 43 44 0D 0A ",
		tcpRsp:   "00010000000619100000000A",
	},
	{
		funccode: FuncCodeReadInput,
		devid:    0x01,
		regaddr:  0x0000,
		reglen:   0x000A,
		arg:      []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		sernum:   0x01,
		rtuReq:   "01 04 00 00 00 0A 70 0D",
		asciiReq: "3A 30 31 30 34 30 30 30 30 30 30 30 41 46 31 0D 0A",
		tcpReq:   "00 01 00 00 00 06 01 04 00 00 00 0A ",
		rtuRsp:   "01 04 14 00 01 00 02 00 03 00 04 00 05 00 06 00 07 00 08 00 09 00 0A B9 F0 ",
		asciiRsp: "3A 30 31 30 34 31 34 30 30 30 31 30 30 30 32 30 30 30 33 30 30 30 34 30 30 " +
			"30 35 30 30 30 36 30 30 30 37 30 30 30 38 30 30 30 39 30 30 30 41 42 30 0D 0A",
		tcpRsp: "00 01 00 00 00 17 01 04 14 00 01 00 02 00 03 00 04 00 05 00 06 00 07 00 08 00 09 00 0A",
	},
	{
		funccode: FuncCodeReadInput,
		devid:    0x01,
		regaddr:  0x0000,
		reglen:   0x0006,
		arg:      []uint16{6, 1, 1, 1, 1, 7},
		sernum:   0x01,
		rtuReq:   "01 04 00 00 00 06 70 08 ",
		asciiReq: "3A 30 31 30 34 30 30 30 30 30 30 30 36 46 35 0D 0A",
		tcpReq:   "00 01 00 00 00 06 01 04 00 00 00 06",
		rtuRsp:   "01 04 0C 00 06 00 01 00 01 00 01 00 01 00 07 BB AD ",
		asciiRsp: "3A 30 31 30 34 30 43 30 30 30 36 30 30 30 31 30 30 30 31 30 30 30 31 30 30 " +
			"30 31 30 30 30 37 44 45 0D 0A ",
		tcpRsp: "00 01 00 00 00 0F 01 04 0C 00 06 00 01 00 01 00 01 00 01 00 07",
	},
}

func filterDevIDDefault(devid uint8, userdata any) bool {
	item := userdata.(testItem)
	return devid == item.devid
}

func checkDefault(regaddr, reglen uint16, isRead bool, userdata any) bool {
	_ = regaddr
	_ = reglen
	_ = isRead
	_ = userdata
	return true
}

// 将十六进制字符串转换为字节切片
func strToHex(s string) []byte {
	b, _ := hex.DecodeString(strings.ReplaceAll(strings.ToLower(s), " ", ""))
	return b
}

// 将字节切片转换为十六进制字符串
func strFromHex(b []byte) string {
	return hex.EncodeToString(b)
}

func testValuesCopy(arg *groArg, it *testItem) {
	switch it.funccode {
	case FuncCodeReadCoil,
		FuncCodeWriteCoil,
		FuncCodeWriteCoils,
		FuncCodeReadDiscrete:

		bits := make([]bool, len(it.arg))
		for i, v := range it.arg {
			bits[i] = v != 0
		}
		arg.SetBits(bits)
	case FuncCodeReadHold,
		FuncCodeWriteHold,
		FuncCodeWriteHolds,
		FuncCodeReadInput:
		arg.SetU16s(it.arg, binary.BigEndian)
	}
}

func testValuesCompare(t *testing.T, index int, arg *groArg, it *testItem) bool {
	switch it.funccode {
	case FuncCodeReadCoil,
		FuncCodeWriteCoil,
		FuncCodeWriteCoils,
		FuncCodeReadDiscrete:
		bits := arg.GetBits()

		if len(bits) < len(it.arg) {
			t.Fatalf("[%d] Arg mismatch: expected %v, got %v.\n", index, it.arg, bits)
			return false
		}
		for i, v := range it.arg {
			if bits[i] != (v != 0) {
				t.Fatalf("[%d] Arg mismatch: expected %v, got %v.\n", index, it.arg, bits)
				return false
			}
		}
		return true
	case FuncCodeReadHold,
		FuncCodeWriteHold,
		FuncCodeWriteHolds,
		FuncCodeReadInput:
		u16s := arg.GetU16s(binary.BigEndian)

		if len(u16s) < len(it.arg) {
			t.Fatalf("[%d] Arg mismatch: expected %v, got %v.\n", index, it.arg, u16s)
			return false
		}
		for i, v := range it.arg {
			if u16s[i] != v {
				t.Fatalf("[%d] Arg mismatch: expected %v, got %v.\n", index, it.arg, u16s)
				return false
			}
		}
		return true
	}
	return false
}

// TestPack 测试打包功能
func testPack(t *testing.T, protocol uint8) {
	m := New()
	m.Access.SetCheckCoil(checkDefault)
	m.Access.SetCheckDiscrete(checkDefault)
	m.Access.SetCheckHold(checkDefault)
	m.Access.SetCheckInput(checkDefault)
	m.Head.SetProtocol(protocol)

	var (
		reqBytes = make([]byte, 1024)
		rspBytes = make([]byte, 1024)
	)

	for i, item := range testItems {
		m.Access.SetUserData(item)
		m.Access.SetFilterDevID(filterDevIDDefault)
		m.Head.SetSerNum(item.sernum)
		m.Head.SetDevId(item.devid)
		m.Arg.SetFuncCode(item.funccode)
		m.Arg.SetRegAddr(item.regaddr)
		m.Arg.SetRegLen(item.reglen)
		m.Arg.SetU16s(item.arg, binary.BigEndian)

		testValuesCopy(&m.Arg, &item)

		// 测试打包请求
		{
			if err := m.PackRequest(reqBytes); err != nil {
				t.Fatalf("[%d] Failed to pack %s request: %v.\n", i, m.Head.GetProtocolString(), err.Error())
			}

			retlen := m.Result.GetRetLen()
			reqHex := strFromHex(reqBytes[:retlen])
			expect := item.packet(protocol, true)
			if !bytes.Equal(reqBytes[:retlen], expect) {
				t.Errorf("[%d] Expected = %x.\n", i, reqBytes[:retlen])
				t.Fatalf("[%d] %s request mismatch: Expected = %s, Actual = %s.\n", i, m.Head.GetProtocolString(), strFromHex(expect), reqHex)
			}
		}

		// 测试打包响应
		{

			if err := m.PackResponse(rspBytes); err != nil {
				t.Fatalf("[%d] Failed to pack %s response: %v.\n", i, m.Head.GetProtocolString(), m.Result.GetResult())
			}

			retlen := m.Result.GetRetLen()
			reqHex := strFromHex(rspBytes[:retlen])
			expect := item.packet(protocol, false)[0:retlen]
			if retlen == 0 || !bytes.Equal(rspBytes[:retlen], expect) {
				t.Fatalf("[%d] %s response mismatch: Expected = %s, Actual = %s.\n", i, m.Head.GetProtocolString(), strFromHex(expect), reqHex)
			}
		}
	}
}

// TestParse 测试解析功能
func testParse(t *testing.T, protocol uint8) {
	m := New()
	m.Access.SetCheckCoil(checkDefault)
	m.Access.SetCheckDiscrete(checkDefault)
	m.Access.SetCheckHold(checkDefault)
	m.Access.SetCheckInput(checkDefault)
	m.Head.SetProtocol(protocol)

	for i, item := range testItems {
		m.Access.SetUserData(item)
		m.Access.SetFilterDevID(filterDevIDDefault)
		m.Head.SetDevId(0)
		m.Result.Reset()
		m.Arg.Reset()

		// 测试解析请求
		{
			reqBytes := item.packet(protocol, true)

			if err := m.ParseRequest(reqBytes); err != nil {
				t.Fatalf("[%d] Failed to parse %s request: %s.\n", i, m.Head.GetProtocolString(), err.Error())
			}
			if m.Result.GetExcepCode() != ExcepNormal {
				t.Fatalf("[%d] ExcepCode mismatch: expected %d, got %d(%s).\n", i, ExcepNormal, m.Result.GetExcepCode(), m.Result.GetExcepCodeString())
			}
			if m.Head.GetDevId() != item.devid {
				t.Fatalf("[%d] DevId mismatch: expected %d, got %d.\n", i, item.devid, m.Head.GetDevId())
			}
			if m.Arg.GetFuncCode() != item.funccode {
				t.Fatalf("[%d] Funccode mismatch: expected %d, got %s.\n", i, item.funccode, m.Arg.GetFuncCodeString())
			}
			if m.Arg.GetRegAddr() != item.regaddr {
				t.Fatalf("[%d] RegAddr mismatch: expected %d, got %d.\n", i, item.regaddr, m.Arg.GetRegAddr())
			}
			if m.Arg.GetRegLen() != item.reglen {
				t.Fatalf("[%d] RegLen mismatch: expected %d, got %d.\n", i, item.reglen, m.Arg.GetRegLen())
			}
		}

		// 测试解析响应
		{
			rspBytes := item.packet(protocol, false)

			if err := m.ParseResponse(rspBytes); err != nil {
				t.Fatalf("[%d] Failed to parse %s response: %v.\n", i, m.Head.GetProtocolString(), err.Error())
			}
			if m.Result.GetExcepCode() != ExcepNormal {
				t.Fatalf("[%d] ExcepCode mismatch: expected %d, got %d(%s).\n", i, ExcepNormal, m.Result.GetExcepCode(), m.Result.GetExcepCodeString())
			}

			if !testValuesCompare(t, i, &m.Arg, &item) {
				t.Fatalf("[%d] Arg mismatch.\n", i)
			}
		}
	}
}
