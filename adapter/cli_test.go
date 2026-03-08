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
