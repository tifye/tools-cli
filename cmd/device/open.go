package device

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
	"os"
	"os/signal"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/internal/automower"
	"github.com/Tifufu/tools-cli/internal/protocol"
	"github.com/spf13/cobra"
)

type openOptions struct {
	address string
	network string
}

func newOpenCommand(tCli *cli.ToolsCli) *cobra.Command {
	opts := &openOptions{}
	openCmd := &cobra.Command{
		Use:   "open",
		Short: "Open a device and print its data",
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := net.Dial(opts.network, opts.address)
			if err != nil {
				tCli.Log.Fatal("Error opening device", "err", err)
			}
			defer conn.Close()

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			device := automower.NewDevice(conn, ctx)
			defer device.Close()

			device.Write([]byte{0x02, 0xFD, 0x0D, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x63, 0x12, 0x8E, 0x8C, 0x0D, 0x2B, 0x60, 0x03})
			for {
				select {
				case rawPacket := <-device.PacketChan:
					if len(rawPacket) <= 0 {
						continue
					}

					if rawPacket[0] != protocol.LinkedPacketType {
						tCli.Log.Debug("Skipping packet, not a linked packet", "packet", protocol.Payload(rawPacket).String())
						continue
					}

					packet, payload, err := protocol.ParseLinkedPacket(rawPacket)
					if err != nil {
						tCli.Log.Error("Error parsing packet", "err", err)
						continue
					}

					tCli.Log.Info("Parse linked packet", "linkId", packet.LinkId, "control", packet.Control, "payloadSize", len(payload))

					if packet.Control == 0x00 {
						buff := bytes.NewBuffer(payload)
						messageId, _ := buff.ReadByte()
						if messageId == 0x13 {
							traceBackLinkId := binary.LittleEndian.Uint32(buff.Next(4))
							nodeType := binary.LittleEndian.Uint32(buff.Next(4))
							nodeName := buff.Bytes()

							tCli.Log.Info("ListAllNodesResponse received", "traceBackLinkId", traceBackLinkId, "nodeType", nodeType, "nodeName", string(nodeName))
						}
					}
				case err := <-device.ErrChan:
					tCli.Log.Error("Error on device", "err", err)
					return
				case <-ctx.Done():
					return
				}
			}
		},
	}

	// todo: Add defaults to config
	openCmd.Flags().StringVarP(&opts.address, "address", "a", "127.0.0.1:4250", "Network address of the device")
	openCmd.Flags().StringVarP(&opts.network, "network", "n", "tcp", "Network type of the device")

	return openCmd
}
