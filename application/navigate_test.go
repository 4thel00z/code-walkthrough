package application_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/4thel00z/code-walkthrough/application"
	"github.com/4thel00z/code-walkthrough/domain"
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
