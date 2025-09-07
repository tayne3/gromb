// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gromb

import (
	"encoding/binary"
)

//	<--------------- Max(1) -------------->
//	<---- Buffer Size ---->
//	+----------+----------+---------------+ 
//	|   Last   |   This   |   Available   |
//	+----------+----------+---------------+ 
//	     |          |             |
//	    (2)        (3)           (4)
//
//	(1)  ... max 是指所有操作的最大空间
//	(2)  ... last 是指此次操作前的已用空间
//	(3)  ... this 是此次操作的使用空间
//	(4)  ... available 是指此次操作后的剩余空间

// 缓冲盒子
type groBox struct {
	buffer *[]uint8 // 缓冲区
	last   uint16   // 已用空间
	max    uint16   // 最大空间 (默认为 1024)
}

// 初始化
func (b *groBox) Init(buffer *[]uint8, max uint16) {
	b.buffer = buffer
	b.last = 0
	b.max = max
}

// 重置
func (b *groBox) Reset() {
	b.Init(nil, 0)
}

// 清空元素
func (b *groBox) Clear() {
	(*b.buffer) = make([]uint8, 0)
	b.last = 0
	b.max = 0
}

// 获取 最大空间
func (b *groBox) GetMax() uint16 {
	return b.max
}

// 设置 最大空间
func (b *groBox) SetMax(max uint16) {
	b.max = max
}

// 获取 已用空间
func (b *groBox) GetLast() uint16 {
	return b.last
}

// 设置 已用空间
func (b *groBox) SetLast(last uint16) {
	b.last = last
}

// 增加 已用空间
func (b *groBox) AddLast(n uint16) {
	b.last += n
}

// 减少 已用空间
func (b *groBox) SubLast(n uint16) {
	b.last -= n
}

// // 获取 使用空间
// func (b *groBox) LastSize() uint16 {
// 	return uint16(len(*b.buffer)) - b.last
// }

// 获取 使用空间
func (b *groBox) ThisSize() uint16 {
	return uint16(len(*b.buffer)) - b.last
}

// 获取 总长度
func (b *groBox) Size() uint16 {
	return uint16(len(*b.buffer))
}

// 获取 可用空间
func (b *groBox) Available() uint16 {
	return b.max - uint16(len(*b.buffer)) - b.last
}

// 获取 总缓冲区
func (b *groBox) GetBuffer(start uint16, end uint16) []uint8 {
	return (*b.buffer)[start:end]
}

// 获取 缓冲区
func (b *groBox) GetThisBuffer(start uint16, end uint16) []uint8 {
	return (*b.buffer)[b.last+start : b.last+end]
}

// 获取一个 uint8
func (b *groBox) GetU8(offset uint16) uint8 {
	return (*b.buffer)[b.last+offset]
}

// 获取一个 uint16
func (b *groBox) GetU16(offset uint16, order binary.ByteOrder) uint16 {
	return order.Uint16((*b.buffer)[b.last+offset : b.last+offset+2])
}

// 获取多个 uint8
func (b *groBox) GetU8s(offset uint16, end uint16) []uint8 {
	return (*b.buffer)[b.last+offset : b.last+end]
}

// 获取多个 uint16
func (b *groBox) GetU16s(offset uint16, end uint16, order binary.ByteOrder) []uint16 {
	return bytesToUint16s((*b.buffer)[b.last+offset:b.last+end*2], order)
}

// 设置一个 uint8
func (b *groBox) SetU8(offset uint16, value uint8) {
	(*b.buffer)[b.last+offset] = value
}

// 设置一个 uint16
func (b *groBox) SetU16(offset uint16, value uint16, order binary.ByteOrder) {
	order.PutUint16((*b.buffer)[b.last+offset:b.last+offset+2], value)
}

// 尾部插入 uint8
func (b *groBox) PutU8(value uint8) {
	(*b.buffer) = append((*b.buffer), value)
}

// 尾部插入 uint16
func (b *groBox) PutU16(value uint16, order binary.ByteOrder) {
	(*b.buffer) = append((*b.buffer), byteFromUint16(value, order)...)
}

// 尾部插入 []uint8
func (b *groBox) PutU8s(arg []uint8) {
	(*b.buffer) = append((*b.buffer), arg...)
}

// 尾部插入 []uint16
func (b *groBox) PutU16s(arg []uint16, order binary.ByteOrder) {
	(*b.buffer) = append((*b.buffer), bytesFromUint16s(arg, order)...)
}

// 设置 []uint8
func (b *groBox) SetU8s(arg []uint8) {
	(*b.buffer) = make([]uint8, len(arg))
	copy((*b.buffer), arg)
}

// 设置 []uint16
func (b *groBox) SetU16s(arg []uint16, order binary.ByteOrder) {
	(*b.buffer) = make([]uint8, len(arg)*2)
	copy((*b.buffer), bytesFromUint16s(arg, order))
}