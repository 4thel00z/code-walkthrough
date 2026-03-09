package adapter

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/4thel00z/code-walkthrough/application"
	"github.com/4thel00z/code-walkthrough/domain"
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
	viewport      viewport.Model
	mode          viewMode
	showDiagram   bool
	showCode      bool
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

	vp := viewport.New(80, 20)
	vp.KeyMap = viewportKeyMap()

	return Model{
		walkthrough: w,
		navigate:    nav,
		search:      srch,
		bookmarks:   bm,
		renderer:    renderer,
		keys:        DefaultKeyMap(),
		styles:      DefaultStyles(),
		viewport:    vp,
		mode:        viewStep,
		showDiagram: true,
		showCode:    true,
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
		m.viewport.Width = msg.Width
		// Viewport height is computed dynamically in View() based on header/status,
		// but we set a reasonable default here for the initial resize.
		m.viewport.Height = max(msg.Height-4, 1)
		return m, nil

	case tea.MouseMsg:
		m = m.syncViewport()
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd

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
			m.viewport.GotoTop()
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
		m.viewport.GotoTop()
	case matchKey(msg, m.keys.Prev):
		m.navigate.StepBackward()
		m.viewport.GotoTop()
	case matchKey(msg, m.keys.ToggleDiag):
		m.showDiagram = !m.showDiagram
	case matchKey(msg, m.keys.ToggleCode):
		m.showCode = !m.showCode
	case matchKey(msg, m.keys.TOC):
		m.mode = viewTOC
		m.tocCursor = 0
		m.viewport.GotoTop()
	case matchKey(msg, m.keys.Search):
		m.mode = viewSearch
		m.searchInput.Focus()
		m.searchResults = nil
		m.viewport.GotoTop()
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
		m.viewport.GotoTop()
	case matchKey(msg, m.keys.Help):
		m.mode = viewHelp
		m.viewport.GotoTop()
	default:
		m = m.syncViewport()
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) updateTOC(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case matchKey(msg, m.keys.Escape), matchKey(msg, m.keys.TOC):
		m.mode = viewStep
		m.viewport.GotoTop()
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
			m.viewport.GotoTop()
		}
	default:
		m = m.syncViewport()
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case matchKey(msg, m.keys.Escape):
		m.mode = viewStep
		m.searchInput.Blur()
		m.viewport.GotoTop()
		return m, nil
	case matchKey(msg, m.keys.Enter):
		if len(m.searchResults) > 0 && m.listCursor < len(m.searchResults) {
			m.navigate.JumpTo(m.searchResults[m.listCursor].StepID)
			m.mode = viewStep
			m.searchInput.Blur()
			m.viewport.GotoTop()
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
		m.viewport.GotoTop()
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
			m.viewport.GotoTop()
		}
	default:
		m = m.syncViewport()
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

// syncViewport sets the viewport's content and dimensions based on the current mode.
// This must be called before forwarding messages to the viewport, so it knows
// the content height and can calculate valid scroll bounds.
func (m Model) syncViewport() Model {
	var header, body, status string
	switch m.mode {
	case viewTOC:
		header, body, status = m.viewTOCParts()
	case viewSearch:
		header, body, status = m.viewSearchParts()
	case viewBookmarks:
		header, body, status = m.viewBookmarksParts()
	case viewHelp:
		header, body, status = m.viewHelpParts()
	default:
		header, body, status = m.viewStepParts()
	}
	headerHeight := lipgloss.Height(header)
	statusHeight := lipgloss.Height(status)
	vpHeight := m.height - headerHeight - statusHeight
	if vpHeight < 1 {
		vpHeight = 1
	}
	m.viewport.Height = vpHeight
	m.viewport.Width = m.width
	m.viewport.SetContent(body)
	return m
}

// View composes the 3-zone layout: fixed header, scrollable body, fixed status bar.
func (m Model) View() string {
	var header, status string

	switch m.mode {
	case viewTOC:
		header, _, status = m.viewTOCParts()
	case viewSearch:
		header, _, status = m.viewSearchParts()
	case viewBookmarks:
		header, _, status = m.viewBookmarksParts()
	case viewHelp:
		header, _, status = m.viewHelpParts()
	default:
		header, _, status = m.viewStepParts()
	}

	m = m.syncViewport()

	return header + "\n" + m.viewport.View() + "\n" + status
}

func (m Model) viewStepParts() (header, body, status string) {
	step, err := m.navigate.Current()
	if err != nil {
		return "", m.styles.Muted.Render("No steps to display"), ""
	}

	sec := m.navigate.CurrentSection()
	sectionIdx := 0
	for i, s := range m.walkthrough.Sections {
		if s.ID == sec.ID {
			sectionIdx = i
			break
		}
	}

	// Header
	h := fmt.Sprintf(" [Section %d/%d] %s    [Step %d/%d]",
		sectionIdx+1, len(m.walkthrough.Sections), sec.Title,
		m.navigate.CurrentIndex()+1, m.navigate.TotalSteps())
	bookmarkIndicator := ""
	if m.bookmarks.IsBookmarked(step.ID) {
		bookmarkIndicator = " *"
	}
	header = m.styles.Header.Render(h + bookmarkIndicator)

	// Body
	var b strings.Builder

	b.WriteString(m.styles.StepTitle.Width(m.width).Render(step.Title) + "\n\n")
	b.WriteString(m.styles.Explanation.Width(m.width).Render(step.Explanation) + "\n")

	if step.CodeSnippet != nil {
		label := fmt.Sprintf("─ %s:%d-%d ", step.CodeSnippet.FilePath, step.CodeSnippet.StartLine, step.CodeSnippet.EndLine)
		if m.showCode {
			code := m.styles.CodeBlock.Render(step.CodeSnippet.Source)
			b.WriteString("\n" + m.styles.CodeBorder.Width(m.width).Render(label+"\n"+code) + "\n")
		} else {
			b.WriteString("\n" + m.styles.CodeBorder.Width(m.width).Render(label+"[hidden]") + "\n")
		}
	}

	if step.Diagram != nil {
		label := fmt.Sprintf("─ %s ", step.Diagram.Type)
		if m.showDiagram {
			ascii, err := m.renderer.Render(*step.Diagram, m.width-4)
			if err == nil && ascii != "" {
				b.WriteString("\n" + m.styles.DiagBorder.Width(m.width).Render(label+"\n"+m.styles.DiagBlock.Render(ascii)) + "\n")
			}
		} else {
			b.WriteString("\n" + m.styles.DiagBorder.Width(m.width).Render(label+"[hidden]") + "\n")
		}
	}
	body = b.String()

	// Status bar
	scrollPct := fmt.Sprintf(" %3.f%%", m.viewport.ScrollPercent()*100)
	statusText := " j/k/←→:step  ↑↓:scroll  g:toc  /:search  c:code  d:diagram  b:bookmark  ?:help" + scrollPct + " "
	status = m.styles.StatusBar.Width(m.width).Render(statusText)

	return header, body, status
}

func (m Model) viewTOCParts() (header, body, status string) {
	header = m.styles.Header.Render(" Table of Contents")

	var b strings.Builder
	for i, sec := range m.walkthrough.Sections {
		cursor := "  "
		style := m.styles.Explanation
		if i == m.tocCursor {
			cursor = "▸ "
			style = m.styles.Highlight
		}
		b.WriteString(style.Render(fmt.Sprintf("%s%s (%d steps)", cursor, sec.Title, len(sec.Steps))) + "\n")
	}
	body = b.String()

	status = m.styles.StatusBar.Width(m.width).Render(" j/k:navigate  enter:select  esc:back ")
	return header, body, status
}

func (m Model) viewSearchParts() (header, body, status string) {
	header = m.styles.Header.Render(" Search") + "\n\n " + m.searchInput.View()

	var b strings.Builder
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
	body = b.String()

	status = m.styles.StatusBar.Width(m.width).Render(" enter:search/select  j/k:navigate results  esc:back ")
	return header, body, status
}

func (m Model) viewBookmarksParts() (header, body, status string) {
	header = m.styles.Header.Render(" Bookmarks")

	var b strings.Builder
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
	body = b.String()

	status = m.styles.StatusBar.Width(m.width).Render(" j/k:navigate  enter:select  esc:back ")
	return header, body, status
}

func (m Model) viewHelpParts() (header, body, status string) {
	header = m.styles.Header.Render(" Help")

	var b strings.Builder
	km := m.keys
	bindings := []struct{ key, desc string }{
		{km.Next.Help().Key, km.Next.Help().Desc},
		{km.Prev.Help().Key, km.Prev.Help().Desc},
		{"↑↓", "scroll content"},
		{km.TOC.Help().Key, km.TOC.Help().Desc},
		{km.ToggleCode.Help().Key, km.ToggleCode.Help().Desc},
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
	body = b.String()

	status = m.styles.StatusBar.Width(m.width).Render(" Press any key to go back ")
	return header, body, status
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
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
