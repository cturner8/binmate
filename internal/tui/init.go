package tui

import (
	"cturner8/binmate/internal/core/config"
	"cturner8/binmate/internal/database/repository"
	"cturner8/binmate/internal/providers/github"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		tea.SetWindowTitle("binmate"),
		loadBinaries(m.dbService),
	}

	// Resolve the GitHub token once at startup unless the config requests per-call resolution.
	askpassMode := ""
	if providers := m.config.Global.Providers; providers != nil {
		if gh, ok := providers["github"]; ok {
			askpassMode = gh.AskpassMode
		}
	}

	if askpassMode != "always" {
		cmds = append(cmds, resolveGithubTokenCmd())
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

// clientForConfig resolves the askpass mode from the config.
func askpassModeFromConfig(cfg *config.Config) string {
	if cfg == nil {
		return ""
	}
	if providers := cfg.Global.Providers; providers != nil {
		if gh, ok := providers["github"]; ok {
			return gh.AskpassMode
		}
	}
	return ""
}
