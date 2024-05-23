package automower

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

const (
	LinkedPacketType    byte = 0xFD
	BroadcastPacketType byte = 0xFC
)

var (
	ErrUnkownPacketType = errors.New("unknown packet type")
)

/*
structure:
 - Packet type (1 byte)
 - Length (2 bytes, uint16)
 - Channel ID (4 bytes, uint32)
 - Control (1 byte)
 - CRC (1 byte)
Bytes are stored in little-endian order.
*/
type Header struct {
	channelId  uint32
	length     uint16
	packetType byte
	protocol   byte
	headerCrc  byte // TODO Learn more about CRC
}

type Payload []byte

type Packet struct {
	Header  Header
	Payload []byte
}

func newHeader(
	packetType byte,
	length uint16,
	channelId uint32,
	protocolId byte,
	headerCrc byte,
) Header {
	return Header{
		channelId:  channelId,
		length:     length,
		packetType: packetType,
		protocol:   protocolId,
		headerCrc:  headerCrc,
	}
}

func (p Payload) String() string {
	hexStr := hex.EncodeToString(p)
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}

	var hexArray []string
	for i := 0; i < len(hexStr); i += 2 {
		hexArray = append(hexArray, hexStr[i:i+2])
	}

	return strings.ToUpper(strings.Join(hexArray, " "))
}

func ParsePacket(data []byte) (Packet, error) {
	buff := bytes.NewBuffer(data)

	packetType, err := buff.ReadByte()
	if err != nil {
		return Packet{}, err
	}
	switch packetType {
	case LinkedPacketType:
	case BroadcastPacketType:
	default:
		return Packet{}, errors.Join(ErrUnkownPacketType, fmt.Errorf("packet type: %d", packetType))
	}

	packetLengthBytes := make([]byte, 2, 2)
	n, err := buff.Read(packetLengthBytes)
	if err != nil {
		return Packet{}, err
	}
	if n != 2 {
		return Packet{}, fmt.Errorf("expected 2 bytes for packet length, got %d", n)
	}
	packetLength := binary.LittleEndian.Uint16(packetLengthBytes)

	channelIdBytes := make([]byte, 4, 4)
	n, err = buff.Read(channelIdBytes)
	if err != nil {
		return Packet{}, err
	}
	if n != 4 {
		return Packet{}, fmt.Errorf("expected 4 bytes for channel ID, got %d", n)
	}
	channelId := binary.LittleEndian.Uint32(channelIdBytes)

	protocolId, err := buff.ReadByte()
	if err != nil {
		return Packet{}, err
	}

	headerCrc, err := buff.ReadByte()
	if err != nil {
		return Packet{}, err
	}

	return Packet{
		Header: newHeader(
			packetType,
			packetLength,
			channelId,
			protocolId,
			headerCrc,
		),
		Payload: buff.Bytes(),
	}, nil
}
