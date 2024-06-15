package linking

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

const (
	RoboticsProtocol1 byte = 0x00
	RoboticsProtocol2 byte = 0x01

	PacketStart byte = 0x02
	PacketEnd   byte = 0x03

	LinkedPacketType    byte = 0xFD
	BroadcastPacketType byte = 0xFC

	ControlLinkManagerCommand byte = 0x00
	ControlPayloadCommand     byte = 0x01

	// Size of frame in bytes = Header+Footer.
	linkedPacketFrameSize int = 10
	// Size of frame in bytes = Header+Footer.
	// broadcastFrameSize int = 10
)

var (
	ErrUnkownPacketType = errors.New("unknown packet type")
)

type Payload []byte

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

// https://confluence-husqvarna.riada.se/pages/viewpage.action?pageId=59116753
type linkedPacket struct {
	// Packet type. Should always be 0xFD.
	PacketType byte

	// The Length of the packet starting from the byte
	// after the Length field and including the end-of-packet byte.
	// The end-of-packet byte is not included in the LinkedPacket struct.
	Length uint16

	// Link id is a unique id of a logical link.
	// When a link is created, a random number is generated.
	//
	// Link id 0 (Zero) means local node, i.e. the packet shall not be routed.
	LinkId LinkId

	// The Control byte decides if this packet contains a link manager command, or if it contains payload data.
	//
	// 0x00 - This is a link manager command.
	//
	// 0x01 - Payload data is simply transferred to the destination specified by the channel id, without any interpretation by the link manager.
	Control byte

	// The header CRC is used to verify the integrity of the header data.
	//
	// The recommended approach to receiving a packet
	// is to read the header first, and verify the integrity
	// using the header CRC. If the CRC is correct, the length
	// can be trusted and the rest of the packet can be read.
	//
	// Calculated starting on PacketType and ending on (including) Control.
	HeaderCrc byte

	Payload Payload

	// Checksum that is used to verify the integrity of the packet.
	// The checksum calculation starts after the sync word and runs until the end of the payload.
	FooterCrc byte
}

func NewLinkedPacket(
	packetType byte,
	length uint16,
	linkId LinkId,
	control byte,
	headerCrc byte,
	payload Payload,
	footerCrc byte,
) linkedPacket {
	return linkedPacket{
		PacketType: packetType,
		Length:     length,
		LinkId:     linkId,
		Control:    control,
		HeaderCrc:  headerCrc,
		Payload:    payload,
		FooterCrc:  footerCrc,
	}
}

func marshallLinkedPacket(lp linkedPacket, buf []byte) error {
	if len(buf) < len(lp.Payload)+linkedPacketFrameSize {
		return fmt.Errorf("not enough space in buf for linked packet; buf len %v", len(buf))
	}

	buf[0] = lp.PacketType
	binary.LittleEndian.PutUint16(buf[1:3], lp.Length)
	binary.LittleEndian.PutUint32(buf[3:7], uint32(lp.LinkId))
	copy(buf[7:9], []byte{lp.Control, lp.HeaderCrc})
	copy(buf[9:], lp.Payload)
	buf[len(buf)-1] = lp.FooterCrc

	return nil
}

// todo: Check heap allocations of this function
func parseLinkedPacket(data []byte) (packet linkedPacket, err error) {
	if len(data) < linkedPacketFrameSize {
		return linkedPacket{}, fmt.Errorf("malformed packet, not enough bytes for header+footer")
	}

	buff := bytes.NewBuffer(data)

	packetType, _ := buff.ReadByte()
	if packetType != LinkedPacketType {
		return packet, fmt.Errorf("invalid packet type, expected LinkedPacket %b but got %b", LinkedPacketType, packetType)
	}

	packetLength := binary.LittleEndian.Uint16(buff.Next(2))
	if buff.Len() != int(packetLength)-1 {
		return packet, fmt.Errorf("malformed packet, length header (%d) does not match size of data (%d - end-of-packet)", buff.Len(), packetLength)
	}

	linkId := binary.LittleEndian.Uint32(buff.Next(4))
	control, _ := buff.ReadByte()
	headerCrc, _ := buff.ReadByte()
	payload := buff.Next(buff.Len() - 1) // - 1 byte for footerCrc
	footerCrc, _ := buff.ReadByte()

	packet = linkedPacket{
		PacketType: packetType,
		Length:     packetLength,
		LinkId:     LinkId(linkId),
		Control:    control,
		HeaderCrc:  headerCrc,
		Payload:    payload,
		FooterCrc:  footerCrc,
	}
	return packet, nil
}

type broadcastPacketFrame struct {
	// Packet type. Should always be 0xFC.
	PacketType byte

	// The Length of the packet starting from the byte
	// after the Length field and including the end-of-packet byte.
	// The end-of-packet byte is not included in the BroadcastPacket struct.
	Length uint16

	// The family id of the message.
	// Matches the family id used in
	// Robotics Protocol.
	MessageFamily uint16

	// Static id of the node that sent this message.
	// Used to identify the sender of the message.
	// A node without an id is not allowed to send broadcast
	// packets, since there is no way to guarantee that the header
	// is unique.
	SenderId byte

	BroadcastChannel byte

	// The Control byte decides if this packet contains a link manager command, or if it contains payload data.
	//
	// 0x00 - This is a link manager command.
	//
	// 0x01 - Payload data is simply transferred to the destination specified by the channel id, without any interpretation by the link manager.
	Control byte

	// The header CRC is used to verify the integrity of the header data.
	//
	// The recommended approach to receiving a packet
	// is to read the header first, and verify the integrity
	// using the header CRC. If the CRC is correct, the length
	// can be trusted and the rest of the packet can be read.
	//
	// Calculated starting on PacketType and ending on (including) Control.
	HeaderCrc byte

	// Checksum that is used to verify the integrity of the packet.
	// The checksum calculation starts after the sync word and runs until the end of the payload.
	FooterCrc byte
}

func parseBroadcastPacket(data []byte) (packet broadcastPacketFrame, payload Payload, err error) {
	if len(data) < linkedPacketFrameSize {
		return broadcastPacketFrame{}, nil, fmt.Errorf("malformed packet, not enough bytes for header+footer")
	}

	buff := bytes.NewBuffer(data)

	packetType, _ := buff.ReadByte()
	if packetType != BroadcastPacketType {
		return packet, payload, fmt.Errorf("invalid packet type, expected BroadcastPacket %b but got %b", BroadcastPacketType, packetType)
	}

	packetLength := binary.LittleEndian.Uint16(buff.Next(2))
	if buff.Len() != int(packetLength)-1 {
		return packet, payload, fmt.Errorf("malformed packet, length header (%d) does not match size of data (%d - end-of-packet)", buff.Len(), packetLength)
	}

	messageFamily := binary.LittleEndian.Uint16(buff.Next(2))
	senderId, _ := buff.ReadByte()
	broadcastChannel, _ := buff.ReadByte()
	control, _ := buff.ReadByte()
	headerCrc, _ := buff.ReadByte()
	payload = buff.Next(buff.Len() - 1) // - 1 byte for footerCrc
	footerCrc, _ := buff.ReadByte()

	packet = broadcastPacketFrame{
		PacketType:       packetType,
		Length:           packetLength,
		MessageFamily:    messageFamily,
		SenderId:         senderId,
		BroadcastChannel: broadcastChannel,
		Control:          control,
		HeaderCrc:        headerCrc,
		FooterCrc:        footerCrc,
	}
	return packet, payload, nil
}
