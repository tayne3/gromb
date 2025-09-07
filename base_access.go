// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gromb

// 过滤请求的设备ID
type AccessFilter func(devid uint8, userdata any) bool

// 检查请求的寄存器地址是否运行读/写
type AccessCheck func(regaddr, reglen uint16, isRead bool, userdata any) bool

// / Modbus-数据访问控制器
type groAccess struct {
	UserData      any          // 用户数据
	FilterDevID   AccessFilter // 过滤请求的设备ID
	CheckCoil     AccessCheck  // 检查函数-线圈状态
	CheckDiscrete AccessCheck  // 检查函数-离散量输入
	CheckHold     AccessCheck  // 检查函数-保持寄存器
	CheckInput    AccessCheck  // 检查函数-输入寄存器
}

func (a *groAccess) Reset() {
	a.UserData = nil
	a.FilterDevID = nil
	a.CheckCoil = nil
	a.CheckDiscrete = nil
	a.CheckHold = nil
	a.CheckInput = nil
}

func (a *groAccess) SetUserData(UserData any) {
	a.UserData = UserData
}

func (a *groAccess) SetFilterDevID(FilterDevID AccessFilter) {
	a.FilterDevID = FilterDevID
}

func (a *groAccess) SetCheckCoil(CheckCoil AccessCheck) {
	a.CheckCoil = CheckCoil
}

func (a *groAccess) SetCheckDiscrete(CheckDiscrete AccessCheck) {
	a.CheckDiscrete = CheckDiscrete
}

func (a *groAccess) SetCheckHold(CheckHold AccessCheck) {
	a.CheckHold = CheckHold
}

func (a *groAccess) SetCheckInput(CheckInput AccessCheck) {
	a.CheckInput = CheckInput
}
