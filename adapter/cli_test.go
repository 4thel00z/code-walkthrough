package adapter_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/4thel00z/code-walkthrough/adapter"
	"github.com/4thel00z/code-walkthrough/application"
	"github.com/4thel00z/code-walkthrough/port"
	"github.com/4thel00z/code-walkthrough/skilldata"
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

func TestInstallCmd_WritesFiles(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "install-test")

	cmd := adapter.NewRootCmd([]byte("# Test Skill"), []byte(`{"test":true}`))
	cmd.SetArgs([]string{"install", "--dir", dir})

	err := cmd.Execute()
	require.NoError(t, err)

	gotSkill, err := os.ReadFile(filepath.Join(dir, "SKILL.md"))
	require.NoError(t, err)
	assert.Equal(t, "# Test Skill", string(gotSkill))

	gotSchema, err := os.ReadFile(filepath.Join(dir, "walkthrough.schema.json"))
	require.NoError(t, err)
	assert.Equal(t, `{"test":true}`, string(gotSchema))
}

func TestEmbeddedSchema_MatchesCanonical(t *testing.T) {
	canonical, err := os.ReadFile("../schema/walkthrough.schema.json")
	require.NoError(t, err)

	assert.Equal(t, string(canonical), string(skilldata.SchemaJSON),
		"skilldata/walkthrough.schema.json is out of sync with schema/walkthrough.schema.json")
}
