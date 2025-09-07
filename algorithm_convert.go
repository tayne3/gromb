// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gromb

import (
	"encoding/binary"
	"math"
)

const (
	hextable        = "0123456789ABCDEF"
	reverseHexTable = ""
		+ "\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"
		+ "\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"
		+ "\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"
		+ "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\xFF\xFF\xFF\xFF\xFF\xFF"
		+ "\xFF\x0A\x0B\x0C\x0D\x0E\x0F\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"
		+ "\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"
		+ "\xFF\x0A\x0B\x0C\x0D\x0E\x0F\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"
		+ "\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"
		+ "\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"
		+ "\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"
		+ "\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"
		+ "\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"
		+ "\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"
		+ "\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"
		+ "\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"
		+ "\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF\xFF"
)

func asciiToHex(src []uint8) []uint8 {
	dst := make([]uint8, len(src)/2)
	{
		i, j := 0, 0
		for j+1 < len(src) {
			if src[j] == ' ' || src[j] == '\r' || src[j] == '\n' {
				j++
				continue
			}

			a := reverseHexTable[src[j]]
			b := reverseHexTable[src[j+1]]

			if a > 0x0f {
				return nil
			}
			if b > 0x0f {
				return nil
			}

			dst[i] = (a << 4) | b

			i++
			j += 2
		}
	}
	return dst
}

func asciiFromHex(src []uint8) []uint8 {
	dst := make([]uint8, len(src)*2)
	{
		j := 0
		for _, v := range src {
			dst[j] = hextable[v>>4]
			dst[j+1] = hextable[v&0x0f]
			j += 2
		}
	}
	return dst
}

func byteToUint16(in []uint8, order binary.ByteOrder) uint16 {
	return order.Uint16(in[:])
}

func byteFromUint16(in uint16, order binary.ByteOrder) []uint8 {
	var out []uint8 = make([]uint8, 2)
	order.PutUint16(out, in)
	return out
}

func bytesFromUint16s(input []uint16, order binary.ByteOrder) []uint8 {
	output := make([]uint8, 2*len(input))
	for i := 0; i < len(input); i++ {
		order.PutUint16(output[i*2:i*2+2], input[i])
	}
	return output
}

func bytesToUint16s(input []uint8, order binary.ByteOrder) []uint16 {
	output := make([]uint16, len(input)/2)
	for i := 0; i < len(input)/2; i++ {
		output[i] = order.Uint16(input[i*2 : i*2+2])
	}
	return output
}

func bytesToFloat32s(input []uint8, order binary.ByteOrder) []float32 {
	output := make([]float32, len(input)/4)
	for i := 0; i < len(input)/4; i++ {
		output[i] = math.Float32frombits(order.Uint32(input[i*4 : i*4+4]))
	}
	return output
}

func bytesFromFloat32s(input []float32, order binary.ByteOrder) []uint8 {
	output := make([]uint8, len(input)*4)
	for i := 0; i < len(input); i++ {
		order.PutUint32(output[i*4:i*4+4], math.Float32bits(input[i]))
	}
	return output
}