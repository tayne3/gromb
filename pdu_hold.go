// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gromb

import (
	"encoding/binary"
)

// <--------- MODBUS Read Holding Registers Request PDU ------------->
// +-------------------+---------------------+-----------------------+
// | Function Code     | Starting Address    | Quantity of Registers |
// | 1 Byte            | 2 Bytes             | 2 Bytes               |
// +-------------------+---------------------+-----------------------+

// 封装请求报文-读取保持寄存器
func packRequestReadHold(result *groResult, arg *groArg, box *groBox) int {
	regaddr := arg.GetRegAddr()
	reglen := arg.GetRegLen()

	// 检查参数
	if reglen < 0x0001 || reglen > 0x007D {
		result.SetResult(ErrResultRegLen)
		return -1
	}

	// 填充报文
	box.PutU8(FuncCodeReadHold)
	box.PutU16(regaddr, binary.BigEndian)
	box.PutU16(reglen, binary.BigEndian)
	return 5
}

// 解析请求报文-读取保持寄存器
func parseRequestReadHold(result *groResult, access *groAccess, arg *groArg, box *groBox) int {
	// 检查报文是否过短
	if box.ThisSize() < 5 {
		result.SetResult(ErrResultTooShort)
		return -1
	}

	regaddr := box.GetU16(1, binary.BigEndian)
	reglen := box.GetU16(3, binary.BigEndian)

	// 检查参数
	if reglen < 0x0001 || reglen > 0x007D {
		result.SetExcepCode(ExcepIllDataValue)
	} else if access.CheckHold == nil {
		result.SetExcepCode(ExcepIllFuncCode)
	} else if !access.CheckHold(regaddr, reglen, true, access.UserData) {
		result.SetExcepCode(ExcepIllDataAddr)
	} else {
		arg.SetRegAddr(regaddr)
		arg.SetRegLen(reglen)
	}
	return 5
}

// <--------- MODBUS Read Holding Registers Response PDU -------->
// +-------------------+--------------------+--------------------+
// | Function Code     | Byte Count         | Register Arg    |
// | 1 Byte            | 1 Byte             | n Bytes            |
// +-------------------+--------------------+--------------------+

// 封装响应报文-读取保持寄存器
func packResponseReadHold(result *groResult, arg *groArg, box *groBox) int {
	reglen := arg.GetRegLen() // 寄存器数量
	number := reglen * 2      // 字节数

	// 检查参数
	if reglen < 0x0001 || reglen > 0x007D {
		result.SetResult(ErrResultRegLen)
		return -1
	}

	// 填充报文
	box.PutU8(FuncCodeReadHold)
	box.PutU8(uint8(number))
	box.PutU8s(arg.GetU8s()[0:number])

	return 2 + int(number)
}

// 解析响应报文-读取保持寄存器
func parseResponseReadHold(result *groResult, arg *groArg, box *groBox) int {
	// 检查报文是否过短
	if box.ThisSize() < 2 {
		result.SetResult(ErrResultTooShort)
		return -1
	}

	// 检查字节数/寄存器数量是否错误
	number := uint16(box.GetU8(1))
	if box.ThisSize() < 2+number {
		result.SetResult(ErrResultTooShort)
		return -1
	}
	if number != arg.GetRegLen()*2 {
		result.SetResult(ErrResultLength)
		return -1
	}

	arg.SetU8s(box.GetThisBuffer(2, 2+number))

	return 2 + int(number)
}

// <----- MODBUS Write Single Holding Register Request PDU ------>
// +-------------------+--------------------+--------------------+
// | Function Code     | Register Address   | Register Value     |
// | 1 Byte            | 2 Bytes            | 2 Bytes            |
// +-------------------+--------------------+--------------------+

// 封装请求报文-写入单个保持寄存器
func packRequestWriteHold(result *groResult, arg *groArg, box *groBox) int {
	_ = result
	regaddr := arg.GetRegAddr()

	// 填充报文
	box.PutU8(FuncCodeWriteHold)
	box.PutU16(regaddr, binary.BigEndian)
	box.PutU8s(arg.GetU8s()[0:2])
	return 5
}

// 解析请求报文-写入单个保持寄存器
func parseRequestWriteHold(result *groResult, access *groAccess, arg *groArg, box *groBox) int {
	// 检查报文是否过短
	if box.ThisSize() < 5 {
		result.SetResult(ErrResultTooShort)
		return -1
	}

	regaddr := box.GetU16(1, binary.BigEndian) // 寄存器地址

	// 检查参数
	if access.CheckHold == nil {
		result.SetExcepCode(ExcepIllFuncCode)
	} else if !access.CheckHold(regaddr, 1, false, access.UserData) {
		result.SetExcepCode(ExcepIllDataAddr)
	} else {
		arg.SetRegAddr(regaddr)
		arg.SetRegLen(1)
		arg.SetU8s(box.GetU8s(3, 5))
	}

	return 5
}

// <----- MODBUS Write Single Holding Register Response PDU ----->
// +-------------------+--------------------+--------------------+
// | Function Code     | Register Address   | Register Value     |
// | 1 Byte            | 2 Bytes            | 2 Bytes            |
// +-------------------+--------------------+--------------------+

// 封装响应报文-写入单个保持寄存器
func packResponseWriteHold(result *groResult, arg *groArg, box *groBox) int {
	_ = result
	regaddr := arg.GetRegAddr() // 寄存器地址

	// 填充报文
	box.PutU8(FuncCodeWriteHold)
	box.PutU16(regaddr, binary.BigEndian)
	box.PutU8s(arg.GetU8s()[0:2])

	return 5
}

// 解析响应报文-写入单个保持寄存器
func parseResponseWriteHold(result *groResult, arg *groArg, box *groBox) int {
	// 检查报文是否过短
	if box.ThisSize() < 5 {
		result.SetResult(ErrResultTooShort)
		return -1
	}

	// 检查寄存器地址
	regaddr := box.GetU16(1, binary.BigEndian)
	if arg.GetRegAddr() != regaddr {
		result.SetResult(ErrResultRegAddr)
		return -1
	}

	// 检查寄存器值
	value := box.GetU16(3, binary.BigEndian)
	if arg.GetU16(0, binary.BigEndian) != value {
		result.SetResult(ErrResultRegValue)
		return -1
	}

	return 5
}

// <--------------------------- MODBUS Write Multiple Holding Registers Request PDU ----------------------->
// +-------------------+--------------------+--------------------+--------------------+--------------------+
// | Function Code     | Starting Address   | Quantity of Regis  | Byte Count         | Register Arg    |
// | 1 Byte            | 2 Bytes            | 2 Bytes            | 1 Byte             | n Bytes            |
// +-------------------+--------------------+--------------------+--------------------+--------------------+

// 封装请求报文-写入多个保持寄存器
func packRequestWriteHolds(result *groResult, arg *groArg, box *groBox) int {
	regaddr := arg.GetRegAddr() // 寄存器地址
	reglen := arg.GetRegLen()   // 寄存器数量
	number := reglen * 2        // 字节数

	if reglen < 0x0001 || reglen > 0x007B {
		result.SetResult(ErrResultRegLen)
		return -1
	}

	box.PutU8(FuncCodeWriteHolds)
	box.PutU16(regaddr, binary.BigEndian)
	box.PutU16(reglen, binary.BigEndian)
	box.PutU8(uint8(number))
	box.PutU8s(arg.GetU8s()[0:number])

	return 6 + int(number)
}

// 解析请求报文-写入多个保持寄存器
func parseRequestWriteHolds(result *groResult, access *groAccess, arg *groArg, box *groBox) int {
	// 检查报文是否过短
	if box.ThisSize() < 6 {
		result.SetResult(ErrResultTooShort)
		return -1
	}

	regaddr := box.GetU16(1, binary.BigEndian)
	reglen := box.GetU16(3, binary.BigEndian)
	number := uint16(box.GetU8(5))

	// 校验报文长度
	if box.ThisSize() < 6+uint16(number) {
		result.SetResult(ErrResultRegLen)
		return -1
	}

	if reglen < 0x0001 || reglen > 0x007B || number != reglen*2 {
		result.SetExcepCode(ExcepIllDataValue)
	} else if access.CheckHold == nil {
		result.SetExcepCode(ExcepIllFuncCode)
	} else if !access.CheckHold(regaddr, reglen, false, access.UserData) {
		result.SetExcepCode(ExcepIllDataAddr)
	} else {
		arg.SetRegAddr(regaddr)
		arg.SetRegLen(reglen)
		arg.SetU8s(box.GetThisBuffer(6, 6+number))
	}

	return 6 + int(number)
}

// <--- MODBUS Write Multiple Holding Registers Response PDU ---->
// +-------------------+--------------------+--------------------+
// | Function Code     | Starting Address   | Quantity of Regis  |
// | 1 Byte            | 2 Bytes            | 2 Bytes            |
// +-------------------+--------------------+--------------------+

// 封装响应报文-写入多个保持寄存器
func packResponseWriteHolds(result *groResult, arg *groArg, box *groBox) int {
	regaddr := arg.GetRegAddr()
	reglen := arg.GetRegLen()

	// 检查参数
	if reglen < 0x0001 || reglen > 0x0078 {
		result.SetResult(ErrResultRegLen)
		return -1
	}

	box.PutU8(FuncCodeWriteHolds)
	box.PutU16(regaddr, binary.BigEndian)
	box.PutU16(reglen, binary.BigEndian)
	return 5
}

// 解析响应报文-写入多个保持寄存器
func parseResponseWriteHolds(result *groResult, arg *groArg, box *groBox) int {
	// 检查报文是否过短
	if box.ThisSize() < 5 {
		result.SetResult(ErrResultTooShort)
		return -1
	}

	// 检查寄存器地址
	regaddr := box.GetU16(1, binary.BigEndian)
	if arg.GetRegAddr() != regaddr {
		result.SetResult(ErrResultRegAddr)
		return -1
	}

	// 检查寄存器数量
	reglen := box.GetU16(3, binary.BigEndian)
	if arg.GetRegLen() != reglen {
		result.SetResult(ErrResultRegLen)
		return -1
	}

	return 5
}
