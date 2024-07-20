package device

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/Tifufu/tools-cli/internal/automower"
	"github.com/Tifufu/tools-cli/internal/protocol/linking"
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

			linkMux := linking.NewLinkMux(device, tCli.Log)
			go func() {
				err := linkMux.Start()
				if err != nil && err != linking.ErrLinkMuxShuttingDown {
					tCli.Log.Fatal("Link host error", "err", err)
				}
			}()

			// time.Sleep(2 * time.Second)
			// _, err = linkMux.DefaultLink.SendLinkedRequest(linking.ControlLinkManagerCommand, []byte{0x12, 0x01, 0x00, 0x00, 0x00})
			// if err != nil {
			// 	tCli.Log.Error("err sending link manager command", "err", err)
			// }
			_, err = linkMux.DefaultLink.SendLinkedRequest(linking.ControlLinkManagerCommand, []byte{0x08, 0x01})
			if err != nil {
				tCli.Log.Error("err sending link manager command", "err", err)
			}

			// linkMux.Write([]byte{0x02, 0xFD, 0x0D, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x63, 0x12, 0x9F, 0x57, 0xF7, 0x70, 0x8A, 0x03})

			<-ctx.Done()
			linkMux.Stop()
		},
	}

	// todo: Add defaults to config
	openCmd.Flags().StringVarP(&opts.address, "address", "a", "127.0.0.1:4250", "Network address of the device")
	openCmd.Flags().StringVarP(&opts.network, "network", "n", "tcp", "Network type of the device")

	return openCmd
}
