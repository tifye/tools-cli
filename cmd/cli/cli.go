package cli

import (
	"github.com/Tifufu/tools-cli/pkg"
	"github.com/charmbracelet/log"
)

type ToolsCli struct {
	User *pkg.UserProfile
	Log  *log.Logger
}
