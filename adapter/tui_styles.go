package adapter

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	Header      lipgloss.Style
	StepTitle   lipgloss.Style
	Explanation lipgloss.Style
	CodeBlock   lipgloss.Style
	CodeBorder  lipgloss.Style
	DiagBlock   lipgloss.Style
	DiagBorder  lipgloss.Style
	StatusBar   lipgloss.Style
	Highlight   lipgloss.Style
	Muted       lipgloss.Style
	Bookmark    lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Header:      lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12")).Padding(0, 1),
		StepTitle:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Padding(1, 1),
		Explanation: lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Padding(0, 1),
		CodeBlock:   lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Padding(0, 2),
		CodeBorder:  lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("8")).Padding(0, 1),
		DiagBlock:   lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Padding(0, 2),
		DiagBorder:  lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("8")).Padding(0, 1),
		StatusBar:   lipgloss.NewStyle().Background(lipgloss.Color("237")).Foreground(lipgloss.Color("7")).Padding(0, 1),
		Highlight:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11")),
		Muted:       lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Bookmark:    lipgloss.NewStyle().Foreground(lipgloss.Color("13")),
	}
}
