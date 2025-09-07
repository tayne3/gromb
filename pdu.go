// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gromb

// 封装 PDU 报文
// 返回值: >=0-PDU报文长度,-1-失败
func Pack(result *groResult, arg *groArg, box *groBox, isReq bool) int {
	// 检查是否回复异常
	if !isReq && result.GetExcepCode() != ExcepNormal {
		return packResponseExcep(result, arg, box)
	}

	switch arg.GetFuncCode() {
	case FuncCodeReadCoil:
		if isReq {
			return packRequestReadCoil(result, arg, box)
		} else {
			return packResponseReadCoil(result, arg, box)
		}
	case FuncCodeWriteCoil:
		if isReq {
			return packRequestWriteCoil(result, arg, box)
		} else {
			return packResponseWriteCoil(result, arg, box)
		}
	case FuncCodeWriteCoils:
		if isReq {
			return packRequestWriteCoils(result, arg, box)
		} else {
			return packResponseWriteCoils(result, arg, box)
		}
	case FuncCodeReadDiscrete:
		if isReq {
			return packRequestReadDiscretes(result, arg, box)
		} else {
			return packResponseReadDiscretes(result, arg, box)
		}
	case FuncCodeReadHold:
		if isReq {
			return packRequestReadHold(result, arg, box)
		} else {
			return packResponseReadHold(result, arg, box)
		}
	case FuncCodeWriteHold:
		if isReq {
			return packRequestWriteHold(result, arg, box)
		} else {
			return packResponseWriteHold(result, arg, box)
		}
	case FuncCodeWriteHolds:
		if isReq {
			return packRequestWriteHolds(result, arg, box)
		} else {
			return packResponseWriteHolds(result, arg, box)
		}
	case FuncCodeReadInput:
		if isReq {
			return packRequestReadInput(result, arg, box)
		} else {
			return packResponseReadInput(result, arg, box)
		}
	default:
		result.SetResult(ErrResultFuncCode)
		return -1
	}
}

// 解析 PDU 报文
// 返回值: >=0-PDU报文长度,-1-失败
func Parse(result *groResult, access *groAccess, arg *groArg, box *groBox, isReq bool) int {
	funccode := box.GetU8(0)
	arg.SetFuncCode(funccode)

	switch funccode {
	case FuncCodeReadCoil:
		if isReq {
			return parseRequestReadCoil(result, access, arg, box)
		} else {
			return parseResponseReadCoil(result, arg, box)
		}
	case FuncCodeWriteCoil:
		if isReq {
			return parseRequestWriteCoil(result, access, arg, box)
		} else {
			return parseResponseWriteCoil(result, arg, box)
		}
	case FuncCodeWriteCoils:
		if isReq {
			return parseRequestWriteCoils(result, access, arg, box)
		} else {
			return parseResponseWriteCoils(result, arg, box)
		}
	case FuncCodeReadDiscrete:
		if isReq {
			return parseRequestReadDiscretes(result, access, arg, box)
		} else {
			return parseResponseReadDiscretes(result, arg, box)
		}
	case FuncCodeReadHold:
		if isReq {
			return parseRequestReadHold(result, access, arg, box)
		} else {
			return parseResponseReadHold(result, arg, box)
		}
	case FuncCodeWriteHold:
		if isReq {
			return parseRequestWriteHold(result, access, arg, box)
		} else {
			return parseResponseWriteHold(result, arg, box)
		}
	case FuncCodeWriteHolds:
		if isReq {
			return parseRequestWriteHolds(result, access, arg, box)
		} else {
			return parseResponseWriteHolds(result, arg, box)
		}
	case FuncCodeReadInput:
		if isReq {
			return parseRequestReadInput(result, access, arg, box)
		} else {
			return parseResponseReadInput(result, arg, box)
		}
	case FuncCodeReadCoil | 0x80,
		FuncCodeWriteCoil | 0x80,
		FuncCodeWriteCoils | 0x80,
		FuncCodeReadDiscrete | 0x80,
		FuncCodeReadHold | 0x80,
		FuncCodeWriteHold | 0x80,
		FuncCodeWriteHolds | 0x80,
		FuncCodeReadInput | 0x80:
		if isReq {
			result.SetResult(ErrResultFuncCode)
			return -1
		} else {
			return parseResponseExcep(result, box)
		}
	default:
		result.SetResult(ErrResultFuncCode)
		return -1
	}
}
