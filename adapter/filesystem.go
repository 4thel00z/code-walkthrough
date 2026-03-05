package adapter

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/tahrioui/code-walkthrough/domain"
)

// FileRepository implements port.WalkthroughRepository
type FileRepository struct{}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (r *FileRepository) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (r *FileRepository) Write(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// JSON deserialization types

type jsonWalkthrough struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Scope       string        `json:"scope"`
	Repository  string        `json:"repository"`
	GeneratedAt string        `json:"generatedAt"`
	Sections    []jsonSection `json:"sections"`
}

type jsonSection struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Steps       []jsonStep `json:"steps"`
}

type jsonStep struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Explanation string           `json:"explanation"`
	CodeSnippet *jsonCodeSnippet `json:"codeSnippet,omitempty"`
	Diagram     *jsonDiagram     `json:"diagram,omitempty"`
}

type jsonCodeSnippet struct {
	FilePath  string `json:"filePath"`
	Language  string `json:"language"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
	Source    string `json:"source"`
}

type jsonDiagram struct {
	Type    string `json:"type"`
	Mermaid string `json:"mermaid"`
}

// JSONLoader implements port.WalkthroughLoader
type JSONLoader struct {
	repo *FileRepository
}

func NewJSONLoader(repo *FileRepository) *JSONLoader {
	return &JSONLoader{repo: repo}
}

func (l *JSONLoader) Load(source string) (domain.Walkthrough, error) {
	data, err := l.repo.Read(source)
	if err != nil {
		return domain.Walkthrough{}, fmt.Errorf("reading walkthrough file: %w", err)
	}

	var jw jsonWalkthrough
	if err := json.Unmarshal(data, &jw); err != nil {
		return domain.Walkthrough{}, fmt.Errorf("parsing walkthrough JSON: %w", err)
	}

	return toDomain(jw), nil
}

func toDomain(jw jsonWalkthrough) domain.Walkthrough {
	w := domain.Walkthrough{
		Title:       jw.Title,
		Description: jw.Description,
		Scope:       domain.Scope(jw.Scope),
		Repository:  jw.Repository,
	}

	if jw.GeneratedAt != "" {
		t, err := time.Parse(time.RFC3339, jw.GeneratedAt)
		if err == nil {
			w.GeneratedAt = t
		}
	}

	for _, js := range jw.Sections {
		sec := domain.NewSection(js.ID, js.Title, js.Description)
		for _, jst := range js.Steps {
			step := domain.NewStep(jst.ID, jst.Title, jst.Explanation)
			if jst.CodeSnippet != nil {
				step.SetCodeSnippet(domain.CodeSnippet{
					FilePath:  jst.CodeSnippet.FilePath,
					Language:  jst.CodeSnippet.Language,
					StartLine: jst.CodeSnippet.StartLine,
					EndLine:   jst.CodeSnippet.EndLine,
					Source:    jst.CodeSnippet.Source,
				})
			}
			if jst.Diagram != nil {
				step.SetDiagram(domain.Diagram{
					Type:    domain.DiagramType(jst.Diagram.Type),
					Mermaid: jst.Diagram.Mermaid,
				})
			}
			sec.AddStep(step)
		}
		w.AddSection(sec)
	}

	return w
}
