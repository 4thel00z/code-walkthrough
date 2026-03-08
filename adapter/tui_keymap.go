package adapter

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Next       key.Binding
	Prev       key.Binding
	TOC        key.Binding
	ToggleDiag key.Binding
	Search     key.Binding
	Enter      key.Binding
	Bookmark   key.Binding
	Bookmarks  key.Binding
	Export     key.Binding
	Escape     key.Binding
	Help       key.Binding
	Quit       key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Next:       key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j/↓", "next step")),
		Prev:       key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k/↑", "prev step")),
		TOC:        key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "table of contents")),
		ToggleDiag: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "toggle diagram")),
		Search:     key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
		Enter:      key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		Bookmark:   key.NewBinding(key.WithKeys("b"), key.WithHelp("b", "toggle bookmark")),
		Bookmarks:  key.NewBinding(key.WithKeys("B"), key.WithHelp("B", "list bookmarks")),
		Export:     key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "export")),
		Escape:     key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
		Help:       key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Quit:       key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}
