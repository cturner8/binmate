package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

// renderBinariesList renders the binaries list view
func (m model) renderBinariesList() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("📦 Binmate - Binary Manager"))
	b.WriteString("\n\n")

	// Tabs
	b.WriteString(m.renderTabs())

	// Show loading state
	if m.loading {
		b.WriteString(loadingStyle.Render("Loading binaries..."))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render(getHelpText(m.currentView)))
		return b.String()
	}

	// Show empty state
	if len(m.binaries) == 0 {
		b.WriteString(emptyStateStyle.Render("No binaries configured"))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Press 'a' to add a binary"))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render(getHelpText(m.currentView)))
		return b.String()
	}

	// Show error if any
	if m.errorMessage != "" {
		b.WriteString(errorStyle.Render("Error: " + m.errorMessage))
		b.WriteString("\n\n")
	}

	// Show success message if any
	if m.successMessage != "" {
		b.WriteString(successStyle.Render("✓ " + m.successMessage))
		b.WriteString("\n\n")
	}

	// Show remove confirmation if active
	if m.confirmingRemove {
		if m.bulkRemoveCount > 0 {
			// Bulk removal confirmation
			b.WriteString(headerStyle.Render(fmt.Sprintf("Remove %d selected binaries?", m.bulkRemoveCount)))
		} else {
			// Single binary removal confirmation
			b.WriteString(headerStyle.Render(fmt.Sprintf("Remove binary '%s'?", m.removeBinaryID)))
		}
		b.WriteString("\n\n")
		b.WriteString("Press 'y' to remove from database only\n")
		b.WriteString("Press 'Y' (Shift+Y) to also delete files from disk\n")
		b.WriteString("Press 'n' or Esc to cancel\n")
		b.WriteString("\n")
		return b.String()
	}

	// Show search mode UI
	if m.searchMode {
		b.WriteString(headerStyle.Render("🔍 Search Binaries"))
		b.WriteString("\n\n")
		b.WriteString(m.searchTextInput.View())
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Type to search (regex supported) • Enter: apply filter • Esc: cancel"))
		b.WriteString("\n\n")
	}

	// Show filter panel if open
	if m.filterPanelOpen {
		b.WriteString(headerStyle.Render("🔧 Filters"))
		b.WriteString("\n\n")

		// Show current filters
		if len(m.activeFilters) > 0 {
			b.WriteString("Active filters:\n")
			for key, value := range m.activeFilters {
				b.WriteString(fmt.Sprintf("  %s: %s\n", key, value))
			}
		} else {
			b.WriteString("No active filters\n")
		}

		b.WriteString("\n")
		b.WriteString("Filter by:\n")
		b.WriteString("  1: Provider (github)\n")
		b.WriteString("  2: Format (.tar.gz, .zip)\n")
		b.WriteString("  3: Status (installed, not-installed)\n")
		b.WriteString("  c: Clear all filters\n")
		b.WriteString("  Esc: Close filter panel\n")
		b.WriteString("\n")
		return b.String()
	}

	// Determine which binaries list to display
	binariesToShow := getDisplayBinaries(m.binaries, m.activeFilters, m.searchQuery, m.sortMode, m.sortAscending)

	// Display filter/sort indicators
	var indicators []string
	if m.searchQuery != "" && !m.searchMode {
		indicators = append(indicators, fmt.Sprintf("🔍 Search: \"%s\"", m.searchQuery))
	}
	if len(m.activeFilters) > 0 {
		filterStr := ""
		for key, value := range m.activeFilters {
			if filterStr != "" {
				filterStr += ", "
			}
			filterStr += fmt.Sprintf("%s=%s", key, value)
		}
		indicators = append(indicators, fmt.Sprintf("🔧 Filters: %s", filterStr))
	}
	sortDir := "↑"
	if !m.sortAscending {
		sortDir = "↓"
	}
	indicators = append(indicators, fmt.Sprintf("📊 Sort: %s %s", m.sortMode, sortDir))

	// Show bulk mode indicator
	if m.bulkSelectMode {
		indicators = append(indicators, fmt.Sprintf("✓ Bulk Mode: %d selected", len(m.selectedBinaries)))
	}

	if len(indicators) > 0 {
		b.WriteString(fmt.Sprintf("%s (%d results)\n\n", strings.Join(indicators, " • "), len(binariesToShow)))
	}

	// Calculate proportional column widths based on available width
	// Default to 80 if width not set
	availableWidth := m.width
	if availableWidth == 0 {
		availableWidth = defaultTerminalWidth
	}

	// Account for padding: 2 chars per column (4 columns = 8 total)
	totalWidth := availableWidth - columnPadding4

	// Allocate proportional widths: Name 35%, Provider 15%, Version 30%, Count 20%
	nameWidth := int(float64(totalWidth) * 0.35)
	providerWidth := int(float64(totalWidth) * 0.15)
	versionWidth := int(float64(totalWidth) * 0.30)
	countWidth := int(float64(totalWidth) * 0.20)

	headers := []string{
		tableHeaderStyle.Width(nameWidth).Render("Name"),
		tableHeaderStyle.Width(providerWidth).Render("Provider"),
		tableHeaderStyle.Width(versionWidth).Render("Active Version"),
		tableHeaderStyle.Width(countWidth).Render("Installed"),
	}
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, headers...))
	b.WriteString("\n")

	// Separator line
	b.WriteString(strings.Repeat("─", nameWidth+providerWidth+versionWidth+countWidth+columnPadding4))
	b.WriteString("\n")

	// Table rows - use binariesToShow instead of m.binaries
	for i, binary := range binariesToShow {
		style := normalStyle
		if i == m.selectedIndex {
			style = selectedStyle
		}

		// Add selection indicator for bulk mode
		selectionIndicator := ""
		if m.bulkSelectMode {
			if m.selectedBinaries[i] {
				selectionIndicator = "☑ "
			} else {
				selectionIndicator = "☐ "
			}
		}

		name := truncateText(selectionIndicator+binary.Binary.Name, nameWidth)
		provider := truncateText(binary.Binary.Provider, providerWidth)
		version := truncateText(binary.ActiveVersion, versionWidth)
		count := fmt.Sprintf("%d", binary.InstallCount)

		row := []string{
			style.Width(nameWidth).Render(name),
			style.Width(providerWidth).Render(provider),
			style.Width(versionWidth).Render(version),
			style.Width(countWidth).Render(count),
		}

		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, row...))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render(getHelpText(m.currentView)))

	return b.String()
}
