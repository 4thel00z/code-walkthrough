package adapter

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tahrioui/code-walkthrough/application"
	"github.com/tahrioui/code-walkthrough/domain"
)

type viewMode int

const (
	viewStep viewMode = iota
	viewTOC
	viewSearch
	viewBookmarks
	viewHelp
)

type Model struct {
	walkthrough   domain.Walkthrough
	navigate      *application.NavigateUseCase
	search        *application.SearchUseCase
	bookmarks     *application.BookmarkUseCase
	renderer      *MermaidRenderer
	keys          KeyMap
	styles        Styles
	mode          viewMode
	showDiagram   bool
	searchInput   textinput.Model
	searchResults []domain.SearchResult
	tocCursor     int
	listCursor    int
	width         int
	height        int
}

func NewModel(
	w domain.Walkthrough,
	nav *application.NavigateUseCase,
	srch *application.SearchUseCase,
	bm *application.BookmarkUseCase,
	renderer *MermaidRenderer,
) Model {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.CharLimit = 100

	return Model{
		walkthrough: w,
		navigate:    nav,
		search:      srch,
		bookmarks:   bm,
		renderer:    renderer,
		keys:        DefaultKeyMap(),
		styles:      DefaultStyles(),
		mode:        viewStep,
		showDiagram: true,
		searchInput: ti,
		width:       80,
		height:      24,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case viewSearch:
			return m.updateSearch(msg)
		case viewTOC:
			return m.updateTOC(msg)
		case viewBookmarks:
			return m.updateBookmarks(msg)
		case viewHelp:
			m.mode = viewStep
			return m, nil
		default:
			return m.updateStep(msg)
		}
	}
	return m, nil
}

func (m Model) updateStep(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case matchKey(msg, m.keys.Quit):
		return m, tea.Quit
	case matchKey(msg, m.keys.Next):
		m.navigate.StepForward()
	case matchKey(msg, m.keys.Prev):
		m.navigate.StepBackward()
	case matchKey(msg, m.keys.ToggleDiag):
		m.showDiagram = !m.showDiagram
	case matchKey(msg, m.keys.TOC):
		m.mode = viewTOC
		m.tocCursor = 0
	case matchKey(msg, m.keys.Search):
		m.mode = viewSearch
		m.searchInput.Focus()
		m.searchResults = nil
	case matchKey(msg, m.keys.Bookmark):
		step, err := m.navigate.Current()
		if err == nil {
			if m.bookmarks.IsBookmarked(step.ID) {
				m.bookmarks.Remove(step.ID)
			} else {
				m.bookmarks.Add(step.ID)
			}
		}
	case matchKey(msg, m.keys.Bookmarks):
		m.mode = viewBookmarks
		m.listCursor = 0
	case matchKey(msg, m.keys.Help):
		m.mode = viewHelp
	}
	return m, nil
}

func (m Model) updateTOC(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case matchKey(msg, m.keys.Escape), matchKey(msg, m.keys.TOC):
		m.mode = viewStep
	case matchKey(msg, m.keys.Quit):
		return m, tea.Quit
	case matchKey(msg, m.keys.Next):
		if m.tocCursor < len(m.walkthrough.Sections)-1 {
			m.tocCursor++
		}
	case matchKey(msg, m.keys.Prev):
		if m.tocCursor > 0 {
			m.tocCursor--
		}
	case matchKey(msg, m.keys.Enter):
		if m.tocCursor < len(m.walkthrough.Sections) {
			sec := m.walkthrough.Sections[m.tocCursor]
			m.navigate.JumpToSection(sec.ID)
			m.mode = viewStep
		}
	}
	return m, nil
}

func (m Model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case matchKey(msg, m.keys.Escape):
		m.mode = viewStep
		m.searchInput.Blur()
		return m, nil
	case matchKey(msg, m.keys.Enter):
		if len(m.searchResults) > 0 && m.listCursor < len(m.searchResults) {
			m.navigate.JumpTo(m.searchResults[m.listCursor].StepID)
			m.mode = viewStep
			m.searchInput.Blur()
			return m, nil
		}
		// Perform search
		m.searchResults = m.search.Search(m.searchInput.Value())
		m.listCursor = 0
		return m, nil
	}

	// Navigate results if we have them
	if len(m.searchResults) > 0 {
		switch {
		case matchKey(msg, m.keys.Next):
			if m.listCursor < len(m.searchResults)-1 {
				m.listCursor++
			}
			return m, nil
		case matchKey(msg, m.keys.Prev):
			if m.listCursor > 0 {
				m.listCursor--
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	return m, cmd
}

func (m Model) updateBookmarks(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	bms := m.bookmarks.List()
	switch {
	case matchKey(msg, m.keys.Escape), matchKey(msg, m.keys.Bookmarks):
		m.mode = viewStep
	case matchKey(msg, m.keys.Quit):
		return m, tea.Quit
	case matchKey(msg, m.keys.Next):
		if m.listCursor < len(bms)-1 {
			m.listCursor++
		}
	case matchKey(msg, m.keys.Prev):
		if m.listCursor > 0 {
			m.listCursor--
		}
	case matchKey(msg, m.keys.Enter):
		if m.listCursor < len(bms) {
			m.navigate.JumpTo(bms[m.listCursor].StepID)
			m.mode = viewStep
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.mode {
	case viewTOC:
		return m.viewTOC()
	case viewSearch:
		return m.viewSearch()
	case viewBookmarks:
		return m.viewBookmarks()
	case viewHelp:
		return m.viewHelp()
	default:
		return m.viewStep()
	}
}

func (m Model) viewStep() string {
	step, err := m.navigate.Current()
	if err != nil {
		return m.styles.Muted.Render("No steps to display")
	}

	sec := m.navigate.CurrentSection()
	sectionIdx := 0
	for i, s := range m.walkthrough.Sections {
		if s.ID == sec.ID {
			sectionIdx = i
			break
		}
	}

	var b strings.Builder

	// Header
	header := fmt.Sprintf(" [Section %d/%d] %s    [Step %d/%d]",
		sectionIdx+1, len(m.walkthrough.Sections), sec.Title,
		m.navigate.CurrentIndex()+1, m.navigate.TotalSteps())
	bookmarkIndicator := ""
	if m.bookmarks.IsBookmarked(step.ID) {
		bookmarkIndicator = " *"
	}
	b.WriteString(m.styles.Header.Render(header+bookmarkIndicator) + "\n")

	// Step title
	b.WriteString(m.styles.StepTitle.Render(step.Title) + "\n\n")

	// Explanation
	b.WriteString(m.styles.Explanation.Render(step.Explanation) + "\n")

	// Code snippet
	if step.CodeSnippet != nil {
		label := fmt.Sprintf("─ %s:%d-%d ", step.CodeSnippet.FilePath, step.CodeSnippet.StartLine, step.CodeSnippet.EndLine)
		code := m.styles.CodeBlock.Render(step.CodeSnippet.Source)
		b.WriteString("\n" + m.styles.CodeBorder.Render(label+"\n"+code) + "\n")
	}

	// Diagram
	if m.showDiagram && step.Diagram != nil {
		ascii, err := m.renderer.Render(*step.Diagram, m.width-4)
		if err == nil && ascii != "" {
			label := fmt.Sprintf("─ %s ", step.Diagram.Type)
			b.WriteString("\n" + m.styles.DiagBorder.Render(label+"\n"+m.styles.DiagBlock.Render(ascii)) + "\n")
		}
	}

	// Status bar
	b.WriteString("\n")
	status := " j/k:navigate  g:toc  /:search  d:diagram  b:bookmark  e:export  ?:help "
	b.WriteString(m.styles.StatusBar.Width(m.width).Render(status))

	return b.String()
}

func (m Model) viewTOC() string {
	var b strings.Builder
	b.WriteString(m.styles.Header.Render(" Table of Contents") + "\n\n")

	for i, sec := range m.walkthrough.Sections {
		cursor := "  "
		style := m.styles.Explanation
		if i == m.tocCursor {
			cursor = "▸ "
			style = m.styles.Highlight
		}
		b.WriteString(style.Render(fmt.Sprintf("%s%s (%d steps)", cursor, sec.Title, len(sec.Steps))) + "\n")
	}

	b.WriteString("\n")
	b.WriteString(m.styles.StatusBar.Width(m.width).Render(" j/k:navigate  enter:select  esc:back "))
	return b.String()
}

func (m Model) viewSearch() string {
	var b strings.Builder
	b.WriteString(m.styles.Header.Render(" Search") + "\n\n")
	b.WriteString(" " + m.searchInput.View() + "\n\n")

	if len(m.searchResults) > 0 {
		for i, r := range m.searchResults {
			cursor := "  "
			style := m.styles.Explanation
			if i == m.listCursor {
				cursor = "▸ "
				style = m.styles.Highlight
			}
			b.WriteString(style.Render(fmt.Sprintf("%s%s", cursor, r.StepTitle)) + "\n")
		}
	} else if m.searchInput.Value() != "" {
		b.WriteString(m.styles.Muted.Render("  No results. Press enter to search.") + "\n")
	}

	b.WriteString("\n")
	b.WriteString(m.styles.StatusBar.Width(m.width).Render(" enter:search/select  j/k:navigate results  esc:back "))
	return b.String()
}

func (m Model) viewBookmarks() string {
	var b strings.Builder
	b.WriteString(m.styles.Header.Render(" Bookmarks") + "\n\n")

	bms := m.bookmarks.List()
	if len(bms) == 0 {
		b.WriteString(m.styles.Muted.Render("  No bookmarks yet. Press 'b' on a step to bookmark it.") + "\n")
	} else {
		for i, bm := range bms {
			cursor := "  "
			style := m.styles.Explanation
			if i == m.listCursor {
				cursor = "▸ "
				style = m.styles.Highlight
			}
			b.WriteString(style.Render(fmt.Sprintf("%s%s", cursor, string(bm.StepID))) + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(m.styles.StatusBar.Width(m.width).Render(" j/k:navigate  enter:select  esc:back "))
	return b.String()
}

func (m Model) viewHelp() string {
	var b strings.Builder
	b.WriteString(m.styles.Header.Render(" Help") + "\n\n")

	km := m.keys
	bindings := []struct{ key, desc string }{
		{km.Next.Help().Key, km.Next.Help().Desc},
		{km.Prev.Help().Key, km.Prev.Help().Desc},
		{km.TOC.Help().Key, km.TOC.Help().Desc},
		{km.ToggleDiag.Help().Key, km.ToggleDiag.Help().Desc},
		{km.Search.Help().Key, km.Search.Help().Desc},
		{km.Bookmark.Help().Key, km.Bookmark.Help().Desc},
		{km.Bookmarks.Help().Key, km.Bookmarks.Help().Desc},
		{km.Export.Help().Key, km.Export.Help().Desc},
		{km.Help.Help().Key, km.Help.Help().Desc},
		{km.Quit.Help().Key, km.Quit.Help().Desc},
	}

	for _, bind := range bindings {
		b.WriteString(fmt.Sprintf("  %s  %s\n",
			m.styles.Highlight.Render(fmt.Sprintf("%-8s", bind.key)),
			m.styles.Explanation.Render(bind.desc)))
	}

	b.WriteString("\n")
	b.WriteString(m.styles.StatusBar.Width(m.width).Render(" Press any key to go back "))
	return b.String()
}

func matchKey(msg tea.KeyMsg, binding key.Binding) bool {
	for _, k := range binding.Keys() {
		if msg.String() == k {
			return true
		}
	}
	return false
}

func RunTUI(m Model) error {
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
