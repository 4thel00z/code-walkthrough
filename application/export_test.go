package application_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/4thel00z/code-walkthrough/application"
	"github.com/4thel00z/code-walkthrough/domain"
	"github.com/4thel00z/code-walkthrough/port"
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
