package linking

import (
	"context"
	"time"

	"github.com/Tifufu/tools-cli/internal/automower"
	"github.com/charmbracelet/log"
)

type LinkMux struct {
	stopChan chan struct{}
	device   *automower.Device
	logger   *log.Logger
}

func NewLinkHost(device *automower.Device, logger *log.Logger) *LinkMux {
	return &LinkMux{
		device:   device,
		logger:   logger,
		stopChan: make(chan struct{}, 1),
	}
}

func (lh *LinkMux) Start(ctx context.Context) error {
	for {
		select {
		case rawPacket := <-lh.device.PacketChan:
			if len(rawPacket) <= 0 {
				continue
			}

			routeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			go lh.routePacket(routeCtx, rawPacket)
		case err := <-lh.device.ErrChan:
			return err
		case <-ctx.Done():
			return ctx.Err()
		case <-lh.stopChan:
			return nil
		}
	}
}

func (lh *LinkMux) Stop() error {
	lh.stopChan <- struct{}{}
	return nil
}

func (lh *LinkMux) routePacket(ctx context.Context, rawPacket []byte) {
	switch rawPacket[0] {
	case BroadcastPacketType:
		packet, payload, err := ParseBroadcastPacket(rawPacket)
		if err != nil {
			lh.logger.Error("Error parsing packet", "err", err)
			return
		}

		lh.logger.Info("Parsed broadcast packet", "channel", packet.BroadcastChannel, "control", packet.Control, "payloadSize", len(payload), "senderId", packet.SenderId, "familyId", packet.MessageFamily)
	case LinkedPacketType:
		packet, payload, err := ParseLinkedPacket(rawPacket)
		if err != nil {
			lh.logger.Error("Error parsing packet", "err", err)
			return
		}

		lh.logger.Info("Parsed linked packet", "linkId", packet.LinkId, "control", packet.Control, "payloadSize", len(payload))
	default:
		lh.logger.Debug("Skipping packet, neither linked or broadcast packet", "packet", Payload(rawPacket).String())
	}
}
