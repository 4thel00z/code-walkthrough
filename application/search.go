package application

import "github.com/4thel00z/code-walkthrough/domain"

type SearchUseCase struct {
	index *domain.SearchIndex
}

func NewSearchUseCase(w domain.Walkthrough) *SearchUseCase {
	idx := domain.NewSearchIndex()
	idx.Build(w)
	return &SearchUseCase{index: idx}
}

func (uc *SearchUseCase) Search(query string) []domain.SearchResult {
	return uc.index.Search(query)
}
