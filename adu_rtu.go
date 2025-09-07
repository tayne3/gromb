package gromb

import (
	"encoding/binary"
)

//	<-------------------- MODBUS RTU ADU ------------------->
//                <--------- MODBUS PDU -------->
//	+-------------+---------------+-------------+-----------+
//	| Address     | Function Code | Data        | CRC       |
//	| 1 Byte      | 1 Byte        | n Bytes     | 2 Byte    |
//	+-------------+---------------+-------------+-----------+

const (
	MaxRTULen = 256
	MinRTULen = 4
)

// 封装报文 (Modbus RTU)
func (m *Modbus) rtuPack(isReq bool) {
	if m.Box.GetMax() < MinRTULen {
		m.Result.SetResult(ErrResultBufTooShort)
		return
	}
	m.Box.SetMax(MaxRTULen)
	m.Box.PutU8(m.Head.GetDevId())
	m.Box.AddLast(1)

	ret := Pack(&m.Result, &m.Arg, &m.Box, isReq)
	if ret < 0 {
		return
	}
	length := uint16(ret)

	{
		crc := CRC16(m.Box.GetBuffer(0, length+1))
		m.Box.PutU16(crc, binary.LittleEndian)
	}

	m.Result.SetResult(nil)
	m.Result.SetRetLen(length + 3)
}

// 解析报文 (Modbus RTU)
func (m *Modbus) rtuParse(isReq bool) {
	if m.Box.Size() < MinRTULen {
		m.Result.SetResult(ErrResultTooShort)
		return
	}

	if isReq {
		if m.Access.FilterDevID != nil && !m.Access.FilterDevID(m.Box.GetU8(0), m.Access.UserData) {
			m.Result.SetResult(ErrResultDevID)
			return
		}
	} else {
		if m.Head.GetDevId() != m.Box.GetU8(0) {
			m.Result.SetResult(ErrResultDevID)
			return
		}
	}
	m.Head.SetDevId(m.Box.GetU8(0))

	m.Box.AddLast(1)
	ret := Parse(&m.Result, &m.Access, &m.Arg, &m.Box, isReq)
	if ret < 0 {
		return
	}
	len := uint16(ret)

	{
		crc1 := CRC16(m.Box.GetBuffer(0, len+1))
		crc2 := m.Box.GetU16(len, binary.LittleEndian)
		if crc1 != crc2 {
			m.Result.SetResult(ErrResultRtuCrc)
			return
		}
	}

	m.Result.SetResult(nil)
	m.Result.SetRetLen(len + 3)
}
