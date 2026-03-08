package adapter

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/4thel00z/code-walkthrough/domain"
)

type jsonBookmark struct {
	StepID    string    `json:"stepId"`
	CreatedAt time.Time `json:"createdAt"`
}

// JSONBookmarkStore implements port.BookmarkStore
type JSONBookmarkStore struct {
	path string
}

func NewJSONBookmarkStore(path string) *JSONBookmarkStore {
	return &JSONBookmarkStore{path: path}
}

func (s *JSONBookmarkStore) Save(bookmarks []domain.Bookmark) error {
	jbs := make([]jsonBookmark, len(bookmarks))
	for i, b := range bookmarks {
		jbs[i] = jsonBookmark{
			StepID:    string(b.StepID),
			CreatedAt: b.CreatedAt,
		}
	}
	data, err := json.MarshalIndent(jbs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *JSONBookmarkStore) Load() ([]domain.Bookmark, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var jbs []jsonBookmark
	if err := json.Unmarshal(data, &jbs); err != nil {
		return nil, err
	}
	bookmarks := make([]domain.Bookmark, len(jbs))
	for i, jb := range jbs {
		bookmarks[i] = domain.Bookmark{
			StepID:    domain.StepID(jb.StepID),
			CreatedAt: jb.CreatedAt,
		}
	}
	return bookmarks, nil
}
