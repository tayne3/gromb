package gromb

import (
	"encoding/binary"
)

//	<------------------------------- MODBUS TCP ADU ------------------------------->
//                                                   <--------- MODBUS PDU -------->
//	+-----------+-----------+------------+-----------+---------------+-------------+
//	| TID       | PID       | Length     | UID       | Function Code | Data        |
//	| 2 Byte    | 2 Byte    | 2 Byte     | 1 Byte    | 1 Byte        | n Byte      |
//	+-----------+-----------+------------+-----------+---------------+-------------+

const (
	MaxTCPLen = 256
	MinTCPLen = 4
)

// 封装报文 (Modbus TCP)
func (m *Modbus) tcpPack(isReq bool) {
	if m.Box.GetMax() < MinTCPLen {
		m.Result.SetResult(ErrResultBufTooShort)
		return
	}

	m.Box.SetMax(MaxTCPLen)
	m.Box.PutU16(m.Head.GetSerNum(), binary.BigEndian)
	m.Box.PutU16(0x0000, binary.BigEndian)
	m.Box.PutU16(0x0000, binary.BigEndian)
	m.Box.PutU8(m.Head.GetDevId())
	m.Box.AddLast(7)

	ret := Pack(&m.Result, &m.Arg, &m.Box, isReq)
	if ret < 0 {
		return
	}
	length := uint16(ret)
	m.Box.SubLast(7)
	m.Box.SetU16(4, length+1, binary.BigEndian)

	m.Result.SetResult(nil)
	m.Result.SetRetLen(length + 7)
}

// 解析报文 (Modbus TCP)
func (m *Modbus) tcpParse(isReq bool) {
	if m.Box.ThisSize() < 7 {
		m.Result.SetResult(ErrResultTooShort)
		return
	}

	sernum := m.Box.GetU16(0, binary.BigEndian)
	if isReq {
		m.Head.SetSerNum(sernum)
	} else if sernum != m.Head.GetSerNum() {
		m.Result.SetResult(ErrResultTcpSerNum)
		return
	}

	if m.Box.GetU16(2, binary.BigEndian) != 0x0000 {
		m.Result.SetResult(ErrResultTcpProtocol)
		return
	}

	number := m.Box.GetU16(4, binary.BigEndian)
	if m.Box.ThisSize() < 6+number {
		m.Result.SetResult(ErrResultTooShort)
		return
	}

	if isReq {
		if m.Access.FilterDevID != nil && !m.Access.FilterDevID(m.Box.GetU8(6), m.Access.UserData) {
			m.Result.SetResult(ErrResultDevID)
			return
		}
	} else {
		if m.Head.GetDevId() != m.Box.GetU8(6) {
			m.Result.SetResult(ErrResultDevID)
			return
		}
	}
	m.Head.SetDevId(m.Box.GetU8(6))

	m.Box.AddLast(7)
	ret := Parse(&m.Result, &m.Access, &m.Arg, &m.Box, isReq)
	if ret < 0 {
		return
	}
	length := uint16(ret)

	m.Result.SetResult(nil)
	m.Result.SetRetLen(length + 7)
}
