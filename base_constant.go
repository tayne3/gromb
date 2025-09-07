package gromb

import "errors"

// Modbus 协议类型 (Modbus Protocol Type)
const (
	ProtocolRTU = iota
	ProtocolAscii
	ProtocolTCP
)

func ProtocolToString(p uint8) string {
	switch p {
	case ProtocolRTU:
		return "modbus RTU"
	case ProtocolAscii:
		return "modbus Ascii"
	case ProtocolTCP:
		return "modbus TCP"
	}
	return "unknown protocol"
}

// Modbus 功能码 (Modbus Function Code)
const (
	FuncCodeReadCoil     = 0x01 // 读线圈
	FuncCodeReadDiscrete = 0x02 // 读离散量输入
	FuncCodeReadHold     = 0x03 // 读保持寄存器
	FuncCodeReadInput    = 0x04 // 读输入寄存器
	FuncCodeWriteCoil    = 0x05 // 写单个线圈寄存器
	FuncCodeWriteHold    = 0x06 // 写单个保持寄存器
	FuncCodeWriteCoils   = 0x0F // 写多个线圈
	FuncCodeWriteHolds   = 0x10 // 写多个保持寄存器
)

func FuncCodeToString(b uint8) string {
	switch b {
	case FuncCodeReadCoil:
		return "read coil"
	case FuncCodeReadDiscrete:
		return "read discrete"
	case FuncCodeReadHold:
		return "read hold"
	case FuncCodeReadInput:
		return "read input"
	case FuncCodeWriteCoil:
		return "write coil"
	case FuncCodeWriteHold:
		return "write hold"
	case FuncCodeWriteCoils:
		return "write coils"
	case FuncCodeWriteHolds:
		return "write holds"
	default:
		return "unknown function code"
	}
}

// Modbus 错误码 (Modbus Exception Code)
const (
	ExcepNormal         = 0x00 // 正常 (Normal)
	ExcepIllFuncCode    = 0x01 // 无效功能码 (Invalid Function Code)
	ExcepIllDataAddr    = 0x02 // 无效数据地址 (Invalid Data Address)
	ExcepIllDataValue   = 0x03 // 无效数据值 (Invalid Data Value)
	ExcepSlaveFail      = 0x04 // 从机设备故障 (Slave Device Failure)
	ExcepAck            = 0x05 // 应答 (Acknowledge)
	ExcepSlaveBusy      = 0x06 // 从机设备忙 (Slave Device Busy)
	ExcepNAck           = 0x07 // 非应答 (Negative Acknowledge)
	ExcepMemoryParity   = 0x08 // 存储器校验错 (Memory Parity Error)
	ExcepGwPathUnav     = 0x0A // 网关路径不可用 (Gateway Path Unavailable)
	ExcepGwDevNoRespond = 0x0B // 网关目标设备无应答 (Gateway Target Device No Response)
)

func ExcepToString(b uint8) string {
	switch b {
	case ExcepNormal:
		return "normal"
	case ExcepIllFuncCode:
		return "illegal function code"
	case ExcepIllDataAddr:
		return "illegal data address"
	case ExcepIllDataValue:
		return "illegal data value"
	case ExcepSlaveFail:
		return "slave device failure"
	case ExcepAck:
		return "acknowledge"
	case ExcepSlaveBusy:
		return "slave device busy"
	case ExcepNAck:
		return "negative acknowledge"
	case ExcepMemoryParity:
		return "memory parity error"
	case ExcepGwPathUnav:
		return "gateway path unavailable"
	case ExcepGwDevNoRespond:
		return "gateway target device no response"
	default:
		return "unknown error"
	}
}

// 处理结果 (Process Result)
const (
	ResultNormal       = iota // 正常
	ResultTooShort            // 报文过短
	ResultProtocol            // 协议错误
	ResultFuncCode            // 功能码错误
	ResultDevID               // 设备标识错误
	ResultRegAddr             // 寄存器地址错误
	ResultRegLen              // 寄存器数量错误
	ResultRegValue            // 寄存器值错误
	ResultLength              // 报文长度错误
	ResultRtuCrc              // CRC校验码错误 (Modbus RTU)
	ResultAsciiStart          // 起始符错误 (Modbus Ascii)
	ResultAsciiEnd            // 结束符错误 (Modbus Ascii)
	ResultAsciiLrc            // LRC校验码错误 (Modbus Ascii)
	ResultTcpSerNum           // 流水号错误 (Modbus TCP)
	ResultTcpProtocol         // 协议标识错误 (Modbus TCP)
	ResultBufTooShort         // 缓冲区过短
	ResultUnknownError        // 未知错误
)

// 处理结果错误 (Process Result Error)
type ErrResult struct {
	Code uint8  // 处理结果码
	Zh   string // 中文提示
	Err  error  // 错误
}

func (e *ErrResult) Error() string {
	return e.Err.Error()
}

func (e *ErrResult) Unwrap() error {
	return e.Err
}

// 定义处理结果错误
var (
	ErrResultTooShort     = &ErrResult{Code: ResultTooShort, Zh: "报文过短", Err: errors.New("message too short")}
	ErrResultProtocol     = &ErrResult{Code: ResultProtocol, Zh: "协议错误", Err: errors.New("protocol error")}
	ErrResultFuncCode     = &ErrResult{Code: ResultFuncCode, Zh: "功能码错误", Err: errors.New("function code error")}
	ErrResultDevID        = &ErrResult{Code: ResultDevID, Zh: "设备标识错误", Err: errors.New("device identifier error")}
	ErrResultRegAddr      = &ErrResult{Code: ResultRegAddr, Zh: "寄存器地址错误", Err: errors.New("register address error")}
	ErrResultRegLen       = &ErrResult{Code: ResultRegLen, Zh: "寄存器数量错误", Err: errors.New("register quantity error")}
	ErrResultRegValue     = &ErrResult{Code: ResultRegValue, Zh: "寄存器值错误", Err: errors.New("register value error")}
	ErrResultLength       = &ErrResult{Code: ResultLength, Zh: "报文长度错误", Err: errors.New("message length error")}
	ErrResultRtuCrc       = &ErrResult{Code: ResultRtuCrc, Zh: "CRC校验码错误 (Modbus RTU)", Err: errors.New("CRC check error (Modbus RTU)")}
	ErrResultAsciiStart   = &ErrResult{Code: ResultAsciiStart, Zh: "起始符错误 (Modbus Ascii)", Err: errors.New("start character error (Modbus Ascii)")}
	ErrResultAsciiEnd     = &ErrResult{Code: ResultAsciiEnd, Zh: "结束符错误 (Modbus Ascii)", Err: errors.New("end character error (Modbus Ascii)")}
	ErrResultAsciiLrc     = &ErrResult{Code: ResultAsciiLrc, Zh: "LRC校验码错误 (Modbus Ascii)", Err: errors.New("LRC check error (Modbus Ascii)")}
	ErrResultTcpSerNum    = &ErrResult{Code: ResultTcpSerNum, Zh: "流水号错误 (Modbus TCP)", Err: errors.New("transaction ID error (Modbus TCP)")}
	ErrResultTcpProtocol  = &ErrResult{Code: ResultTcpProtocol, Zh: "协议标识错误 (Modbus TCP)", Err: errors.New("protocol identifier error (Modbus TCP)")}
	ErrResultBufTooShort  = &ErrResult{Code: ResultBufTooShort, Zh: "缓冲区过短", Err: errors.New("buffer too short")}
	ErrResultUnknownError = &ErrResult{Code: ResultUnknownError, Zh: "未知错误", Err: errors.New("unknown error")}
)
