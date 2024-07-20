package linking

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseLinkedPacket(t *testing.T) {
	t.Parallel()
	input := []byte{0xFD, 0x0D, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x63, 0x12, 0x5F, 0x4D, 0xB6, 0x16, 0x74}
	expected := linkedPacket{
		PacketType: LinkedPacketType,
		Length:     13,
		LinkId:     LinkId(0),
		Control:    ControlLinkManagerCommand,
		HeaderCrc:  0x63,
		Payload:    []byte{0x12, 0x5F, 0x4D, 0xB6, 0x16},
		FooterCrc:  0x74,
	}

	packet, err := parseLinkedPacket(input)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(expected.Payload, packet.Payload) {
		t.Errorf("payload mismatch expected %v but got %v", expected.Payload, packet.Payload)
	}
	if !cmp.Equal(expected, packet) {
		t.Errorf("expected %v but got %v", expected, packet)
	}
}

func TestMarshallLinkedPacket(t *testing.T) {
	t.Parallel()
	input := linkedPacket{
		PacketType: LinkedPacketType,
		Length:     13,
		LinkId:     LinkId(0),
		Control:    ControlLinkManagerCommand,
		HeaderCrc:  0x63,
		Payload:    []byte{0x12, 0x5F, 0x4D, 0xB6, 0x16},
		FooterCrc:  0x74,
	}
	expected := []byte{0xFD, 0x0D, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x63, 0x12, 0x5F, 0x4D, 0xB6, 0x16, 0x74}

	buf := make([]byte, linkedPacketSize(len(input.Payload)))
	err := marshallLinkedPacket(input, buf)
	if err != nil {
		t.Errorf("expected no error but got %q", err)
	}
	if !bytes.Equal(buf, expected) {
		t.Errorf("marshallLinkedPacket(%+v) returned %v; expected %v", input, buf, expected)
	}
}
