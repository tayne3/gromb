package gromb

import (
	"encoding/binary"
)

// <--------------- MODBUS Read Coils Request PDU ------------------->
// +-------------------+---------------------+-----------------------+
// | Function Code     | Starting Address    | Quantity of Coils     |
// | 1 Byte            | 2 Bytes             | 2 Bytes               |
// +-------------------+---------------------+-----------------------+

// 封装请求报文-读取线圈
func packRequestReadCoil(result *groResult, arg *groArg, box *groBox) int {
	regaddr := arg.GetRegAddr()
	reglen := arg.GetRegLen()

	// 检查参数
	if reglen < 0x0001 || reglen > 0x07D0 {
		result.SetResult(ErrResultRegLen)
		return -1
	}

	// 填充报文
	box.PutU8(FuncCodeReadCoil)
	box.PutU16(regaddr, binary.BigEndian)
	box.PutU16(reglen, binary.BigEndian)
	return 5
}

// 解析请求报文-读取线圈
func parseRequestReadCoil(result *groResult, access *groAccess, arg *groArg, box *groBox) int {
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
	} else if access.CheckCoil == nil {
		result.SetExcepCode(ExcepIllFuncCode)
	} else if !access.CheckCoil(regaddr, reglen, true, access.UserData) {
		result.SetExcepCode(ExcepIllDataAddr)
	} else {
		arg.SetRegAddr(regaddr)
		arg.SetRegLen(reglen)
	}
	return 5
}

// <--------------- MODBUS Read Coils Response PDU -------------->
// +-------------------+--------------------+--------------------+
// | Function Code     | Byte Count         | Coil Status        |
// | 1 Byte            | 1 Byte             | n Bytes            |
// +-------------------+--------------------+--------------------+

// 封装响应报文-读取线圈
func packResponseReadCoil(result *groResult, arg *groArg, box *groBox) int {
	reglen := arg.GetRegLen()  // 线圈数量
	number := (reglen + 7) / 8 // 字节数

	// 检查参数
	if reglen < 0x0001 || reglen > 0x07D0 {
		result.SetResult(ErrResultRegLen)
		return -1
	}

	// 填充报文
	box.PutU8(FuncCodeReadCoil)
	box.PutU8(uint8(number))
	box.PutU8s(arg.GetU8s()[0:number])

	return 2 + int(number)
}

// 解析响应报文-读取线圈
func parseResponseReadCoil(result *groResult, arg *groArg, box *groBox) int {
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

// <----------- MODBUS Write Single Coil Request PDU ------------>
// +-------------------+--------------------+--------------------+
// | Function Code     | Output Address     | Output Value       |
// | 1 Byte            | 2 Bytes            | 2 Bytes            |
// +-------------------+--------------------+--------------------+

// 封装请求报文-写入线圈
func packRequestWriteCoil(result *groResult, arg *groArg, box *groBox) int {
	_ = result
	box.PutU8(FuncCodeWriteCoil)
	box.PutU16(arg.GetRegAddr(), binary.BigEndian)

	if arg.GetU8(0)&0x01 != 0 {
		box.PutU16(0xFF00, binary.BigEndian)
	} else {
		box.PutU16(0x0000, binary.BigEndian)
	}

	return 5
}

// 解析请求报文-写入线圈
func parseRequestWriteCoil(result *groResult, access *groAccess, arg *groArg, box *groBox) int {
	// 检查报文是否过短
	if box.ThisSize() < 5 {
		result.SetResult(ErrResultTooShort)
		return -1
	}

	regaddr := box.GetU16(1, binary.BigEndian) // 寄存器地址
	value := box.GetU16(3, binary.BigEndian)   // 线圈状态值

	// 检查参数
	if access.CheckCoil == nil {
		result.SetExcepCode(ExcepIllFuncCode)
	} else if !access.CheckCoil(regaddr, 1, false, access.UserData) {
		result.SetExcepCode(ExcepIllDataAddr)
	} else if value == 0x0000 {
		arg.SetRegAddr(regaddr)
		arg.SetRegLen(1)
		arg.SetU8s([]uint8{0x00})
	} else if value == 0xFF00 {
		arg.SetRegAddr(regaddr)
		arg.SetRegLen(1)
		arg.SetU8s([]uint8{0x01})
	} else {
		result.SetExcepCode(ExcepIllDataValue)
	}

	return 5
}

// <----------- MODBUS Write Single Coil Response PDU ----------->
// +-------------------+--------------------+--------------------+
// | Function Code     | Output Address     | Output Value       |
// | 1 Byte            | 2 Bytes            | 2 Bytes            |
// +-------------------+--------------------+--------------------+

// 封装响应报文-写入单个线圈
func packResponseWriteCoil(result *groResult, arg *groArg, box *groBox) int {
	_ = result
	box.PutU8(FuncCodeWriteCoil)
	box.PutU16(arg.GetRegAddr(), binary.BigEndian)

	if arg.GetU8(0)&0x01 != 0 {
		box.PutU16(0xFF00, binary.BigEndian)
	} else {
		box.PutU16(0x0000, binary.BigEndian)
	}

	return 5
}

// 解析响应报文-写入单个线圈
func parseResponseWriteCoil(result *groResult, arg *groArg, box *groBox) int {
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

	// 检查线圈状态值
	{
		v1 := box.GetU16(3, binary.BigEndian) == 0xFF00
		v2 := arg.GetU8(0)&0x01 != 0
		if v1 != v2 {
			result.SetResult(ErrResultRegValue)
			return -1
		}
	}

	return 5
}

// <--------------------------------- MODBUS Write Multiple Coils Request PDU ----------------------------->
// +-------------------+--------------------+--------------------+--------------------+--------------------+
// | Function Code     | Starting Address   | Quantity of Coils  | Byte Count         | Coil Status        |
// | 1 Byte            | 2 Bytes            | 2 Bytes            | 1 Byte             | n Bytes            |
// +-------------------+--------------------+--------------------+--------------------+--------------------+

// 封装请求报文-写入多个线圈
func packRequestWriteCoils(result *groResult, arg *groArg, box *groBox) int {
	regaddr := arg.GetRegAddr() // 寄存器地址
	reglen := arg.GetRegLen()
	number := (reglen + 7) / 8

	if reglen < 0x0001 || reglen > 0x07B0 {
		result.SetResult(ErrResultRegLen)
		return -1
	}

	box.PutU8(FuncCodeWriteCoils)
	box.PutU16(regaddr, binary.BigEndian)
	box.PutU16(reglen, binary.BigEndian)
	box.PutU8(uint8(number))
	box.PutU8s(arg.GetU8s()[0:number])

	return 6 + int(number)
}

// 解析请求报文-写入多个线圈
func parseRequestWriteCoils(result *groResult, access *groAccess, arg *groArg, box *groBox) int {
	// 检查报文是否过短
	if box.ThisSize() < 6 {
		result.SetResult(ErrResultTooShort)
		return -1
	}

	regaddr := box.GetU16(1, binary.BigEndian)
	reglen := box.GetU16(3, binary.BigEndian)
	number := (reglen + 7) / 8

	if box.ThisSize() < 6+number {
		result.SetResult(ErrResultLength)
		return -1
	}

	if reglen < 0x0001 || reglen > 0x07B0 || number != uint16(box.GetU8(5)) {
		result.SetExcepCode(ExcepIllDataValue)
	} else if access.CheckCoil == nil {
		result.SetExcepCode(ExcepIllFuncCode)
	} else if !access.CheckCoil(regaddr, reglen, false, access.UserData) {
		result.SetExcepCode(ExcepIllDataAddr)
	} else {
		arg.SetRegAddr(regaddr)
		arg.SetRegLen(reglen)
		arg.SetU8s(box.GetThisBuffer(6, 6+number))
	}

	return 6 + int(number)
}

// <--------- MODBUS Write Multiple Coils Response PDU ---------->
// +-------------------+--------------------+--------------------+
// | Function Code     | Starting Address   | Quantity of Coils  |
// | 1 Byte            | 2 Bytes            | 2 Bytes            |
// +-------------------+--------------------+--------------------+

// 封装响应报文-写入多个线圈
func packResponseWriteCoils(result *groResult, arg *groArg, box *groBox) int {
	_ = result
	regaddr := arg.GetRegAddr()
	reglen := arg.GetRegLen()

	box.PutU8(FuncCodeWriteCoils)
	box.PutU16(regaddr, binary.BigEndian)
	box.PutU16(reglen, binary.BigEndian)
	return 5
}

// 解析响应报文-写入多个线圈
func parseResponseWriteCoils(result *groResult, arg *groArg, box *groBox) int {
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
