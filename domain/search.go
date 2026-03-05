package domain

import "strings"

type SearchResult struct {
	StepID    StepID
	StepTitle string
	SectionID SectionID
	MatchText string
}

type searchEntry struct {
	stepID    StepID
	stepTitle string
	sectionID SectionID
	text      string // lowercased concatenation of searchable fields
}

type SearchIndex struct {
	entries []searchEntry
}

func NewSearchIndex() *SearchIndex {
	return &SearchIndex{}
}

func (si *SearchIndex) Build(w Walkthrough) {
	si.entries = nil
	for _, sec := range w.Sections {
		for _, step := range sec.Steps {
			var parts []string
			parts = append(parts, step.Title, step.Explanation)
			if step.CodeSnippet != nil {
				parts = append(parts, step.CodeSnippet.Source, step.CodeSnippet.FilePath)
			}
			if step.Diagram != nil {
				parts = append(parts, step.Diagram.Mermaid)
			}
			si.entries = append(si.entries, searchEntry{
				stepID:    step.ID,
				stepTitle: step.Title,
				sectionID: sec.ID,
				text:      strings.ToLower(strings.Join(parts, " ")),
			})
		}
	}
}

func (si *SearchIndex) Search(query string) []SearchResult {
	q := strings.ToLower(query)
	var results []SearchResult
	for _, entry := range si.entries {
		if strings.Contains(entry.text, q) {
			results = append(results, SearchResult{
				StepID:    entry.stepID,
				StepTitle: entry.stepTitle,
				SectionID: entry.sectionID,
			})
		}
	}
	return results
}

func (si *SearchIndex) Size() int {
	return len(si.entries)
}
