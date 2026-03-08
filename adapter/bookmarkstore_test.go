package adapter_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/4thel00z/code-walkthrough/adapter"
	"github.com/4thel00z/code-walkthrough/domain"
)

func TestJSONBookmarkStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bookmarks.json")
	store := adapter.NewJSONBookmarkStore(path)

	bookmarks := []domain.Bookmark{
		{StepID: "s1", CreatedAt: time.Now()},
		{StepID: "s3", CreatedAt: time.Now()},
	}

	err := store.Save(bookmarks)
	require.NoError(t, err)

	loaded, err := store.Load()
	require.NoError(t, err)
	assert.Len(t, loaded, 2)
	assert.Equal(t, domain.StepID("s1"), loaded[0].StepID)
	assert.Equal(t, domain.StepID("s3"), loaded[1].StepID)
}

func TestJSONBookmarkStore_Load_NoFile(t *testing.T) {
	store := adapter.NewJSONBookmarkStore("/nonexistent/bookmarks.json")

	loaded, err := store.Load()
	require.NoError(t, err)
	assert.Empty(t, loaded)
}
