package tui

import tea "charm.land/bubbletea/v2"

func (m model) View() tea.View {
	var content string
	switch m.currentView {
	case viewBinariesList:
		content = m.renderBinariesList()
	case viewVersions:
		content = m.renderVersions()
	case viewAddBinaryURL:
		content = m.renderAddBinaryURL()
	case viewAddBinaryForm:
		content = m.renderAddBinaryForm()
	case viewInstallBinary:
		content = m.renderInstallBinary()
	case viewImportBinary:
		content = m.renderImportBinary()
	case viewConfiguration:
		content = m.renderConfiguration()
	case viewHelp:
		content = m.renderHelp()
	case viewReleaseNotes:
		content = m.renderReleaseNotes()
	case viewAvailableVersions:
		content = m.renderAvailableVersions()
	case viewRepositoryInfo:
		content = m.renderRepositoryInfo()
	default:
		content = "Unknown view"
	}
	v := tea.NewView(content)
	v.WindowTitle = "binmate"
	return v
}
