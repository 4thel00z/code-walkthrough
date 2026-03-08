package application

import (
	"time"

	"github.com/tahrioui/code-walkthrough/domain"
	"github.com/tahrioui/code-walkthrough/port"
)

type BookmarkUseCase struct {
	store     port.BookmarkStore
	bookmarks []domain.Bookmark
	lookup    map[domain.StepID]bool
}

func NewBookmarkUseCase(store port.BookmarkStore) *BookmarkUseCase {
	return &BookmarkUseCase{
		store:  store,
		lookup: make(map[domain.StepID]bool),
	}
}

func (uc *BookmarkUseCase) Add(id domain.StepID) error {
	if uc.lookup[id] {
		return nil
	}
	b := domain.Bookmark{StepID: id, CreatedAt: time.Now()}
	uc.bookmarks = append(uc.bookmarks, b)
	uc.lookup[id] = true
	return uc.store.Save(uc.bookmarks)
}

func (uc *BookmarkUseCase) Remove(id domain.StepID) error {
	if !uc.lookup[id] {
		return nil
	}
	var filtered []domain.Bookmark
	for _, b := range uc.bookmarks {
		if b.StepID != id {
			filtered = append(filtered, b)
		}
	}
	uc.bookmarks = filtered
	delete(uc.lookup, id)
	return uc.store.Save(uc.bookmarks)
}

func (uc *BookmarkUseCase) List() []domain.Bookmark {
	return uc.bookmarks
}

func (uc *BookmarkUseCase) IsBookmarked(id domain.StepID) bool {
	return uc.lookup[id]
}

func (uc *BookmarkUseCase) LoadFromStore() error {
	loaded, err := uc.store.Load()
	if err != nil {
		return err
	}
	uc.bookmarks = loaded
	uc.lookup = make(map[domain.StepID]bool)
	for _, b := range loaded {
		uc.lookup[b.StepID] = true
	}
	return nil
}
