package linking

import "io"

type LinkId uint32

const ()

type Link struct {
	id     LinkId
	writer io.Writer
}

func (l Link) Id() LinkId {
	return l.id
}

func (l Link) write(data []byte) (n int, err error) {
	return l.writer.Write(data)
}

func (l Link) SendLinkedRequest(ctrl byte, payload []byte) (<-chan []byte, error) {
	payloadSize := len(payload)
	length := 1 + uint16(linkIdSize+controlSize+headerCrcSize+payloadSize+footerCrcSize)
	packet := linkedPacket{
		PacketType: LinkedPacketType,
		Length:     length,
		LinkId:     l.id,
		Control:    ctrl,
		HeaderCrc:  0x63,
		Payload:    payload,
		FooterCrc:  0x74,
	}

	packetSize := linkedPacketSize(payloadSize)
	buf := make([]byte, packetSize)
	err := marshallLinkedPacket(packet, buf)
	if err != nil {
		return nil, err
	}

	const headerCrcIdx int = 8 // packetTypeSize + lengthSize + linkIdSize + controlSize
	headerCrc := calcCrc8(buf[0:headerCrcIdx])
	buf[headerCrcIdx] = headerCrc // headerCrc comes right after control byte

	footerCrc := calcCrc8(buf[0 : packetSize-1])
	buf[packetSize-1] = footerCrc

	packetBuf := make([]byte, packetSize+2)
	packetBuf[0] = 0x02
	copy(packetBuf[1:], buf)
	packetBuf[packetSize+1] = 0x03

	_, err = l.write(packetBuf)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
