package port

import "github.com/4thel00z/code-walkthrough/domain"

type WalkthroughRepository interface {
	Read(path string) ([]byte, error)
	Write(path string, data []byte) error
}

type DiagramRenderer interface {
	Render(diagram domain.Diagram, width int) (string, error)
}

type Presenter interface {
	RenderStep(step domain.Step, sectionTitle string, stepIndex, totalSteps, sectionIndex, totalSections int)
	RenderDiagram(ascii string)
	RenderTOC(sections []domain.Section)
	RenderSearchResults(results []domain.SearchResult)
	RenderBookmarks(bookmarks []domain.Bookmark)
}

type BookmarkStore interface {
	Save(bookmarks []domain.Bookmark) error
	Load() ([]domain.Bookmark, error)
}

type SchemaValidator interface {
	Validate(data []byte) error
}

type SkillInstaller interface {
	Install(dir string, skill []byte, schema []byte) error
}
