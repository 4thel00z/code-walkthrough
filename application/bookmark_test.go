package application_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/4thel00z/code-walkthrough/application"
	"github.com/4thel00z/code-walkthrough/domain"
)

type mockBookmarkStore struct {
	bookmarks []domain.Bookmark
}

func (m *mockBookmarkStore) Save(bookmarks []domain.Bookmark) error {
	m.bookmarks = bookmarks
	return nil
}

func (m *mockBookmarkStore) Load() ([]domain.Bookmark, error) {
	return m.bookmarks, nil
}

func TestBookmarkUseCase_Add(t *testing.T) {
	store := &mockBookmarkStore{}
	uc := application.NewBookmarkUseCase(store)

	err := uc.Add(domain.StepID("s1"))
	require.NoError(t, err)

	assert.True(t, uc.IsBookmarked(domain.StepID("s1")))
	assert.Len(t, uc.List(), 1)
}

func TestBookmarkUseCase_Add_Duplicate(t *testing.T) {
	store := &mockBookmarkStore{}
	uc := application.NewBookmarkUseCase(store)

	uc.Add(domain.StepID("s1"))
	err := uc.Add(domain.StepID("s1"))
	require.NoError(t, err)

	assert.Len(t, uc.List(), 1) // no duplicate
}

func TestBookmarkUseCase_Remove(t *testing.T) {
	store := &mockBookmarkStore{}
	uc := application.NewBookmarkUseCase(store)

	uc.Add(domain.StepID("s1"))
	err := uc.Remove(domain.StepID("s1"))
	require.NoError(t, err)

	assert.False(t, uc.IsBookmarked(domain.StepID("s1")))
	assert.Empty(t, uc.List())
}

func TestBookmarkUseCase_LoadFromStore(t *testing.T) {
	store := &mockBookmarkStore{
		bookmarks: []domain.Bookmark{
			{StepID: "s2"},
		},
	}
	uc := application.NewBookmarkUseCase(store)

	err := uc.LoadFromStore()
	require.NoError(t, err)
	assert.True(t, uc.IsBookmarked(domain.StepID("s2")))
}
