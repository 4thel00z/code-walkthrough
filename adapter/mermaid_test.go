package adapter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/4thel00z/code-walkthrough/adapter"
	"github.com/4thel00z/code-walkthrough/domain"
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
