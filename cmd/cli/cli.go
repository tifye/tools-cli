package cli

import (
	"net/http"

	"github.com/Tifufu/tools-cli/pkg"
	"github.com/Tifufu/tools-cli/pkg/security"
	"github.com/charmbracelet/log"
)

type ToolsCli struct {
	User             *security.UserProfile
	Log              *log.Logger
	WinMowerRegistry *pkg.WinMowerRegistry
	BundleRegistry   *pkg.BundleRegistry
	Client           *http.Client
}
