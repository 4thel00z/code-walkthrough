package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tahrioui/code-walkthrough/domain"
)

func newTestWalkthrough() domain.Walkthrough {
	w := domain.NewWalkthrough("Test", "Desc", domain.ScopeFlow, "/repo")

	sec1 := domain.NewSection("sec-1", "Section 1", "First section")
	sec1.AddStep(domain.NewStep("s1", "Step 1", "Exp 1"))
	sec1.AddStep(domain.NewStep("s2", "Step 2", "Exp 2"))

	sec2 := domain.NewSection("sec-2", "Section 2", "Second section")
	sec2.AddStep(domain.NewStep("s3", "Step 3", "Exp 3"))

	w.AddSection(sec1)
	w.AddSection(sec2)
	return w
}

func TestNavigator_Init(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	step, err := nav.Current()
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s1"), step.ID)
	assert.Equal(t, 0, nav.CurrentIndex())
}

func TestNavigator_Next(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	step, err := nav.Next()
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s2"), step.ID)
}

func TestNavigator_Next_CrossesSection(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	nav.Next() // s2
	step, err := nav.Next() // s3, crosses into sec-2
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s3"), step.ID)
}

func TestNavigator_Next_AtEnd(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	nav.Next() // s2
	nav.Next() // s3
	_, err := nav.Next()
	assert.Error(t, err)
}

func TestNavigator_Prev(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	nav.Next() // s2
	step, err := nav.Prev()
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s1"), step.ID)
}

func TestNavigator_Prev_AtStart(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	_, err := nav.Prev()
	assert.Error(t, err)
}

func TestNavigator_JumpTo(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	step, err := nav.JumpTo(domain.StepID("s3"))
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s3"), step.ID)
	assert.Equal(t, 2, nav.CurrentIndex())
}

func TestNavigator_JumpTo_NotFound(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	_, err := nav.JumpTo(domain.StepID("nonexistent"))
	assert.Error(t, err)
}

func TestNavigator_CurrentSection(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	sec := nav.CurrentSection()
	assert.Equal(t, domain.SectionID("sec-1"), sec.ID)

	nav.Next() // s2, still sec-1
	sec = nav.CurrentSection()
	assert.Equal(t, domain.SectionID("sec-1"), sec.ID)

	nav.Next() // s3, now sec-2
	sec = nav.CurrentSection()
	assert.Equal(t, domain.SectionID("sec-2"), sec.ID)
}

func TestNavigator_JumpToSection(t *testing.T) {
	w := newTestWalkthrough()
	nav := domain.NewNavigator(w)

	step, err := nav.JumpToSection(domain.SectionID("sec-2"))
	require.NoError(t, err)
	assert.Equal(t, domain.StepID("s3"), step.ID)
}

func TestNavigator_EmptyWalkthrough(t *testing.T) {
	w := domain.NewWalkthrough("Empty", "", domain.ScopeFlow, "/repo")
	nav := domain.NewNavigator(w)

	_, err := nav.Current()
	assert.Error(t, err)
}
