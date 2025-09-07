// Copyright 2025 The Gromb Authors. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gromb

// Modbus处理套件-处理结果
type groResult struct {
	excepCode uint8  // 异常码
	errResult error  // 处理结果
	retlen    uint16 // 封装/解析实际长度
}

func (r *groResult) Reset() {
	r.errResult = nil
	r.excepCode = ExcepNormal
	r.retlen = 0
}

func (r *groResult) SetExcepCode(code uint8) {
	r.excepCode = code
}

func (r *groResult) SetResult(errResult error) {
	r.errResult = errResult
}

func (r *groResult) SetRetLen(retlen uint16) {
	r.retlen = retlen
}

func (r *groResult) GetExcepCode() uint8 {
	return r.excepCode
}

func (r *groResult) GetResult() error {
	return r.errResult
}

func (r *groResult) GetRetLen() uint16 {
	return r.retlen
}

func (r *groResult) GetExcepCodeString() string {
	return ExcepToString(r.excepCode)
}