package adapter

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
)

type KeyMap struct {
	Next       key.Binding
	Prev       key.Binding
	TOC        key.Binding
	ToggleDiag key.Binding
	ToggleCode key.Binding
	Search     key.Binding
	Enter      key.Binding
	Bookmark   key.Binding
	Bookmarks  key.Binding
	Export     key.Binding
	Escape     key.Binding
	Help       key.Binding
	Quit       key.Binding
}

func viewportKeyMap() viewport.KeyMap {
	return viewport.KeyMap{
		Up:           key.NewBinding(key.WithKeys("up")),
		Down:         key.NewBinding(key.WithKeys("down")),
		HalfPageUp:   key.NewBinding(key.WithKeys("ctrl+u")),
		HalfPageDown: key.NewBinding(key.WithKeys("ctrl+d")),
	}
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Next:       key.NewBinding(key.WithKeys("j", "right"), key.WithHelp("j/→", "next step")),
		Prev:       key.NewBinding(key.WithKeys("k", "left"), key.WithHelp("k/←", "prev step")),
		TOC:        key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "table of contents")),
		ToggleDiag: key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "toggle diagram")),
		ToggleCode: key.NewBinding(key.WithKeys("c"), key.WithHelp("c", "toggle code")),
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
