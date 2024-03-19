package cli

import (
	"github.com/charmbracelet/log"

	"github.com/charmbracelet/lipgloss"
)

func SubProcessLogStyle(color lipgloss.Color) *log.Styles {
	style := log.DefaultStyles()
	style.Prefix = lipgloss.NewStyle().Foreground(color)
	style.Timestamp = lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderLeftBackground(color).
		BorderLeftForeground(color).
		BorderLeft(true).
		PaddingLeft(1)
	return style
}
