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

type LinkId uint32

type Link struct {
	id LinkId
}

func (l Link) Id() LinkId {
	return l.id
}

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
	defaultLink   *Link
}

func NewLinkMux(device *automower.Device, logger *log.Logger) *LinkMux {
	return &LinkMux{
		device:        device,
		logger:        logger,
		writeChan:     make(chan []byte),
		inShutdown:    atomic.Bool{},
		linkIdCounter: atomic.Uint32{},
		defaultLink:   &Link{id: 0},
	}
}

func (lh *LinkMux) Start() error {
	go func() {
		err := writeWorker(lh.writeChan, lh.device)
		if err != nil {
			lh.logger.Error("err is writer worker", "err", err)
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

		go lh.routePacket(context.Background(), rawPacket)
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
	lh.logger.Debug("writing to device", "data", data)
	return len(data), nil
}

func (lh *LinkMux) shuttingDown() bool {
	return lh.inShutdown.Load()
}

func (lh *LinkMux) readFromDevice() (rawPacket []byte, err error) {
	select {
	case rawPacket = <-lh.device.PacketChan:
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

		lh.logger.Debug("Parsed linked packet", "linkId", packet.LinkId, "control", packet.Control, "payloadSize", len(packet.Payload))
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
			return fmt.Errorf("was unable to write entire data expected to write %d bytes but wrote %d bytes", len(data), n)
		}
	}
	return nil
}
