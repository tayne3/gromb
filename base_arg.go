package gromb

import (
	"encoding/binary"
)

// Modbus 处理套件-处理值
type groArg struct {
	funccode uint8   // 功能码
	regaddr  uint16  // 寄存器地址
	reglen   uint16  // 寄存器长度
	all      []uint8 // 寄存器值
}

func (a *groArg) Init(funccode uint8, regaddr, reglen uint16) {
	a.funccode = funccode
	a.regaddr = regaddr
	a.reglen = reglen
	a.all = nil
}

func (a *groArg) Reset() {
	a.all = nil
}

func (a *groArg) SetFuncCode(funccode uint8) {
	a.funccode = funccode
}

func (a *groArg) SetRegAddr(regaddr uint16) {
	a.regaddr = regaddr
}

func (a *groArg) SetRegLen(reglen uint16) {
	a.reglen = reglen
}

func (a *groArg) GetFuncCode() uint8 {
	return a.funccode
}

func (a *groArg) GetRegAddr() uint16 {
	return a.regaddr
}

func (a *groArg) GetRegLen() uint16 {
	return a.reglen
}

func (a *groArg) GetFuncCodeString() string {
	return FuncCodeToString(a.funccode)
}

func (a *groArg) GetU8(offset int) uint8 {
	return a.all[offset]
}

func (a *groArg) GetU16(offset int, order binary.ByteOrder) uint16 {
	return byteToUint16(a.all[offset:offset+2], order)
}

// func (a *groArg) SetU8(offset int, data uint8) {
// 	a.all[offset] = data
// }

// func (a *groArg) SetU16(offset int, data uint16, order binary.ByteOrder) {
// 	order.PutUint16(a.all[offset:offset+2], data)
// }

func (a *groArg) GetU8s() []uint8 {
	return a.all
}

func (a *groArg) GetU16s(order binary.ByteOrder) []uint16 {
	return bytesToUint16s(a.all, order)
}

func (a *groArg) GetFloat32s(order binary.ByteOrder) []float32 {
	return bytesToFloat32s(a.all, order)
}

func (a *groArg) SetBits(data []bool) {
	sum := len(data)
	number := (sum + 7) / 8
	a.all = make([]uint8, number)

	s := 0
	for n := 0; s < sum; n++ {
		a.all[n] = 0x00
		for b := 0; b < 8 && s < sum; {
			if data[s] {
				a.all[n] |= 1 << b
			}
			b++
			s++
		}
	}
}

func (a *groArg) GetBits() []bool {
	sum := len(a.all) * 8
	data := make([]bool, sum)

	s := 0
	for n := 0; s < sum; n++ {
		for b := 0; b < 8 && s < sum; b++ {
			if (a.all[n] & (1 << b)) != 0 {
				data[s] = true
			} else {
				data[s] = false
			}
			s++
		}
	}
	return data
}

func (a *groArg) SetU8s(data []uint8) {
	a.all = data
}

func (a *groArg) SetU16s(data []uint16, order binary.ByteOrder) {
	a.all = bytesFromUint16s(data, order)
}

func (a *groArg) SetFloat32s(data []float32, order binary.ByteOrder) {
	a.all = bytesFromFloat32s(data, order)
}
