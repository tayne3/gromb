<p align="center" style="font-size: 24px; font-weight: bold; color: #797979;">Gromb</p>

**English** | [中文](README_zh.md)

Gromb是一个基于Go语言实现的Modbus协议库。它支持多种Modbus协议和数据传输方式，包括RTU、ASCII和TCP，可以与各种Modbus主站和从站设备进行通信。

## ✨ 功能特性

-   支持Modbus RTU、ASCII和TCP协议
-   可作为Modbus主站或从站使用
-   支持自定义寄存器访问控制
-   支持读/写线圈、输入/保持寄存器等常用功能码

## 🚀快速开始

### 使用示例

示例伪代码:

```go
package main

import (
	"fmt"
	"gromb"
)

func main() {
    // 创建Modbus协议栈实例
    m := gromb.New()

    // 初始化RTU协议头
    m.Head.InitRtu(0x01)

    // 设置功能码和寄存器地址
    m.Arg.SetFuncCode(gromb.FuncCodeReadHold)
    m.Arg.SetRegAddr(0x0000)
    m.Arg.SetRegLen(0x000A)

    // 打包请求报文
    reqBytes, ok := m.PackRequest()
    if !ok {
    	fmt.Println("pack request failed:", m.Result.GetResultString())
    	return
    }

    // 发送请求报文
    Write(reqBytes)

    // 接收响应报文
    rspBytes := Read()

    // 解析响应报文
    if ok := m.ParseResponse(rspBytes); !ok {
    	fmt.Println("parse response failed:", m.Result.GetResultString())
    	return
    }

    // 获取响应数据
    arg := m.Arg.GetU16s(binary.BigEndian)
    fmt.Println("response arg:", arg)
}
```

