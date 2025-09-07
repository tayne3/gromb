<p align="center" style="font-size: 24px; font-weight: bold; color: #797979;">Gromb</p>

**English** | [中文](README_zh.md)

Gromb is a Modbus protocol library implemented in the Go language. It supports multiple Modbus protocols and data transmission methods, including RTU, ASCII, and TCP, and can communicate with various Modbus master and slave devices.

## ✨ Features

- Supports Modbus RTU, ASCII, and TCP protocols
- Can be used as a Modbus master or slave
- Supports custom register access control
- Supports common function codes such as read/write coils, input/holding registers

## 🚀 Quick Start

### Example Usage

Example pseudo-code:

```go
package main

import (
	"fmt"
	"gromb"
)

func main() {
    // Create a Modbus protocol stack instance
    m := gromb.New()

    // Initialize RTU protocol header
    m.Head.InitRtu(0x01)

    // Set function code and register address
    m.Arg.SetFuncCode(gromb.FuncCodeReadHold)
    m.Arg.SetRegAddr(0x0000)
    m.Arg.SetRegLen(0x000A)

    // Pack the request message
    reqBytes, ok := m.PackRequest()
    if !ok {
    	fmt.Println("pack request failed:", m.Result.GetResultString())
    	return
    }

    // Send the request message
    Write(reqBytes)

    // Receive the response message
    rspBytes := Read()

    // Parse the response message
    if ok := m.ParseResponse(rspBytes); !ok {
    	fmt.Println("parse response failed:", m.Result.GetResultString())
    	return
    }

    // Get the response data
    arg := m.Arg.GetU16s(binary.BigEndian)
    fmt.Println("response arg:", arg)
}
```
