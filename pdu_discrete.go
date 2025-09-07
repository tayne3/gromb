// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gromb

import (
	"encoding/binary"
)

// <---------- MODBUS Read Discrete Inputs Request PDU -------------->
// +-------------------+---------------------+-----------------------+
// | Function Code     | Starting Address    | Quantity of Discretes |
// | 1 Byte            | 2 Bytes             | 2 Bytes               |
// +-------------------+---------------------+-----------------------+

// 封装请求报文-读取离散量输入
func packRequestReadDiscretes(result *groResult, arg *groArg, box *groBox) int {
	regaddr := arg.GetRegAddr()
	reglen := arg.GetRegLen()

	// 检查参数
	if reglen < 0x0001 || reglen > 0x07D0 {
		result.SetResult(ErrResultRegLen)
		return -1
	}

	// 填充报文
	box.PutU8(FuncCodeReadDiscrete)
	box.PutU16(regaddr, binary.BigEndian)
	box.PutU16(reglen, binary.BigEndian)
	return 5
}

// 解析请求报文-读取离散量输入
func parseRequestReadDiscretes(result *groResult, access *groAccess, arg *groArg, box *groBox) int {
	// 检查报文是否过短
	if box.ThisSize() < 5 {
		result.SetResult(ErrResultTooShort)
		return -1
	}

	regaddr := box.GetU16(1, binary.BigEndian)
	reglen := box.GetU16(3, binary.BigEndian)

	// 检查参数
	if reglen < 0x0001 || reglen > 0x07D0 {
		result.SetExcepCode(ExcepIllDataValue)
	} else if access.CheckDiscrete == nil {
		result.SetExcepCode(ExcepIllFuncCode)
	} else if !access.CheckDiscrete(regaddr, reglen, true, access.UserData) {
		result.SetExcepCode(ExcepIllDataAddr)
	} else {
		arg.SetRegAddr(regaddr)
		arg.SetRegLen(reglen)
	}
	return 5
}

// <---------- MODBUS Read Discrete Inputs Response PDU --------->
// +-------------------+--------------------+--------------------+
// | Function Code     | Byte Count         | Discrete Status    |
// | 1 Byte            | 1 Byte             | n Bytes            |
// +-------------------+--------------------+--------------------+

// 封装响应报文-读取离散量输入
func packResponseReadDiscretes(result *groResult, arg *groArg, box *groBox) int {
	reglen := arg.GetRegLen()  // 离散量数量
	number := (reglen + 7) / 8 // 字节数

	// 检查参数
	if reglen < 0x0001 || reglen > 0x07D0 {
		result.SetResult(ErrResultRegLen)
		return -1
	}

	// 填充报文
	box.PutU8(FuncCodeReadDiscrete)
	box.PutU8(uint8(number))
	box.PutU8s(arg.GetU8s()[0:number])

	return 2 + int(number)
}

// 解析响应报文-读取离散量输入
func parseResponseReadDiscretes(result *groResult, arg *groArg, box *groBox) int {
	// 检查报文是否过短
	if box.ThisSize() < 2 {
		result.SetResult(ErrResultTooShort)
		return -1
	}

	// 检查报文是否过短
	number := uint16(box.GetU8(1))
	if box.ThisSize() < 2+number {
		result.SetResult(ErrResultTooShort)
		return -1
	}

	// 检查字节数/寄存器数量是否错误
	reglen := arg.GetRegLen()
	if number != (reglen+7)/8 {
		result.SetResult(ErrResultLength)
		return -1
	}

	arg.SetU8s(box.GetThisBuffer(2, 2+number))

	return 2 + int(number)
}
