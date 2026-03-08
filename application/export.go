package application

import (
	"fmt"
	"io"

	"github.com/4thel00z/code-walkthrough/domain"
	"github.com/4thel00z/code-walkthrough/port"
)

type ExportUseCase struct {
	walkthrough domain.Walkthrough
	renderer    port.DiagramRenderer
}

func NewExportUseCase(w domain.Walkthrough, renderer port.DiagramRenderer) *ExportUseCase {
	return &ExportUseCase{walkthrough: w, renderer: renderer}
}

func (uc *ExportUseCase) Export(format port.ExportFormat, w io.Writer) error {
	switch format {
	case port.ExportMarkdown:
		return uc.exportMarkdown(w)
	case port.ExportHTML:
		return uc.exportHTML(w)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

func (uc *ExportUseCase) exportMarkdown(w io.Writer) error {
	wt := uc.walkthrough
	fmt.Fprintf(w, "# %s\n\n", wt.Title)
	if wt.Description != "" {
		fmt.Fprintf(w, "%s\n\n", wt.Description)
	}

	for _, sec := range wt.Sections {
		fmt.Fprintf(w, "## %s\n\n", sec.Title)
		if sec.Description != "" {
			fmt.Fprintf(w, "%s\n\n", sec.Description)
		}
		for _, step := range sec.Steps {
			fmt.Fprintf(w, "### %s\n\n", step.Title)
			fmt.Fprintf(w, "%s\n\n", step.Explanation)
			if step.CodeSnippet != nil {
				fmt.Fprintf(w, "```%s\n%s\n```\n\n", step.CodeSnippet.Language, step.CodeSnippet.Source)
			}
			if step.Diagram != nil {
				ascii, err := uc.renderer.Render(*step.Diagram, 80)
				if err == nil {
					fmt.Fprintf(w, "%s\n\n", ascii)
				}
			}
		}
	}
	return nil
}

func (uc *ExportUseCase) exportHTML(w io.Writer) error {
	wt := uc.walkthrough
	fmt.Fprintf(w, `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>%s</title>
<style>body{font-family:monospace;max-width:800px;margin:0 auto;padding:2rem}
pre{background:#f5f5f5;padding:1rem;overflow-x:auto}</style></head><body>
`, wt.Title)
	fmt.Fprintf(w, "<h1>%s</h1>\n", wt.Title)
	if wt.Description != "" {
		fmt.Fprintf(w, "<p>%s</p>\n", wt.Description)
	}

	for _, sec := range wt.Sections {
		fmt.Fprintf(w, "<h2>%s</h2>\n", sec.Title)
		for _, step := range sec.Steps {
			fmt.Fprintf(w, "<h3>%s</h3>\n", step.Title)
			fmt.Fprintf(w, "<p>%s</p>\n", step.Explanation)
			if step.CodeSnippet != nil {
				fmt.Fprintf(w, "<pre><code>%s</code></pre>\n", step.CodeSnippet.Source)
			}
			if step.Diagram != nil {
				ascii, err := uc.renderer.Render(*step.Diagram, 80)
				if err == nil {
					fmt.Fprintf(w, "<pre>%s</pre>\n", ascii)
				}
			}
		}
	}

	fmt.Fprint(w, "</body></html>\n")
	return nil
}
