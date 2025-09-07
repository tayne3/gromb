package gromb

//  <------- MODBUS Exception PDU ------->
//	+-----------------+------------------+
//	| Function Code   | Exception Code   |
//	| 1 Byte          | 1 Byte           |
//	+-----------------+------------------+

// 封装响应报文-异常
func packResponseExcep(result *groResult, arg *groArg, box *groBox) int {
	box.PutU8(arg.GetFuncCode() | 0x80)
	box.PutU8(result.GetExcepCode())
	return 2
}

// 解析响应报文-异常
func parseResponseExcep(result *groResult, box *groBox) int {
	if box.ThisSize() < 2 {
		result.SetResult(ErrResultTooShort)
		return -1
	}

	result.SetExcepCode(box.GetU8(1))
	return 2
}
