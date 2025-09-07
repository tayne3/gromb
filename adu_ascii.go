// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gromb

//	<--------------------------------- MODBUS Ascii ADU -------------------------------->
//	                            <--------- MODBUS PDU -------->
//	+-------------+-------------+---------------+-------------+------------+------------+
//	| Start       | Address     | Function Code | Data        | LRC        | End        |
//	| 1 Byte      | 1 Byte * 2  | 1 Byte * 2    | n Bytes * 2 | 1 Byte * 2 | 2 Byte     |
//	+-------------+-------------+---------------+-------------+------------+------------+

const (
	MaxAsciiLen = 256
	MinAsciiLen = 4
	StartChar   = 0x3a
	EndHigh     = 0x0d
	EndLow      = 0x0a
)

// 封装报文 (Modbus Ascii)
func (m *Modbus) asciiPack(isReq bool) {
	if m.Box.GetMax() < MinAsciiLen {
		m.Result.SetResult(ErrResultBufTooShort)
		return
	}

	var hex []uint8
	var box groBox
	box.Init(&hex, MaxAsciiLen)
	box.PutU8(m.Head.GetDevId())
	box.AddLast(1)

	ret := Pack(&m.Result, &m.Arg, &box, isReq)
	if ret < 0 {
		return
	}
	length := uint16(ret)
	box.PutU8(LRCCalcul(box.GetBuffer(0, length+1)))

	m.Box.PutU8(StartChar)
	m.Box.PutU8s(asciiFromHex(box.GetBuffer(0, length+2)))
	m.Box.PutU8(EndHigh)
	m.Box.PutU8(EndLow)

	m.Result.SetResult(nil)
	m.Result.SetRetLen(length*2 + 7)
}

// 解析报文 (Modbus Ascii)
func (m *Modbus) asciiParse(isReq bool) {
	if m.Box.Size() < MinAsciiLen {
		m.Result.SetResult(ErrResultTooShort)
		return
	}

	if m.Box.GetU8(0) != StartChar {
		m.Result.SetResult(ErrResultAsciiStart)
		return
	}

	hex := asciiToHex(m.Box.GetBuffer(1, m.Box.Size()-2))

	var box groBox
	box.Init(&hex, 1024)

	if isReq {
		if m.Access.FilterDevID != nil && !m.Access.FilterDevID(box.GetU8(0), m.Access.UserData) {
			m.Result.SetResult(ErrResultDevID)
			return
		}
	} else {
		if m.Head.GetDevId() != box.GetU8(0) {
			m.Result.SetResult(ErrResultDevID)
			return
		}
	}
	m.Head.SetDevId(box.GetU8(0))

	box.AddLast(1)
	ret := Parse(&m.Result, &m.Access, &m.Arg, &box, isReq)
	if ret < 0 {
		return
	}
	length := uint16(ret)
	{
		box.SubLast(1)
		lrc := LRCCalcul(box.GetBuffer(0, length+1))
		if lrc != box.GetU8(length+1) {
			m.Result.SetResult(ErrResultAsciiLrc)
			return
		}
	}

	if m.Box.GetU8(length*2+5) != EndHigh || m.Box.GetU8(length*2+6) != EndLow {
		m.Result.SetResult(ErrResultAsciiEnd)
		return
	}

	m.Result.SetResult(nil)
	m.Result.SetRetLen(length*2 + 7)
}