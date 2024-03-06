package profile

import (
	"fmt"

	"github.com/Tifufu/tools-cli/cmd/cli"
	"github.com/spf13/cobra"
)

func newListCommand(toolsCli *cli.ToolsCli) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List various user profile information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ID %s\n", toolsCli.User.ID)
			fmt.Printf("Email %s\n", toolsCli.User.Profile.Email)
			fmt.Printf("Fullname %s\n", toolsCli.User.Profile.Fullname)
			fmt.Printf("Human %t\n", toolsCli.User.Profile.Human)

			fmt.Printf("GlobalAdmin %t\n", toolsCli.User.GlobalAdmin)
			fmt.Printf("Developer %t\n", toolsCli.User.Developer)
			fmt.Printf("SystemTest %t\n", toolsCli.User.SystemTest)
			fmt.Printf("BetaTest %t\n", toolsCli.User.BetaTest)
			fmt.Printf("InternalTest %t\n", toolsCli.User.InternalTest)
			for role, accessLevel := range toolsCli.User.Roles {
				fmt.Printf("%s %s\n", role, accessLevel)
			}
		},
	}

	return listCmd
}
