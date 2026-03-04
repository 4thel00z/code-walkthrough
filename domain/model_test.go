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
