package automower

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

const (
	packetStart = 0x02
	packetEnd   = 0x03
)

type Device struct {
	stream     io.ReadWriteCloser
	ErrChan    chan error
	PacketChan chan []byte
}

func NewDevice(stream io.ReadWriteCloser, ctx context.Context) *Device {
	device := &Device{
		stream:     stream,
		ErrChan:    make(chan error, 1),
		PacketChan: make(chan []byte),
	}
	go device.watch(ctx)
	return device
}

func (d Device) Write(b []byte) (int, error) {
	return d.stream.Write(b)
}

func (d *Device) watch(ctx context.Context) {
	bufReader := bufio.NewReader(d.stream)
	for {
		_, err := bufReader.ReadBytes(packetStart)
		if err != nil {
			d.ErrChan <- fmt.Errorf("failed to read packet start: %w", err)
			return
		}

		packet, err := bufReader.ReadBytes(packetEnd)
		if err != nil {
			d.ErrChan <- fmt.Errorf("failed to read to end of packet end: %w", err)
			return
		}
		if len(packet) < 2 {
			continue
		}

		packet = packet[:len(packet)-1]

		select {
		case d.PacketChan <- packet:
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (d Device) Close() error {
	if d.stream == nil {
		return nil
	}

	err := d.stream.Close()
	if err != nil {
		return fmt.Errorf("failed to close device stream: %w", err)
	}

	return nil
}
