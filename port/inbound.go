package port

import (
	"io"

	"github.com/tahrioui/code-walkthrough/domain"
)

type ExportFormat string

const (
	ExportMarkdown ExportFormat = "markdown"
	ExportHTML     ExportFormat = "html"
)

type WalkthroughLoader interface {
	Load(source string) (domain.Walkthrough, error)
}

type NavigationPort interface {
	Current() (domain.Step, error)
	Next() (domain.Step, error)
	Prev() (domain.Step, error)
	JumpTo(id domain.StepID) (domain.Step, error)
	JumpToSection(id domain.SectionID) (domain.Step, error)
	CurrentSection() domain.Section
	CurrentIndex() int
	TotalSteps() int
}

type SearchPort interface {
	Search(query string) []domain.SearchResult
}

type BookmarkPort interface {
	Add(id domain.StepID) error
	Remove(id domain.StepID) error
	List() []domain.Bookmark
	IsBookmarked(id domain.StepID) bool
}

type ExportPort interface {
	Export(format ExportFormat, w io.Writer) error
}
