package gromb

type groHead struct {
	protocol uint8  // 协议类型
	devid    uint8  // 设备标识
	sernum   uint16 // 序列号
}

func (h *groHead) InitRtu(devid uint8) {
	h.protocol = ProtocolRTU
	h.devid = devid
	h.sernum = 0
}

func (h *groHead) InitAscii(devid uint8) {
	h.protocol = ProtocolAscii
	h.devid = devid
	h.sernum = 0
}

func (h *groHead) InitTcp(devid uint8, sernum uint16) {
	h.protocol = ProtocolTCP
	h.devid = devid
	h.sernum = sernum
}

func (h *groHead) Reset() {
	h.InitRtu(0)
}

func (h *groHead) SetProtocol(protocol uint8) {
	h.protocol = protocol
}

func (h *groHead) SetDevId(devid uint8) {
	h.devid = devid
}

func (h *groHead) SetSerNum(sernum uint16) {
	h.sernum = sernum
}

func (h *groHead) IncSerNum() {
	h.sernum++
}

func (h *groHead) GetProtocol() uint8 {
	return h.protocol
}

func (h *groHead) GetDevId() uint8 {
	return h.devid
}

func (h *groHead) GetSerNum() uint16 {
	return h.sernum
}

func (h *groHead) GetProtocolString() string {
	return ProtocolToString(h.protocol)
}
