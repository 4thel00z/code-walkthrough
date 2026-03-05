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
