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
