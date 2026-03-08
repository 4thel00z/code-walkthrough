package application_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/4thel00z/code-walkthrough/application"
	"github.com/4thel00z/code-walkthrough/domain"
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
