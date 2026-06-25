package tui

import (
	"cturner8/binmate/internal/core/config"
	"cturner8/binmate/internal/database/repository"
	"cturner8/binmate/internal/providers/github"

	tea "charm.land/bubbletea/v2"
)

func (m model) Init() tea.Cmd {
	cmds := []tea.Cmd{
		loadBinaries(m.dbService),
	}

	// Resolve the GitHub token once at startup unless the config requests per-call resolution.
	if askpassModeFromConfig(m.config) == "startup" {
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

// askpassModeFromConfig returns the configured askpass mode for the GitHub
// provider, or an empty string if not set.
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
