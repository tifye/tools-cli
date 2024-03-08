package sites

import "github.com/spf13/cobra"

type addOptions struct {
	name string
	url  string
}

func newAddCommand() *cobra.Command {
	opts := &addOptions{}
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a site",
		Run: func(cmd *cobra.Command, args []string) {
			// Add site
		},
	}

	cmd.Flags().StringVarP(&opts.name, "site", "s", "", "The site name")
	cmd.Flags().StringVarP(&opts.url, "url", "u", "", "The site URL")
	cmd.MarkFlagRequired("site")
	cmd.MarkFlagRequired("url")

	return cmd
}
