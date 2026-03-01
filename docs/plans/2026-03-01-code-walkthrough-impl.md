# Code Walkthrough Implementation Plan

> **For the implementor:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a Go TUI application that renders interactive, step-by-step code walkthroughs with ASCII-rendered Mermaid diagrams, driven by a JSON input file produced by an AI agent skill.

**Architecture:** Hexagonal architecture (ports & adapters) with DDD domain model. The domain layer has zero external dependencies. Application layer orchestrates use cases through port interfaces. Adapters implement those interfaces for Bubble Tea TUI, filesystem, Mermaid rendering, and export.

**Tech Stack:** Go 1.22+, Charmbracelet Bubble Tea/Lip Gloss/Bubbles, Cobra CLI, standard library testing + testify

---

### Task 1: Project Scaffolding

**Files:**
- Create: `go.mod`
- Create: `cmd/walkthrough/main.go`
- Create: `domain/`, `application/`, `port/`, `adapter/`, `schema/` directories

**Step 1: Initialize Go module**

Run: `go mod init github.com/tahrioui/code-walkthrough`

**Step 2: Install core dependencies**

Run:
```bash
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/charmbracelet/bubbles@latest
go get github.com/spf13/cobra@latest
go get github.com/stretchr/testify@latest
```

**Step 3: Create directory structure and stub main**

```go
// cmd/walkthrough/main.go
package main

import "fmt"

func main() {
	fmt.Println("code-walkthrough")
}
```

**Step 4: Verify it compiles**

Run: `go build ./cmd/walkthrough`
Expected: no errors, produces `walkthrough` binary

**Step 5: Commit**

```bash
git add go.mod go.sum cmd/ domain/ application/ port/ adapter/ schema/
git commit -m "chore: scaffold project structure"
```

---

### Task 2: Domain Model

**Files:**
- Create: `domain/model.go`
- Create: `domain/model_test.go`

**Step 1: Write the failing test**

```go
// domain/model_test.go
package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tahrioui/code-walkthrough/domain"
)

func TestNewWalkthrough(t *testing.T) {
	w := domain.NewWalkthrough("Test Title", "A description", domain.ScopeFlow, "/repo")

	assert.Equal(t, "Test Title", w.Title)
	assert.Equal(t, "A description", w.Description)
	assert.Equal(t, domain.ScopeFlow, w.Scope)
	assert.Equal(t, "/repo", w.Repository)
	assert.Empty(t, w.Sections)
}

func TestWalkthrough_AddSection(t *testing.T) {
	w := domain.NewWalkthrough("Title", "Desc", domain.ScopeOverview, "/repo")
	sec := domain.NewSection("sec-1", "Auth Flow", "How auth works")

	w.AddSection(sec)

	assert.Len(t, w.Sections, 1)
	assert.Equal(t, "Auth Flow", w.Sections[0].Title)
}

func TestSection_AddStep(t *testing.T) {
	sec := domain.NewSection("sec-1", "Auth Flow", "How auth works")
	step := domain.NewStep("step-1", "Router entry", "Request hits the router")

	sec.AddStep(step)

	assert.Len(t, sec.Steps, 1)
	assert.Equal(t, "Router entry", sec.Steps[0].Title)
}

func TestStep_WithCodeSnippet(t *testing.T) {
	step := domain.NewStep("step-1", "Router", "Explanation")
	snippet := domain.CodeSnippet{
		FilePath:  "router.go",
		Language:  "go",
		StartLine: 10,
		EndLine:   20,
		Source:    "func Route() {}",
	}

	step.SetCodeSnippet(snippet)

	assert.NotNil(t, step.CodeSnippet)
	assert.Equal(t, "router.go", step.CodeSnippet.FilePath)
}

func TestStep_WithDiagram(t *testing.T) {
	step := domain.NewStep("step-1", "Router", "Explanation")
	diagram := domain.Diagram{
		Type:    domain.DiagramSequence,
		Mermaid: "sequenceDiagram\n    A->>B: hello",
	}

	step.SetDiagram(diagram)

	assert.NotNil(t, step.Diagram)
	assert.Equal(t, domain.DiagramSequence, step.Diagram.Type)
}

func TestWalkthrough_TotalSteps(t *testing.T) {
	w := domain.NewWalkthrough("Title", "Desc", domain.ScopeFlow, "/repo")

	sec1 := domain.NewSection("sec-1", "Section 1", "")
	sec1.AddStep(domain.NewStep("s1", "Step 1", ""))
	sec1.AddStep(domain.NewStep("s2", "Step 2", ""))

	sec2 := domain.NewSection("sec-2", "Section 2", "")
	sec2.AddStep(domain.NewStep("s3", "Step 3", ""))

	w.AddSection(sec1)
	w.AddSection(sec2)

	assert.Equal(t, 3, w.TotalSteps())
}

func TestWalkthrough_AllSteps(t *testing.T) {
	w := domain.NewWalkthrough("Title", "Desc", domain.ScopeFlow, "/repo")

	sec := domain.NewSection("sec-1", "Section 1", "")
	sec.AddStep(domain.NewStep("s1", "Step 1", "Exp 1"))
	sec.AddStep(domain.NewStep("s2", "Step 2", "Exp 2"))
	w.AddSection(sec)

	steps := w.AllSteps()

	assert.Len(t, steps, 2)
	assert.Equal(t, domain.StepID("s1"), steps[0].ID)
	assert.Equal(t, domain.StepID("s2"), steps[1].ID)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./domain/ -v`
Expected: FAIL — types and functions not defined

**Step 3: Write the domain model**

```go
// domain/model.go
package domain

import "time"

type StepID string
type SectionID string

type Scope string

const (
	ScopeFlow     Scope = "flow"
	ScopeOverview Scope = "overview"
)

type DiagramType string

const (
	DiagramSequence  DiagramType = "sequence"
	DiagramFlowchart DiagramType = "flowchart"
	DiagramClass     DiagramType = "classDiagram"
	DiagramGraph     DiagramType = "graph"
)

type CodeSnippet struct {
	FilePath  string
	Language  string
	StartLine int
	EndLine   int
	Source    string
}

type Diagram struct {
	Type    DiagramType
	Mermaid string
}

type Bookmark struct {
	StepID    StepID
	CreatedAt time.Time
}

type Step struct {
	ID          StepID
	Title       string
	Explanation string
	CodeSnippet *CodeSnippet
	Diagram     *Diagram
}

func NewStep(id string, title, explanation string) Step {
	return Step{
		ID:          StepID(id),
		Title:       title,
		Explanation: explanation,
	}
}

func (s *Step) SetCodeSnippet(snippet CodeSnippet) {
	s.CodeSnippet = &snippet
}

func (s *Step) SetDiagram(diagram Diagram) {
	s.Diagram = &diagram
}

type Section struct {
	ID          SectionID
	Title       string
	Description string
	Steps       []Step
}

func NewSection(id, title, description string) Section {
	return Section{
		ID:          SectionID(id),
		Title:       title,
		Description: description,
	}
}

func (s *Section) AddStep(step Step) {
	s.Steps = append(s.Steps, step)
}

type Walkthrough struct {
	Title       string
	Description string
	Scope       Scope
	Repository  string
	GeneratedAt time.Time
	Sections    []Section
}

func NewWalkthrough(title, description string, scope Scope, repository string) Walkthrough {
	return Walkthrough{
		Title:       title,
		Description: description,
		Scope:       scope,
		Repository:  repository,
		GeneratedAt: time.Now(),
	}
}

func (w *Walkthrough) AddSection(section Section) {
	w.Sections = append(w.Sections, section)
}

func (w *Walkthrough) TotalSteps() int {
	total := 0
	for _, sec := range w.Sections {
		total += len(sec.Steps)
	}
	return total
}

func (w *Walkthrough) AllSteps() []Step {
	var steps []Step
	for _, sec := range w.Sections {
		steps = append(steps, sec.Steps...)
	}
	return steps
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./domain/ -v`
Expected: PASS — all 6 tests green

**Step 5: Commit**

```bash
git add domain/model.go domain/model_test.go
git commit -m "feat(domain): add walkthrough aggregate root and value objects"
```

---

### Task 3: Domain Navigator Service

**Files:**
- Create: `domain/navigator.go`
- Create: `domain/navigator_test.go`

**Step 1: Write the failing test**

```go
// domain/navigator_test.go
package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahrioui/code-walkthrough/domain"
)

func newTestWalkthrough() domain.Walkthrough {
	w := domain.NewWalkthrough("Test", "Desc", domain.ScopeFlow, "/repo")

	sec1 := domain.NewSection("sec-1", "Section 1", "First section")
	sec1.AddStep(domain.NewStep("s1", "Step 1", "Exp 1"))
	sec1.AddStep(domain.NewStep("s2", "Step 2", "Exp 2"))

	sec2 := domain.NewSection("sec-2", "Section 2", "Second section")
	sec2.AddStep(domain.NewStep("s3", "Step 3", "Exp 3"))

	w.AddSection(sec1)
	w.AddSection(sec2)
	return w
}

func TestNavigator_Init(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	step, err := nav.Current()
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s1"), step.ID)
	assert.Equal(t, 0, nav.CurrentIndex())
}

func TestNavigator_Next(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	step, err := nav.Next()
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s2"), step.ID)
}

func TestNavigator_Next_CrossesSection(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	nav.Next() // s2
	step, err := nav.Next() // s3, crosses into sec-2
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s3"), step.ID)
}

func TestNavigator_Next_AtEnd(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	nav.Next() // s2
	nav.Next() // s3
	_, err := nav.Next()
	assert.Error(t, err)
}

func TestNavigator_Prev(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	nav.Next() // s2
	step, err := nav.Prev()
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s1"), step.ID)
}

func TestNavigator_Prev_AtStart(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	_, err := nav.Prev()
	assert.Error(t, err)
}

func TestNavigator_JumpTo(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	step, err := nav.JumpTo(domain.StepID("s3"))
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s3"), step.ID)
	assert.Equal(t, 2, nav.CurrentIndex())
}

func TestNavigator_JumpTo_NotFound(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	_, err := nav.JumpTo(domain.StepID("nonexistent"))
	assert.Error(t, err)
}

func TestNavigator_CurrentSection(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	sec := nav.CurrentSection()
	assert.Equal(t, domain.SectionID("sec-1"), sec.ID)

	nav.Next() // s2, still sec-1
	sec = nav.CurrentSection()
	assert.Equal(t, domain.SectionID("sec-1"), sec.ID)

	nav.Next() // s3, now sec-2
	sec = nav.CurrentSection()
	assert.Equal(t, domain.SectionID("sec-2"), sec.ID)
}

func TestNavigator_JumpToSection(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	step, err := nav.JumpToSection(domain.SectionID("sec-2"))
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s3"), step.ID)
}

func TestNavigator_EmptyWalkthrough(t *testing.T) {
	w := domain.NewWalkthrough("Empty", "", domain.ScopeFlow, "/repo")
	nav := domain.NewNavigator(w)

	_, err := nav.Current()
	assert.Error(t, err)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./domain/ -v -run TestNavigator`
Expected: FAIL — `NewNavigator` not defined

**Step 3: Write the navigator**

```go
// domain/navigator.go
package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNoSteps     = errors.New("walkthrough has no steps")
	ErrAtEnd       = errors.New("already at last step")
	ErrAtStart     = errors.New("already at first step")
	ErrStepNotFound    = errors.New("step not found")
	ErrSectionNotFound = errors.New("section not found")
)

type Navigator struct {
	walkthrough Walkthrough
	flatSteps   []Step
	sectionMap  map[int]SectionID // flatSteps index -> section ID
	index       int
}

func NewNavigator(w Walkthrough) *Navigator {
	var flat []Step
	sectionMap := make(map[int]SectionID)
	for _, sec := range w.Sections {
		for _, step := range sec.Steps {
			sectionMap[len(flat)] = sec.ID
			flat = append(flat, step)
		}
	}
	return &Navigator{
		walkthrough: w,
		flatSteps:   flat,
		sectionMap:  sectionMap,
		index:       0,
	}
}

func (n *Navigator) Current() (Step, error) {
	if len(n.flatSteps) == 0 {
		return Step{}, ErrNoSteps
	}
	return n.flatSteps[n.index], nil
}

func (n *Navigator) CurrentIndex() int {
	return n.index
}

func (n *Navigator) Next() (Step, error) {
	if len(n.flatSteps) == 0 {
		return Step{}, ErrNoSteps
	}
	if n.index >= len(n.flatSteps)-1 {
		return Step{}, ErrAtEnd
	}
	n.index++
	return n.flatSteps[n.index], nil
}

func (n *Navigator) Prev() (Step, error) {
	if len(n.flatSteps) == 0 {
		return Step{}, ErrNoSteps
	}
	if n.index <= 0 {
		return Step{}, ErrAtStart
	}
	n.index--
	return n.flatSteps[n.index], nil
}

func (n *Navigator) JumpTo(id StepID) (Step, error) {
	for i, step := range n.flatSteps {
		if step.ID == id {
			n.index = i
			return step, nil
		}
	}
	return Step{}, fmt.Errorf("%w: %s", ErrStepNotFound, id)
}

func (n *Navigator) JumpToSection(id SectionID) (Step, error) {
	for i, secID := range n.sectionMap {
		if secID == id {
			n.index = i
			return n.flatSteps[i], nil
		}
	}
	return Step{}, fmt.Errorf("%w: %s", ErrSectionNotFound, id)
}

func (n *Navigator) CurrentSection() Section {
	secID := n.sectionMap[n.index]
	for _, sec := range n.walkthrough.Sections {
		if sec.ID == secID {
			return sec
		}
	}
	return Section{}
}

func (n *Navigator) TotalSteps() int {
	return len(n.flatSteps)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./domain/ -v -run TestNavigator`
Expected: PASS — all 11 tests green

**Step 5: Commit**

```bash
git add domain/navigator.go domain/navigator_test.go
git commit -m "feat(domain): add navigator service with step traversal"
```

---

### Task 4: Domain Search Service

**Files:**
- Create: `domain/search.go`
- Create: `domain/search_test.go`

**Step 1: Write the failing test**

```go
// domain/search_test.go
package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tahrioui/code-walkthrough/domain"
)

func TestSearchIndex_Build(t *testing.T) {
	w := newTestWalkthrough()
	idx := domain.NewSearchIndex()
	idx.Build(w)

	assert.Equal(t, 3, idx.Size())
}

func TestSearchIndex_Search_ByTitle(t *testing.T) {
	w := newTestWalkthrough()
	idx := domain.NewSearchIndex()
	idx.Build(w)

	results := idx.Search("Step 1")
	assert.Len(t, results, 1)
	assert.Equal(t, domain.StepID("s1"), results[0].StepID)
}

func TestSearchIndex_Search_ByExplanation(t *testing.T) {
	w := newTestWalkthrough()
	idx := domain.NewSearchIndex()
	idx.Build(w)

	results := idx.Search("Exp 2")
	assert.Len(t, results, 1)
	assert.Equal(t, domain.StepID("s2"), results[0].StepID)
}

func TestSearchIndex_Search_CaseInsensitive(t *testing.T) {
	w := newTestWalkthrough()
	idx := domain.NewSearchIndex()
	idx.Build(w)

	results := idx.Search("step 1")
	assert.Len(t, results, 1)
}

func TestSearchIndex_Search_MultipleResults(t *testing.T) {
	w := newTestWalkthrough()
	idx := domain.NewSearchIndex()
	idx.Build(w)

	results := idx.Search("Step")
	assert.Len(t, results, 3)
}

func TestSearchIndex_Search_NoResults(t *testing.T) {
	w := newTestWalkthrough()
	idx := domain.NewSearchIndex()
	idx.Build(w)

	results := idx.Search("nonexistent query")
	assert.Empty(t, results)
}

func TestSearchIndex_Search_ByCodeSnippet(t *testing.T) {
	w := domain.NewWalkthrough("Test", "Desc", domain.ScopeFlow, "/repo")
	sec := domain.NewSection("sec-1", "Section", "")
	step := domain.NewStep("s1", "Router", "Handles requests")
	step.SetCodeSnippet(domain.CodeSnippet{
		FilePath: "router.go",
		Language: "go",
		Source:   "func HandleLogin() {}",
	})
	sec.AddStep(step)
	w.AddSection(sec)

	idx := domain.NewSearchIndex()
	idx.Build(w)

	results := idx.Search("HandleLogin")
	assert.Len(t, results, 1)
	assert.Equal(t, domain.StepID("s1"), results[0].StepID)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./domain/ -v -run TestSearchIndex`
Expected: FAIL — `NewSearchIndex` not defined

**Step 3: Write the search index**

```go
// domain/search.go
package domain

import "strings"

type SearchResult struct {
	StepID    StepID
	StepTitle string
	SectionID SectionID
	MatchText string
}

type searchEntry struct {
	stepID    StepID
	stepTitle string
	sectionID SectionID
	text      string // lowercased concatenation of searchable fields
}

type SearchIndex struct {
	entries []searchEntry
}

func NewSearchIndex() *SearchIndex {
	return &SearchIndex{}
}

func (si *SearchIndex) Build(w Walkthrough) {
	si.entries = nil
	for _, sec := range w.Sections {
		for _, step := range sec.Steps {
			var parts []string
			parts = append(parts, step.Title, step.Explanation)
			if step.CodeSnippet != nil {
				parts = append(parts, step.CodeSnippet.Source, step.CodeSnippet.FilePath)
			}
			if step.Diagram != nil {
				parts = append(parts, step.Diagram.Mermaid)
			}
			si.entries = append(si.entries, searchEntry{
				stepID:    step.ID,
				stepTitle: step.Title,
				sectionID: sec.ID,
				text:      strings.ToLower(strings.Join(parts, " ")),
			})
		}
	}
}

func (si *SearchIndex) Search(query string) []SearchResult {
	q := strings.ToLower(query)
	var results []SearchResult
	for _, entry := range si.entries {
		if strings.Contains(entry.text, q) {
			results = append(results, SearchResult{
				StepID:    entry.stepID,
				StepTitle: entry.stepTitle,
				SectionID: entry.sectionID,
			})
		}
	}
	return results
}

func (si *SearchIndex) Size() int {
	return len(si.entries)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./domain/ -v -run TestSearchIndex`
Expected: PASS — all 7 tests green

**Step 5: Commit**

```bash
git add domain/search.go domain/search_test.go
git commit -m "feat(domain): add search index service with full-text matching"
```

---

### Task 5: Port Interfaces

**Files:**
- Create: `port/inbound.go`
- Create: `port/outbound.go`

**Step 1: Write inbound ports**

```go
// port/inbound.go
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
```

**Step 2: Write outbound ports**

```go
// port/outbound.go
package port

import "github.com/tahrioui/code-walkthrough/domain"

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
```

**Step 3: Verify compilation**

Run: `go build ./port/`
Expected: no errors

**Step 4: Commit**

```bash
git add port/inbound.go port/outbound.go
git commit -m "feat(port): define inbound and outbound port interfaces"
```

---

### Task 6: JSON Schema

**Files:**
- Create: `schema/walkthrough.schema.json`
- Create: `adapter/filesystem.go`
- Create: `adapter/filesystem_test.go`

**Step 1: Write the JSON schema**

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Code Walkthrough",
  "type": "object",
  "required": ["title", "scope", "sections"],
  "properties": {
    "title": { "type": "string" },
    "description": { "type": "string" },
    "scope": { "enum": ["flow", "overview"] },
    "repository": { "type": "string" },
    "generatedAt": { "type": "string", "format": "date-time" },
    "sections": {
      "type": "array",
      "items": {
        "type": "object",
        "required": ["id", "title", "steps"],
        "properties": {
          "id": { "type": "string" },
          "title": { "type": "string" },
          "description": { "type": "string" },
          "steps": {
            "type": "array",
            "items": {
              "type": "object",
              "required": ["id", "title", "explanation"],
              "properties": {
                "id": { "type": "string" },
                "title": { "type": "string" },
                "explanation": { "type": "string" },
                "codeSnippet": {
                  "type": "object",
                  "required": ["filePath", "language", "source"],
                  "properties": {
                    "filePath": { "type": "string" },
                    "language": { "type": "string" },
                    "startLine": { "type": "integer" },
                    "endLine": { "type": "integer" },
                    "source": { "type": "string" }
                  }
                },
                "diagram": {
                  "type": "object",
                  "required": ["type", "mermaid"],
                  "properties": {
                    "type": { "enum": ["sequence", "flowchart", "classDiagram", "graph"] },
                    "mermaid": { "type": "string" }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
```

**Step 2: Write the filesystem adapter test**

```go
// adapter/filesystem_test.go
package adapter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahrioui/code-walkthrough/adapter"
	"github.com/tahrioui/code-walkthrough/domain"
)

const testJSON = `{
  "title": "Test Walkthrough",
  "description": "A test",
  "scope": "flow",
  "repository": "/test/repo",
  "generatedAt": "2026-03-01T12:00:00Z",
  "sections": [
    {
      "id": "sec-1",
      "title": "Section 1",
      "description": "First section",
      "steps": [
        {
          "id": "step-1",
          "title": "Step One",
          "explanation": "This is step one",
          "codeSnippet": {
            "filePath": "main.go",
            "language": "go",
            "startLine": 1,
            "endLine": 5,
            "source": "package main"
          },
          "diagram": {
            "type": "sequence",
            "mermaid": "sequenceDiagram\n    A->>B: hello"
          }
        }
      ]
    }
  ]
}`

func TestFileRepository_Read(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "walkthrough.json")
	err := os.WriteFile(path, []byte(testJSON), 0644)
	require.NoError(t, err)

	repo := adapter.NewFileRepository()
	data, err := repo.Read(path)

	require.NoError(t, err)
	assert.Contains(t, string(data), "Test Walkthrough")
}

func TestFileRepository_Write(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.json")

	repo := adapter.NewFileRepository()
	err := repo.Write(path, []byte(`{"title":"written"}`))
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), "written")
}

func TestFileRepository_Read_NotFound(t *testing.T) {
	repo := adapter.NewFileRepository()
	_, err := repo.Read("/nonexistent/path.json")
	assert.Error(t, err)
}

func TestJSONLoader_Load(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "walkthrough.json")
	err := os.WriteFile(path, []byte(testJSON), 0644)
	require.NoError(t, err)

	repo := adapter.NewFileRepository()
	loader := adapter.NewJSONLoader(repo)

	w, err := loader.Load(path)

	require.NoError(t, err)
	assert.Equal(t, "Test Walkthrough", w.Title)
	assert.Equal(t, domain.ScopeFlow, w.Scope)
	assert.Len(t, w.Sections, 1)
	assert.Len(t, w.Sections[0].Steps, 1)
	assert.Equal(t, "Step One", w.Sections[0].Steps[0].Title)
	assert.NotNil(t, w.Sections[0].Steps[0].CodeSnippet)
	assert.Equal(t, "main.go", w.Sections[0].Steps[0].CodeSnippet.FilePath)
	assert.NotNil(t, w.Sections[0].Steps[0].Diagram)
	assert.Equal(t, domain.DiagramSequence, w.Sections[0].Steps[0].Diagram.Type)
}

func TestJSONLoader_Load_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	err := os.WriteFile(path, []byte(`{not json`), 0644)
	require.NoError(t, err)

	repo := adapter.NewFileRepository()
	loader := adapter.NewJSONLoader(repo)

	_, err = loader.Load(path)
	assert.Error(t, err)
}
```

**Step 3: Run test to verify it fails**

Run: `go test ./adapter/ -v`
Expected: FAIL — types not defined

**Step 4: Write the filesystem adapter and JSON loader**

```go
// adapter/filesystem.go
package adapter

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/tahrioui/code-walkthrough/domain"
)

// FileRepository implements port.WalkthroughRepository
type FileRepository struct{}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (r *FileRepository) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (r *FileRepository) Write(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// JSON deserialization types

type jsonWalkthrough struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Scope       string        `json:"scope"`
	Repository  string        `json:"repository"`
	GeneratedAt string        `json:"generatedAt"`
	Sections    []jsonSection `json:"sections"`
}

type jsonSection struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Steps       []jsonStep `json:"steps"`
}

type jsonStep struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Explanation string           `json:"explanation"`
	CodeSnippet *jsonCodeSnippet `json:"codeSnippet,omitempty"`
	Diagram     *jsonDiagram     `json:"diagram,omitempty"`
}

type jsonCodeSnippet struct {
	FilePath  string `json:"filePath"`
	Language  string `json:"language"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
	Source    string `json:"source"`
}

type jsonDiagram struct {
	Type    string `json:"type"`
	Mermaid string `json:"mermaid"`
}

// JSONLoader implements port.WalkthroughLoader
type JSONLoader struct {
	repo *FileRepository
}

func NewJSONLoader(repo *FileRepository) *JSONLoader {
	return &JSONLoader{repo: repo}
}

func (l *JSONLoader) Load(source string) (domain.Walkthrough, error) {
	data, err := l.repo.Read(source)
	if err != nil {
		return domain.Walkthrough{}, fmt.Errorf("reading walkthrough file: %w", err)
	}

	var jw jsonWalkthrough
	if err := json.Unmarshal(data, &jw); err != nil {
		return domain.Walkthrough{}, fmt.Errorf("parsing walkthrough JSON: %w", err)
	}

	return toDomain(jw), nil
}

func toDomain(jw jsonWalkthrough) domain.Walkthrough {
	w := domain.Walkthrough{
		Title:       jw.Title,
		Description: jw.Description,
		Scope:       domain.Scope(jw.Scope),
		Repository:  jw.Repository,
	}

	if jw.GeneratedAt != "" {
		t, err := time.Parse(time.RFC3339, jw.GeneratedAt)
		if err == nil {
			w.GeneratedAt = t
		}
	}

	for _, js := range jw.Sections {
		sec := domain.NewSection(js.ID, js.Title, js.Description)
		for _, jst := range js.Steps {
			step := domain.NewStep(jst.ID, jst.Title, jst.Explanation)
			if jst.CodeSnippet != nil {
				step.SetCodeSnippet(domain.CodeSnippet{
					FilePath:  jst.CodeSnippet.FilePath,
					Language:  jst.CodeSnippet.Language,
					StartLine: jst.CodeSnippet.StartLine,
					EndLine:   jst.CodeSnippet.EndLine,
					Source:    jst.CodeSnippet.Source,
				})
			}
			if jst.Diagram != nil {
				step.SetDiagram(domain.Diagram{
					Type:    domain.DiagramType(jst.Diagram.Type),
					Mermaid: jst.Diagram.Mermaid,
				})
			}
			sec.AddStep(step)
		}
		w.AddSection(sec)
	}

	return w
}
```

**Step 5: Run test to verify it passes**

Run: `go test ./adapter/ -v`
Expected: PASS — all 5 tests green

**Step 6: Commit**

```bash
git add schema/walkthrough.schema.json adapter/filesystem.go adapter/filesystem_test.go
git commit -m "feat: add JSON schema and filesystem adapter with loader"
```

---

### Task 7: Bookmark Store Adapter

**Files:**
- Create: `adapter/bookmarkstore.go`
- Create: `adapter/bookmarkstore_test.go`

**Step 1: Write the failing test**

```go
// adapter/bookmarkstore_test.go
package adapter_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahrioui/code-walkthrough/adapter"
	"github.com/tahrioui/code-walkthrough/domain"
)

func TestJSONBookmarkStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bookmarks.json")
	store := adapter.NewJSONBookmarkStore(path)

	bookmarks := []domain.Bookmark{
		{StepID: "s1", CreatedAt: time.Now()},
		{StepID: "s3", CreatedAt: time.Now()},
	}

	err := store.Save(bookmarks)
	require.NoError(t, err)

	loaded, err := store.Load()
	require.NoError(t, err)
	assert.Len(t, loaded, 2)
	assert.Equal(t, domain.StepID("s1"), loaded[0].StepID)
	assert.Equal(t, domain.StepID("s3"), loaded[1].StepID)
}

func TestJSONBookmarkStore_Load_NoFile(t *testing.T) {
	store := adapter.NewJSONBookmarkStore("/nonexistent/bookmarks.json")

	loaded, err := store.Load()
	require.NoError(t, err)
	assert.Empty(t, loaded)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./adapter/ -v -run TestJSONBookmarkStore`
Expected: FAIL — `NewJSONBookmarkStore` not defined

**Step 3: Write the bookmark store**

```go
// adapter/bookmarkstore.go
package adapter

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/tahrioui/code-walkthrough/domain"
)

type jsonBookmark struct {
	StepID    string    `json:"stepId"`
	CreatedAt time.Time `json:"createdAt"`
}

// JSONBookmarkStore implements port.BookmarkStore
type JSONBookmarkStore struct {
	path string
}

func NewJSONBookmarkStore(path string) *JSONBookmarkStore {
	return &JSONBookmarkStore{path: path}
}

func (s *JSONBookmarkStore) Save(bookmarks []domain.Bookmark) error {
	jbs := make([]jsonBookmark, len(bookmarks))
	for i, b := range bookmarks {
		jbs[i] = jsonBookmark{
			StepID:    string(b.StepID),
			CreatedAt: b.CreatedAt,
		}
	}
	data, err := json.MarshalIndent(jbs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *JSONBookmarkStore) Load() ([]domain.Bookmark, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var jbs []jsonBookmark
	if err := json.Unmarshal(data, &jbs); err != nil {
		return nil, err
	}
	bookmarks := make([]domain.Bookmark, len(jbs))
	for i, jb := range jbs {
		bookmarks[i] = domain.Bookmark{
			StepID:    domain.StepID(jb.StepID),
			CreatedAt: jb.CreatedAt,
		}
	}
	return bookmarks, nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./adapter/ -v -run TestJSONBookmarkStore`
Expected: PASS

**Step 5: Commit**

```bash
git add adapter/bookmarkstore.go adapter/bookmarkstore_test.go
git commit -m "feat(adapter): add JSON bookmark store with persistence"
```

---

### Task 8: Application — Navigate Use Cases

**Files:**
- Create: `application/navigate.go`
- Create: `application/navigate_test.go`

**Step 1: Write the failing test**

```go
// application/navigate_test.go
package application_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahrioui/code-walkthrough/application"
	"github.com/tahrioui/code-walkthrough/domain"
)

func newTestWalkthrough() domain.Walkthrough {
	w := domain.NewWalkthrough("Test", "Desc", domain.ScopeFlow, "/repo")
	sec1 := domain.NewSection("sec-1", "Section 1", "")
	sec1.AddStep(domain.NewStep("s1", "Step 1", "Exp 1"))
	sec1.AddStep(domain.NewStep("s2", "Step 2", "Exp 2"))
	sec2 := domain.NewSection("sec-2", "Section 2", "")
	sec2.AddStep(domain.NewStep("s3", "Step 3", "Exp 3"))
	w.AddSection(sec1)
	w.AddSection(sec2)
	return w
}

func TestNavigateUseCase_Init(t *testing.T) {
	w := newTestWalkthrough()
	uc := application.NewNavigateUseCase(w)

	step, err := uc.Current()
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s1"), step.ID)
	assert.Equal(t, 0, uc.CurrentIndex())
	assert.Equal(t, 3, uc.TotalSteps())
}

func TestNavigateUseCase_StepForward(t *testing.T) {
	w := newTestWalkthrough()
	uc := application.NewNavigateUseCase(w)

	step, err := uc.StepForward()
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s2"), step.ID)
}

func TestNavigateUseCase_StepBackward(t *testing.T) {
	w := newTestWalkthrough()
	uc := application.NewNavigateUseCase(w)

	uc.StepForward()
	step, err := uc.StepBackward()
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s1"), step.ID)
}

func TestNavigateUseCase_JumpTo(t *testing.T) {
	w := newTestWalkthrough()
	uc := application.NewNavigateUseCase(w)

	step, err := uc.JumpTo(domain.StepID("s3"))
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s3"), step.ID)
}

func TestNavigateUseCase_JumpToSection(t *testing.T) {
	w := newTestWalkthrough()
	uc := application.NewNavigateUseCase(w)

	step, err := uc.JumpToSection(domain.SectionID("sec-2"))
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s3"), step.ID)
}

func TestNavigateUseCase_CurrentSection(t *testing.T) {
	w := newTestWalkthrough()
	uc := application.NewNavigateUseCase(w)

	sec := uc.CurrentSection()
	assert.Equal(t, domain.SectionID("sec-1"), sec.ID)
}

func TestNavigateUseCase_ViewTOC(t *testing.T) {
	w := newTestWalkthrough()
	uc := application.NewNavigateUseCase(w)

	sections := uc.ViewTOC()
	assert.Len(t, sections, 2)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./application/ -v -run TestNavigateUseCase`
Expected: FAIL

**Step 3: Write the navigate use case**

```go
// application/navigate.go
package application

import "github.com/tahrioui/code-walkthrough/domain"

type NavigateUseCase struct {
	nav         *domain.Navigator
	walkthrough domain.Walkthrough
}

func NewNavigateUseCase(w domain.Walkthrough) *NavigateUseCase {
	return &NavigateUseCase{
		nav:         domain.NewNavigator(w),
		walkthrough: w,
	}
}

func (uc *NavigateUseCase) Current() (domain.Step, error) {
	return uc.nav.Current()
}

func (uc *NavigateUseCase) StepForward() (domain.Step, error) {
	return uc.nav.Next()
}

func (uc *NavigateUseCase) StepBackward() (domain.Step, error) {
	return uc.nav.Prev()
}

func (uc *NavigateUseCase) JumpTo(id domain.StepID) (domain.Step, error) {
	return uc.nav.JumpTo(id)
}

func (uc *NavigateUseCase) JumpToSection(id domain.SectionID) (domain.Step, error) {
	return uc.nav.JumpToSection(id)
}

func (uc *NavigateUseCase) CurrentSection() domain.Section {
	return uc.nav.CurrentSection()
}

func (uc *NavigateUseCase) CurrentIndex() int {
	return uc.nav.CurrentIndex()
}

func (uc *NavigateUseCase) TotalSteps() int {
	return uc.nav.TotalSteps()
}

func (uc *NavigateUseCase) ViewTOC() []domain.Section {
	return uc.walkthrough.Sections
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./application/ -v -run TestNavigateUseCase`
Expected: PASS — all 7 tests green

**Step 5: Commit**

```bash
git add application/navigate.go application/navigate_test.go
git commit -m "feat(application): add navigation use cases"
```

---

### Task 9: Application — Search Use Cases

**Files:**
- Create: `application/search.go`
- Create: `application/search_test.go`

**Step 1: Write the failing test**

```go
// application/search_test.go
package application_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tahrioui/code-walkthrough/application"
	"github.com/tahrioui/code-walkthrough/domain"
)

func TestSearchUseCase_Search(t *testing.T) {
	w := newTestWalkthrough()
	uc := application.NewSearchUseCase(w)

	results := uc.Search("Step 1")
	assert.Len(t, results, 1)
	assert.Equal(t, domain.StepID("s1"), results[0].StepID)
}

func TestSearchUseCase_Search_NoResults(t *testing.T) {
	w := newTestWalkthrough()
	uc := application.NewSearchUseCase(w)

	results := uc.Search("nonexistent")
	assert.Empty(t, results)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./application/ -v -run TestSearchUseCase`
Expected: FAIL

**Step 3: Write the search use case**

```go
// application/search.go
package application

import "github.com/tahrioui/code-walkthrough/domain"

type SearchUseCase struct {
	index *domain.SearchIndex
}

func NewSearchUseCase(w domain.Walkthrough) *SearchUseCase {
	idx := domain.NewSearchIndex()
	idx.Build(w)
	return &SearchUseCase{index: idx}
}

func (uc *SearchUseCase) Search(query string) []domain.SearchResult {
	return uc.index.Search(query)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./application/ -v -run TestSearchUseCase`
Expected: PASS

**Step 5: Commit**

```bash
git add application/search.go application/search_test.go
git commit -m "feat(application): add search use case"
```

---

### Task 10: Application — Bookmark Use Cases

**Files:**
- Create: `application/bookmark.go`
- Create: `application/bookmark_test.go`

**Step 1: Write the failing test**

```go
// application/bookmark_test.go
package application_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahrioui/code-walkthrough/application"
	"github.com/tahrioui/code-walkthrough/domain"
)

type mockBookmarkStore struct {
	bookmarks []domain.Bookmark
}

func (m *mockBookmarkStore) Save(bookmarks []domain.Bookmark) error {
	m.bookmarks = bookmarks
	return nil
}

func (m *mockBookmarkStore) Load() ([]domain.Bookmark, error) {
	return m.bookmarks, nil
}

func TestBookmarkUseCase_Add(t *testing.T) {
	store := &mockBookmarkStore{}
	uc := application.NewBookmarkUseCase(store)

	err := uc.Add(domain.StepID("s1"))
	require.NoError(t, err)

	assert.True(t, uc.IsBookmarked(domain.StepID("s1")))
	assert.Len(t, uc.List(), 1)
}

func TestBookmarkUseCase_Add_Duplicate(t *testing.T) {
	store := &mockBookmarkStore{}
	uc := application.NewBookmarkUseCase(store)

	uc.Add(domain.StepID("s1"))
	err := uc.Add(domain.StepID("s1"))
	require.NoError(t, err)

	assert.Len(t, uc.List(), 1) // no duplicate
}

func TestBookmarkUseCase_Remove(t *testing.T) {
	store := &mockBookmarkStore{}
	uc := application.NewBookmarkUseCase(store)

	uc.Add(domain.StepID("s1"))
	err := uc.Remove(domain.StepID("s1"))
	require.NoError(t, err)

	assert.False(t, uc.IsBookmarked(domain.StepID("s1")))
	assert.Empty(t, uc.List())
}

func TestBookmarkUseCase_LoadFromStore(t *testing.T) {
	store := &mockBookmarkStore{
		bookmarks: []domain.Bookmark{
			{StepID: "s2"},
		},
	}
	uc := application.NewBookmarkUseCase(store)

	err := uc.LoadFromStore()
	require.NoError(t, err)
	assert.True(t, uc.IsBookmarked(domain.StepID("s2")))
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./application/ -v -run TestBookmarkUseCase`
Expected: FAIL

**Step 3: Write the bookmark use case**

```go
// application/bookmark.go
package application

import (
	"time"

	"github.com/tahrioui/code-walkthrough/domain"
	"github.com/tahrioui/code-walkthrough/port"
)

type BookmarkUseCase struct {
	store     port.BookmarkStore
	bookmarks []domain.Bookmark
	lookup    map[domain.StepID]bool
}

func NewBookmarkUseCase(store port.BookmarkStore) *BookmarkUseCase {
	return &BookmarkUseCase{
		store:  store,
		lookup: make(map[domain.StepID]bool),
	}
}

func (uc *BookmarkUseCase) Add(id domain.StepID) error {
	if uc.lookup[id] {
		return nil
	}
	b := domain.Bookmark{StepID: id, CreatedAt: time.Now()}
	uc.bookmarks = append(uc.bookmarks, b)
	uc.lookup[id] = true
	return uc.store.Save(uc.bookmarks)
}

func (uc *BookmarkUseCase) Remove(id domain.StepID) error {
	if !uc.lookup[id] {
		return nil
	}
	var filtered []domain.Bookmark
	for _, b := range uc.bookmarks {
		if b.StepID != id {
			filtered = append(filtered, b)
		}
	}
	uc.bookmarks = filtered
	delete(uc.lookup, id)
	return uc.store.Save(uc.bookmarks)
}

func (uc *BookmarkUseCase) List() []domain.Bookmark {
	return uc.bookmarks
}

func (uc *BookmarkUseCase) IsBookmarked(id domain.StepID) bool {
	return uc.lookup[id]
}

func (uc *BookmarkUseCase) LoadFromStore() error {
	loaded, err := uc.store.Load()
	if err != nil {
		return err
	}
	uc.bookmarks = loaded
	uc.lookup = make(map[domain.StepID]bool)
	for _, b := range loaded {
		uc.lookup[b.StepID] = true
	}
	return nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./application/ -v -run TestBookmarkUseCase`
Expected: PASS

**Step 5: Commit**

```bash
git add application/bookmark.go application/bookmark_test.go
git commit -m "feat(application): add bookmark use cases with store integration"
```

---

### Task 11: Application — Export Use Cases

**Files:**
- Create: `application/export.go`
- Create: `application/export_test.go`

**Step 1: Write the failing test**

```go
// application/export_test.go
package application_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahrioui/code-walkthrough/application"
	"github.com/tahrioui/code-walkthrough/domain"
	"github.com/tahrioui/code-walkthrough/port"
)

type mockDiagramRenderer struct{}

func (m *mockDiagramRenderer) Render(d domain.Diagram, width int) (string, error) {
	return "[diagram:" + d.Mermaid + "]", nil
}

func walkthroughWithDiagram() domain.Walkthrough {
	w := domain.NewWalkthrough("Export Test", "Testing export", domain.ScopeFlow, "/repo")
	sec := domain.NewSection("sec-1", "Auth", "Auth flow")
	step := domain.NewStep("s1", "Login", "User logs in")
	step.SetCodeSnippet(domain.CodeSnippet{
		FilePath:  "auth.go",
		Language:  "go",
		StartLine: 10,
		EndLine:   15,
		Source:    "func Login() {}",
	})
	step.SetDiagram(domain.Diagram{
		Type:    domain.DiagramSequence,
		Mermaid: "A->>B: login",
	})
	sec.AddStep(step)
	w.AddSection(sec)
	return w
}

func TestExportMarkdown(t *testing.T) {
	w := walkthroughWithDiagram()
	renderer := &mockDiagramRenderer{}
	uc := application.NewExportUseCase(w, renderer)

	var buf bytes.Buffer
	err := uc.Export(port.ExportMarkdown, &buf)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "# Export Test")
	assert.Contains(t, output, "## Auth")
	assert.Contains(t, output, "### Login")
	assert.Contains(t, output, "User logs in")
	assert.Contains(t, output, "```go")
	assert.Contains(t, output, "func Login() {}")
	assert.Contains(t, output, "[diagram:")
}

func TestExportHTML(t *testing.T) {
	w := walkthroughWithDiagram()
	renderer := &mockDiagramRenderer{}
	uc := application.NewExportUseCase(w, renderer)

	var buf bytes.Buffer
	err := uc.Export(port.ExportHTML, &buf)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "<html")
	assert.Contains(t, output, "Export Test")
	assert.Contains(t, output, "Login")
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./application/ -v -run TestExport`
Expected: FAIL

**Step 3: Write the export use case**

```go
// application/export.go
package application

import (
	"fmt"
	"io"

	"github.com/tahrioui/code-walkthrough/domain"
	"github.com/tahrioui/code-walkthrough/port"
)

type ExportUseCase struct {
	walkthrough domain.Walkthrough
	renderer    port.DiagramRenderer
}

func NewExportUseCase(w domain.Walkthrough, renderer port.DiagramRenderer) *ExportUseCase {
	return &ExportUseCase{walkthrough: w, renderer: renderer}
}

func (uc *ExportUseCase) Export(format port.ExportFormat, w io.Writer) error {
	switch format {
	case port.ExportMarkdown:
		return uc.exportMarkdown(w)
	case port.ExportHTML:
		return uc.exportHTML(w)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

func (uc *ExportUseCase) exportMarkdown(w io.Writer) error {
	wt := uc.walkthrough
	fmt.Fprintf(w, "# %s\n\n", wt.Title)
	if wt.Description != "" {
		fmt.Fprintf(w, "%s\n\n", wt.Description)
	}

	for _, sec := range wt.Sections {
		fmt.Fprintf(w, "## %s\n\n", sec.Title)
		if sec.Description != "" {
			fmt.Fprintf(w, "%s\n\n", sec.Description)
		}
		for _, step := range sec.Steps {
			fmt.Fprintf(w, "### %s\n\n", step.Title)
			fmt.Fprintf(w, "%s\n\n", step.Explanation)
			if step.CodeSnippet != nil {
				fmt.Fprintf(w, "```%s\n%s\n```\n\n", step.CodeSnippet.Language, step.CodeSnippet.Source)
			}
			if step.Diagram != nil {
				ascii, err := uc.renderer.Render(*step.Diagram, 80)
				if err == nil {
					fmt.Fprintf(w, "%s\n\n", ascii)
				}
			}
		}
	}
	return nil
}

func (uc *ExportUseCase) exportHTML(w io.Writer) error {
	wt := uc.walkthrough
	fmt.Fprintf(w, `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>%s</title>
<style>body{font-family:monospace;max-width:800px;margin:0 auto;padding:2rem}
pre{background:#f5f5f5;padding:1rem;overflow-x:auto}</style></head><body>
`, wt.Title)
	fmt.Fprintf(w, "<h1>%s</h1>\n", wt.Title)
	if wt.Description != "" {
		fmt.Fprintf(w, "<p>%s</p>\n", wt.Description)
	}

	for _, sec := range wt.Sections {
		fmt.Fprintf(w, "<h2>%s</h2>\n", sec.Title)
		for _, step := range sec.Steps {
			fmt.Fprintf(w, "<h3>%s</h3>\n", step.Title)
			fmt.Fprintf(w, "<p>%s</p>\n", step.Explanation)
			if step.CodeSnippet != nil {
				fmt.Fprintf(w, "<pre><code>%s</code></pre>\n", step.CodeSnippet.Source)
			}
			if step.Diagram != nil {
				ascii, err := uc.renderer.Render(*step.Diagram, 80)
				if err == nil {
					fmt.Fprintf(w, "<pre>%s</pre>\n", ascii)
				}
			}
		}
	}

	fmt.Fprint(w, "</body></html>\n")
	return nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./application/ -v -run TestExport`
Expected: PASS

**Step 5: Commit**

```bash
git add application/export.go application/export_test.go
git commit -m "feat(application): add markdown and HTML export use cases"
```

---

### Task 12: Mermaid ASCII Renderer Adapter

**Files:**
- Create: `adapter/mermaid.go`
- Create: `adapter/mermaid_test.go`

**Step 1: Write the failing test**

```go
// adapter/mermaid_test.go
package adapter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahrioui/code-walkthrough/adapter"
	"github.com/tahrioui/code-walkthrough/domain"
)

func TestMermaidRenderer_Sequence(t *testing.T) {
	r := adapter.NewMermaidRenderer()
	d := domain.Diagram{
		Type:    domain.DiagramSequence,
		Mermaid: "sequenceDiagram\n    Client->>Router: POST /login\n    Router->>Auth: validate",
	}

	output, err := r.Render(d, 60)
	require.NoError(t, err)
	assert.Contains(t, output, "Client")
	assert.Contains(t, output, "Router")
	assert.Contains(t, output, "Auth")
	assert.Contains(t, output, "POST /login")
}

func TestMermaidRenderer_Flowchart(t *testing.T) {
	r := adapter.NewMermaidRenderer()
	d := domain.Diagram{
		Type:    domain.DiagramFlowchart,
		Mermaid: "flowchart TD\n    A[Start] --> B[Process]\n    B --> C[End]",
	}

	output, err := r.Render(d, 60)
	require.NoError(t, err)
	assert.Contains(t, output, "Start")
	assert.Contains(t, output, "Process")
	assert.Contains(t, output, "End")
}

func TestMermaidRenderer_EmptyDiagram(t *testing.T) {
	r := adapter.NewMermaidRenderer()
	d := domain.Diagram{
		Type:    domain.DiagramSequence,
		Mermaid: "",
	}

	output, err := r.Render(d, 60)
	require.NoError(t, err)
	assert.Empty(t, output)
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./adapter/ -v -run TestMermaidRenderer`
Expected: FAIL

**Step 3: Write the Mermaid renderer**

This is a simplified ASCII renderer that parses common Mermaid patterns. A production version could shell out to `mmdc` or use a full parser. For now, we handle sequence diagrams and flowcharts with a basic text-based renderer.

```go
// adapter/mermaid.go
package adapter

import (
	"fmt"
	"strings"

	"github.com/tahrioui/code-walkthrough/domain"
)

// MermaidRenderer implements port.DiagramRenderer
type MermaidRenderer struct{}

func NewMermaidRenderer() *MermaidRenderer {
	return &MermaidRenderer{}
}

func (r *MermaidRenderer) Render(d domain.Diagram, width int) (string, error) {
	source := strings.TrimSpace(d.Mermaid)
	if source == "" {
		return "", nil
	}

	switch d.Type {
	case domain.DiagramSequence:
		return r.renderSequence(source, width), nil
	case domain.DiagramFlowchart:
		return r.renderFlowchart(source, width), nil
	default:
		return r.renderRaw(source, width), nil
	}
}

func (r *MermaidRenderer) renderSequence(source string, width int) string {
	lines := strings.Split(source, "\n")
	var participants []string
	var messages []struct{ from, to, label string }

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "sequenceDiagram" || line == "" {
			continue
		}
		// Parse "A->>B: label" or "A-->>B: label"
		for _, sep := range []string{"->>", "-->>", "->", "-->"} {
			if idx := strings.Index(line, sep); idx >= 0 {
				from := strings.TrimSpace(line[:idx])
				rest := line[idx+len(sep):]
				parts := strings.SplitN(rest, ":", 2)
				to := strings.TrimSpace(parts[0])
				label := ""
				if len(parts) > 1 {
					label = strings.TrimSpace(parts[1])
				}
				messages = append(messages, struct{ from, to, label string }{from, to, label})
				addUnique(&participants, from)
				addUnique(&participants, to)
				break
			}
		}
	}

	if len(participants) == 0 {
		return r.renderRaw(source, width)
	}

	var sb strings.Builder
	// Header
	sb.WriteString("  ")
	for i, p := range participants {
		if i > 0 {
			sb.WriteString("          ")
		}
		sb.WriteString(fmt.Sprintf("%-10s", p))
	}
	sb.WriteString("\n")

	// Separator
	sb.WriteString("  ")
	for i := range participants {
		if i > 0 {
			sb.WriteString("          ")
		}
		sb.WriteString("    |     ")
	}
	sb.WriteString("\n")

	// Messages
	for _, msg := range messages {
		fromIdx := indexOf(participants, msg.from)
		toIdx := indexOf(participants, msg.to)
		sb.WriteString("  ")
		for i := range participants {
			if i == fromIdx && fromIdx < toIdx {
				arrow := fmt.Sprintf("──▶ %s", msg.label)
				sb.WriteString(fmt.Sprintf("    |%s", padTo(arrow, 10*(toIdx-fromIdx)+5)))
			} else if i == toIdx && fromIdx < toIdx {
				sb.WriteString("|")
			} else if i == fromIdx && fromIdx > toIdx {
				sb.WriteString("|")
			} else if i == toIdx && fromIdx > toIdx {
				arrow := fmt.Sprintf("◀── %s", msg.label)
				sb.WriteString(fmt.Sprintf("    |%s", padTo(arrow, 10*(fromIdx-toIdx)+5)))
			} else if i > min(fromIdx, toIdx) && i < max(fromIdx, toIdx) {
				continue
			} else {
				sb.WriteString("    |     ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (r *MermaidRenderer) renderFlowchart(source string, width int) string {
	lines := strings.Split(source, "\n")
	var nodes []string
	var edges []struct{ from, to string }

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "flowchart") || strings.HasPrefix(line, "graph") || line == "" {
			continue
		}
		// Parse "A[Label] --> B[Label]" or "A --> B"
		for _, sep := range []string{"-->", "---"} {
			if idx := strings.Index(line, sep); idx >= 0 {
				from := extractNodeLabel(strings.TrimSpace(line[:idx]))
				to := extractNodeLabel(strings.TrimSpace(line[idx+len(sep):]))
				addUnique(&nodes, from)
				addUnique(&nodes, to)
				edges = append(edges, struct{ from, to string }{from, to})
				break
			}
		}
	}

	if len(nodes) == 0 {
		return r.renderRaw(source, width)
	}

	var sb strings.Builder
	for i, edge := range edges {
		fromBox := fmt.Sprintf("[ %s ]", edge.from)
		toBox := fmt.Sprintf("[ %s ]", edge.to)
		if i == 0 {
			sb.WriteString(fmt.Sprintf("  %s\n", fromBox))
		}
		sb.WriteString("      |\n")
		sb.WriteString("      v\n")
		sb.WriteString(fmt.Sprintf("  %s\n", toBox))
	}

	return sb.String()
}

func (r *MermaidRenderer) renderRaw(source string, width int) string {
	return source
}

func extractNodeLabel(s string) string {
	if idx := strings.Index(s, "["); idx >= 0 {
		end := strings.Index(s, "]")
		if end > idx {
			return s[idx+1 : end]
		}
	}
	if idx := strings.Index(s, "("); idx >= 0 {
		end := strings.Index(s, ")")
		if end > idx {
			return s[idx+1 : end]
		}
	}
	return s
}

func addUnique(slice *[]string, val string) {
	for _, v := range *slice {
		if v == val {
			return
		}
	}
	*slice = append(*slice, val)
}

func indexOf(slice []string, val string) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}

func padTo(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(" ", length-len(s))
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./adapter/ -v -run TestMermaidRenderer`
Expected: PASS

**Step 5: Commit**

```bash
git add adapter/mermaid.go adapter/mermaid_test.go
git commit -m "feat(adapter): add Mermaid to ASCII diagram renderer"
```

---

### Task 13: TUI Core — Bubble Tea App

**Files:**
- Create: `adapter/tui.go`
- Create: `adapter/tui_keymap.go`
- Create: `adapter/tui_styles.go`

**Step 1: Write the key map**

```go
// adapter/tui_keymap.go
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
```

**Step 2: Write the styles**

```go
// adapter/tui_styles.go
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
```

**Step 3: Write the TUI model**

```go
// adapter/tui.go
package adapter

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	walkthrough  domain.Walkthrough
	navigate     *application.NavigateUseCase
	search       *application.SearchUseCase
	bookmarks    *application.BookmarkUseCase
	renderer     *MermaidRenderer
	keys         KeyMap
	styles       Styles
	mode         viewMode
	showDiagram  bool
	searchInput  textinput.Model
	searchResults []domain.SearchResult
	tocCursor    int
	listCursor   int
	width        int
	height       int
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
```

**Step 4: Verify compilation**

Run: `go build ./adapter/`
Expected: no errors

**Step 5: Commit**

```bash
git add adapter/tui.go adapter/tui_keymap.go adapter/tui_styles.go
git commit -m "feat(adapter): add Bubble Tea TUI with navigation, search, bookmarks, help"
```

---

### Task 14: CLI Adapter — Cobra Commands

**Files:**
- Create: `adapter/cli.go`
- Modify: `cmd/walkthrough/main.go`

**Step 1: Write the CLI adapter**

```go
// adapter/cli.go
package adapter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tahrioui/code-walkthrough/application"
	"github.com/tahrioui/code-walkthrough/port"
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "walkthrough",
		Short: "Interactive code walkthrough viewer",
	}

	root.AddCommand(newViewCmd())
	root.AddCommand(newExportCmd())

	return root
}

func newViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view [file.json]",
		Short: "Open an interactive walkthrough in the TUI",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			repo := NewFileRepository()
			loader := NewJSONLoader(repo)

			w, err := loader.Load(path)
			if err != nil {
				return fmt.Errorf("loading walkthrough: %w", err)
			}

			nav := application.NewNavigateUseCase(w)
			srch := application.NewSearchUseCase(w)

			bmPath := filepath.Join(filepath.Dir(path), ".bookmarks.json")
			bmStore := NewJSONBookmarkStore(bmPath)
			bm := application.NewBookmarkUseCase(bmStore)
			bm.LoadFromStore()

			renderer := NewMermaidRenderer()

			model := NewModel(w, nav, srch, bm, renderer)
			return RunTUI(model)
		},
	}
}

func newExportCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "export [file.json] [output]",
		Short: "Export a walkthrough to markdown or HTML",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputPath := args[0]
			outputPath := args[1]

			repo := NewFileRepository()
			loader := NewJSONLoader(repo)

			w, err := loader.Load(inputPath)
			if err != nil {
				return fmt.Errorf("loading walkthrough: %w", err)
			}

			renderer := NewMermaidRenderer()
			exportUC := application.NewExportUseCase(w, renderer)

			f, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("creating output file: %w", err)
			}
			defer f.Close()

			exportFormat := port.ExportMarkdown
			if format == "html" {
				exportFormat = port.ExportHTML
			}

			if err := exportUC.Export(exportFormat, f); err != nil {
				return fmt.Errorf("exporting: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Exported to %s\n", outputPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "markdown", "Export format: markdown or html")
	return cmd
}
```

**Step 2: Update main.go**

```go
// cmd/walkthrough/main.go
package main

import (
	"fmt"
	"os"

	"github.com/tahrioui/code-walkthrough/adapter"
)

func main() {
	if err := adapter.NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

**Step 3: Verify it compiles**

Run: `go build ./cmd/walkthrough/`
Expected: no errors, binary created

**Step 4: Commit**

```bash
git add adapter/cli.go cmd/walkthrough/main.go
git commit -m "feat: add Cobra CLI with view and export commands"
```

---

### Task 15: Integration Test — End to End

**Files:**
- Create: `testdata/sample.json`
- Create: `adapter/cli_test.go`

**Step 1: Create sample walkthrough JSON**

```json
{
  "title": "Sample Walkthrough",
  "description": "A sample walkthrough for testing",
  "scope": "flow",
  "repository": "/test/repo",
  "generatedAt": "2026-03-01T12:00:00Z",
  "sections": [
    {
      "id": "sec-1",
      "title": "Entry Point",
      "description": "Where requests enter the system",
      "steps": [
        {
          "id": "step-1-1",
          "title": "Router receives request",
          "explanation": "The router matches the incoming HTTP request to a handler.",
          "codeSnippet": {
            "filePath": "router/routes.go",
            "language": "go",
            "startLine": 10,
            "endLine": 15,
            "source": "func SetupRoutes(r *mux.Router) {\n    r.HandleFunc(\"/api/login\", Login).Methods(\"POST\")\n}"
          },
          "diagram": {
            "type": "sequence",
            "mermaid": "sequenceDiagram\n    Client->>Router: POST /api/login\n    Router->>Handler: Login()"
          }
        },
        {
          "id": "step-1-2",
          "title": "Handler validates input",
          "explanation": "The login handler parses the request body and validates credentials."
        }
      ]
    },
    {
      "id": "sec-2",
      "title": "Authentication",
      "description": "How the auth system validates users",
      "steps": [
        {
          "id": "step-2-1",
          "title": "Token generation",
          "explanation": "After validation, a JWT token is generated and returned.",
          "codeSnippet": {
            "filePath": "auth/token.go",
            "language": "go",
            "startLine": 25,
            "endLine": 35,
            "source": "func GenerateToken(userID string) (string, error) {\n    claims := jwt.MapClaims{\"sub\": userID}\n    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)\n    return token.SignedString(secretKey)\n}"
          },
          "diagram": {
            "type": "flowchart",
            "mermaid": "flowchart TD\n    A[Validate Credentials] --> B[Generate JWT]\n    B --> C[Return Token]"
          }
        }
      ]
    }
  ]
}
```

**Step 2: Write the integration test for export**

```go
// adapter/cli_test.go
package adapter_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahrioui/code-walkthrough/adapter"
	"github.com/tahrioui/code-walkthrough/application"
	"github.com/tahrioui/code-walkthrough/port"
)

func TestEndToEnd_LoadAndExportMarkdown(t *testing.T) {
	repo := adapter.NewFileRepository()
	loader := adapter.NewJSONLoader(repo)

	w, err := loader.Load("../testdata/sample.json")
	require.NoError(t, err)

	assert.Equal(t, "Sample Walkthrough", w.Title)
	assert.Len(t, w.Sections, 2)
	assert.Equal(t, 3, w.TotalSteps())

	renderer := adapter.NewMermaidRenderer()
	exportUC := application.NewExportUseCase(w, renderer)

	var buf bytes.Buffer
	err = exportUC.Export(port.ExportMarkdown, &buf)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "# Sample Walkthrough")
	assert.Contains(t, output, "## Entry Point")
	assert.Contains(t, output, "### Router receives request")
	assert.Contains(t, output, "```go")
	assert.Contains(t, output, "SetupRoutes")
	assert.Contains(t, output, "## Authentication")
	assert.Contains(t, output, "GenerateToken")
}
```

**Step 3: Run test to verify it passes**

Run: `go test ./adapter/ -v -run TestEndToEnd`
Expected: PASS

**Step 4: Commit**

```bash
git add testdata/sample.json adapter/cli_test.go
git commit -m "test: add sample walkthrough data and end-to-end integration test"
```

---

### Task 16: Run Full Test Suite

**Step 1: Run all tests**

Run: `go test ./... -v`
Expected: ALL PASS across domain/, application/, adapter/

**Step 2: Build final binary**

Run: `go build -o walkthrough ./cmd/walkthrough/`
Expected: binary created successfully

**Step 3: Smoke test**

Run: `./walkthrough export testdata/sample.json /tmp/test-export.md -f markdown && head -20 /tmp/test-export.md`
Expected: markdown output with walkthrough content

**Step 4: Commit**

```bash
git commit --allow-empty -m "chore: verify full test suite and binary build"
```

---

## Summary

| Task | Component | Tests |
|------|-----------|-------|
| 1 | Project scaffolding | — |
| 2 | Domain model | 6 unit tests |
| 3 | Domain navigator | 11 unit tests |
| 4 | Domain search | 7 unit tests |
| 5 | Port interfaces | compile check |
| 6 | JSON schema + filesystem adapter | 5 integration tests |
| 7 | Bookmark store adapter | 2 integration tests |
| 8 | Navigate use cases | 7 unit tests |
| 9 | Search use cases | 2 unit tests |
| 10 | Bookmark use cases | 4 unit tests |
| 11 | Export use cases | 2 unit tests |
| 12 | Mermaid ASCII renderer | 3 unit tests |
| 13 | TUI core (Bubble Tea) | compile check |
| 14 | CLI adapter (Cobra) | compile check |
| 15 | End-to-end integration | 1 integration test |
| 16 | Full suite verification | all tests |
