package tui

import (
	"cturner8/binmate/internal/database/repository"
	"cturner8/binmate/internal/providers/github"

	"log"

	tea "charm.land/bubbletea/v2"
)

func (m model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		loadBinaries(m.dbService),
	}

	// Resolve the GitHub token once at startup unless the config requests per-call resolution.
	if m.config.TUI.AskpassMode == "startup" {
		log.Println("Resolving authentication tokens...")
		cmds = append(cmds, resolveGithubTokenCmd())
	} else {
		log.Printf("Skipping authentication token resolution at startup; will resolve on demand. Mode: %s", m.config.TUI.AskpassMode)
	}

	return tea.Batch(cmds...)
}

// loadBinaries returns a command to load binaries from the database
func loadBinaries(dbService *repository.Service) tea.Cmd {
	return func() tea.Msg {
		binaries, err := getBinariesWithMetadata(dbService)
		return binariesLoadedMsg{binaries: binaries, err: err}
	}
}

// resolveGithubTokenCmd returns a Bubble Tea command that resolves the GitHub
// token in the background and returns a githubTokenResolvedMsg.
func resolveGithubTokenCmd() tea.Cmd {
	return func() tea.Msg {
		token, err := github.ResolveToken()
		return githubTokenResolvedMsg{token: token, err: err}
	}
}
