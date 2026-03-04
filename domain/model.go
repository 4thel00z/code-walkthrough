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
