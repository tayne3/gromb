package gromb

import (
	"encoding/binary"
)

// <---------- MODBUS Read Input Registers Request PDU -------------->
// +-------------------+---------------------+-----------------------+
// | Function Code     | Starting Address    | Quantity of Registers |
// | 1 Byte            | 2 Bytes             | 2 Bytes               |
// +-------------------+---------------------+-----------------------+

// 封装请求报文-读取输入寄存器
func packRequestReadInput(result *groResult, arg *groArg, box *groBox) int {
	regaddr := arg.GetRegAddr()
	reglen := arg.GetRegLen()

	// 检查参数
	if reglen < 0x0001 || reglen > 0x007D {
		result.SetResult(ErrResultRegLen)
		return -1
	}

	// 填充报文
	box.PutU8(FuncCodeReadInput)
	box.PutU16(regaddr, binary.BigEndian)
	box.PutU16(reglen, binary.BigEndian)
	return 5
}

// 解析请求报文-读取输入寄存器
func parseRequestReadInput(result *groResult, access *groAccess, arg *groArg, box *groBox) int {
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
	} else if access.CheckInput == nil {
		result.SetExcepCode(ExcepIllFuncCode)
	} else if !access.CheckInput(regaddr, reglen, true, access.UserData) {
		result.SetExcepCode(ExcepIllDataAddr)
	} else {
		arg.SetRegAddr(regaddr)
		arg.SetRegLen(reglen)
	}
	return 5
}

// <---------- MODBUS Read Input Registers Response PDU --------->
// +-------------------+--------------------+--------------------+
// | Function Code     | Byte Count         | Register Arg    |
// | 1 Byte            | 1 Byte             | n Bytes            |
// +-------------------+--------------------+--------------------+

// 封装响应报文-读取输入寄存器
func packResponseReadInput(result *groResult, arg *groArg, box *groBox) int {
	reglen := arg.GetRegLen() // 寄存器数量
	number := reglen * 2      // 字节数

	// 检查参数
	if reglen < 0x0001 || reglen > 0x007D {
		result.SetResult(ErrResultRegLen)
		return -1
	}

	// 填充报文
	box.PutU8(FuncCodeReadInput)
	box.PutU8(uint8(number))
	box.PutU8s(arg.GetU8s()[0:number])

	return 2 + int(number)
}

// 解析响应报文-读取输入寄存器
func parseResponseReadInput(result *groResult, arg *groArg, box *groBox) int {
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
