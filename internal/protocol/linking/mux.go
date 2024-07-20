package linking

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync/atomic"

	"github.com/Tifufu/tools-cli/internal/automower"
	"github.com/charmbracelet/log"
)

const (
	DefaultLinkId LinkId = 0
)

var (
	ErrLinkMuxShuttingDown = errors.New("link mux closed")
)

type LinkMux struct {
	device     *automower.Device
	logger     *log.Logger
	inShutdown atomic.Bool
	writeChan  chan []byte

	linkIdCounter atomic.Uint32
	DefaultLink   *Link
}

func NewLinkMux(device *automower.Device, logger *log.Logger) *LinkMux {
	mux := &LinkMux{
		device:        device,
		logger:        logger,
		writeChan:     make(chan []byte),
		inShutdown:    atomic.Bool{},
		linkIdCounter: atomic.Uint32{},
	}
	mux.DefaultLink = &Link{
		id:     DefaultLinkId,
		writer: mux,
	}
	return mux
}

func (lh *LinkMux) Start() error {
	go func() {
		err := writeWorker(lh.writeChan, lh.device)
		if err != nil {
			lh.logger.Error("err from writer worker", "err", err)
		}
	}()

	for {
		rawPacket, err := lh.readFromDevice()
		if err != nil {
			if lh.shuttingDown() {
				return ErrLinkMuxShuttingDown
			}
			return err
		}

		go lh.routePacket(context.TODO(), rawPacket)
	}
}

func (lh *LinkMux) Stop() error {
	lh.inShutdown.Store(true)
	close(lh.writeChan)
	lh.device.Close()
	return nil
}

func (lh *LinkMux) Write(data []byte) (n int, err error) {
	lh.writeChan <- data
	lh.logger.Debug("writing to device", "data", Payload(data))
	return len(data), nil
}

func (lh *LinkMux) shuttingDown() bool {
	return lh.inShutdown.Load()
}

func (lh *LinkMux) readFromDevice() (rawPacket []byte, err error) {
	select {
	case rawPacket = <-lh.device.PacketChan:
		lh.logger.Debug("received from device", "packet", Payload(rawPacket).String())
		return rawPacket, nil
	case err = <-lh.device.ErrChan:
		return nil, err
	}
}

func (lh *LinkMux) routePacket(_ context.Context, rawPacket []byte) {
	switch rawPacket[0] {
	case BroadcastPacketType:
		packet, payload, err := parseBroadcastPacket(rawPacket)
		if err != nil {
			lh.logger.Error("Error parsing packet", "err", err)
			return
		}

		lh.logger.Debug("Parsed broadcast packet", "channel", packet.BroadcastChannel, "control", packet.Control, "payloadSize", len(payload), "senderId", packet.SenderId, "familyId", packet.MessageFamily)
	case LinkedPacketType:
		packet, err := parseLinkedPacket(rawPacket)
		if err != nil {
			lh.logger.Error("Error parsing packet", "err", err)
			return
		}

		lh.logger.Debug("Parsed linked packet", "linkId", packet.LinkId, "control", packet.Control, "payloadSize", len(packet.Payload), "payload", Payload(packet.Payload).String())
	default:
		lh.logger.Debug("Skipping packet, neither linked or broadcast packet", "packet", Payload(rawPacket).String())
	}
}

func writeWorker(input <-chan []byte, out io.Writer) error {
	for data := range input {
		n, err := out.Write(data)
		if err != nil {
			return err
		}
		if n != len(data) {
			return fmt.Errorf("was unable to write entire data, expected to write %d bytes but only wrote %d bytes", len(data), n)
		}
	}
	return nil
}
