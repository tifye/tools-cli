package device

import (
	"context"
	"net"
	"os"
	"os/signal"
	"time"

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

			linkHost := linking.NewLinkHost(device, tCli.Log)
			go func() {
				timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				if err := linkHost.Start(timeoutCtx); err != nil {
					tCli.Log.Fatal("Link host error", "err", err)
				}
			}()

			//device.Write([]byte{0x02, 0xFD, 0x0D, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x63, 0x12, 0x8E, 0x8C, 0x0D, 0x2B, 0x60, 0x03})

			<-ctx.Done()

			linkHost.Stop()
		},
	}

	// todo: Add defaults to config
	openCmd.Flags().StringVarP(&opts.address, "address", "a", "127.0.0.1:4250", "Network address of the device")
	openCmd.Flags().StringVarP(&opts.network, "network", "n", "tcp", "Network type of the device")

	return openCmd
}
