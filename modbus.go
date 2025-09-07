package gromb

type Modbus struct {
	Arg    groArg    // 处理寄存器值
	Access groAccess // 数据访问控制器
	Result groResult // 处理参数
	Head   groHead   // 协议头参数
	Box    groBox    // 处理报文盒子
}

func New() *Modbus {
	m := &Modbus{}
	m.Reset()
	return m
}

func (m *Modbus) Reset() {
	m.Arg.Reset()
	m.Access.Reset()
	m.Result.Reset()
	m.Head.Reset()
	m.Box.Reset()
}

func (m *Modbus) PackRequest(b []uint8) error {
	b = b[:0]
	m.Box.Init(&b, 1024)

	switch m.Head.GetProtocol() {
	case ProtocolRTU:
		m.rtuPack(true)
	case ProtocolAscii:
		m.asciiPack(true)
	case ProtocolTCP:
		m.tcpPack(true)
	default:
		m.Result.SetResult(ErrResultProtocol)
	}
	return m.Result.GetResult()
}

func (m *Modbus) PackResponse(b []uint8) error {
	b = b[:0]
	m.Box.Init(&b, 1024)

	switch m.Head.GetProtocol() {
	case ProtocolRTU:
		m.rtuPack(false)
	case ProtocolAscii:
		m.asciiPack(false)
	case ProtocolTCP:
		m.tcpPack(false)
	default:
		m.Result.SetResult(ErrResultProtocol)
	}
	return m.Result.GetResult()
}

func (m *Modbus) ParseRequest(b []uint8) error {
	m.Box.Init(&b, uint16(len(b)))
	m.Result.Reset()

	switch m.Head.GetProtocol() {
	case ProtocolRTU:
		m.rtuParse(true)
	case ProtocolAscii:
		m.asciiParse(true)
	case ProtocolTCP:
		m.tcpParse(true)
	default:
		m.Result.SetResult(ErrResultProtocol)
	}
	return m.Result.GetResult()
}

func (m *Modbus) ParseResponse(b []uint8) error {
	m.Box.Init(&b, uint16(len(b)))
	m.Result.Reset()

	switch m.Head.GetProtocol() {
	case ProtocolRTU:
		m.rtuParse(false)
	case ProtocolAscii:
		m.asciiParse(false)
	case ProtocolTCP:
		m.tcpParse(false)
	default:
		m.Result.SetResult(ErrResultProtocol)
	}
	return m.Result.GetResult()
}
