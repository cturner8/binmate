package tui

import (
	"cturner8/binmate/internal/core/config"
	"cturner8/binmate/internal/core/format"
)

// getDateFormat returns the configured date format or a sensible default.
func getDateFormat(cfg *config.Config) string {
	dateFormat := format.GetDefaultDateFormat()
	if cfg != nil && cfg.TUI.DateFormat != "" {
		dateFormat = cfg.TUI.DateFormat
	}
	return dateFormat
}
