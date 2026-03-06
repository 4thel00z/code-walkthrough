package application

import "github.com/tahrioui/code-walkthrough/domain"

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
