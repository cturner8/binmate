package tui

import (
	tea "charm.land/bubbletea/v2"
	"cturner8/binmate/internal/core/config"
	"cturner8/binmate/internal/database/repository"
)

func InitProgram(dbService *repository.Service, cfg *config.Config) *tea.Program {
	return tea.NewProgram(initialModel(dbService, cfg))
}
