package adapter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tahrioui/code-walkthrough/application"
	"github.com/tahrioui/code-walkthrough/port"
)

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "walkthrough",
		Short: "Interactive code walkthrough viewer",
	}

	root.AddCommand(newViewCmd())
	root.AddCommand(newExportCmd())

	return root
}

func newViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view [file.json]",
		Short: "Open an interactive walkthrough in the TUI",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]

			repo := NewFileRepository()
			loader := NewJSONLoader(repo)

			w, err := loader.Load(path)
			if err != nil {
				return fmt.Errorf("loading walkthrough: %w", err)
			}

			nav := application.NewNavigateUseCase(w)
			srch := application.NewSearchUseCase(w)

			bmPath := filepath.Join(filepath.Dir(path), ".bookmarks.json")
			bmStore := NewJSONBookmarkStore(bmPath)
			bm := application.NewBookmarkUseCase(bmStore)
			bm.LoadFromStore()

			renderer := NewMermaidRenderer()

			model := NewModel(w, nav, srch, bm, renderer)
			return RunTUI(model)
		},
	}
}

func newExportCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "export [file.json] [output]",
		Short: "Export a walkthrough to markdown or HTML",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputPath := args[0]
			outputPath := args[1]

			repo := NewFileRepository()
			loader := NewJSONLoader(repo)

			w, err := loader.Load(inputPath)
			if err != nil {
				return fmt.Errorf("loading walkthrough: %w", err)
			}

			renderer := NewMermaidRenderer()
			exportUC := application.NewExportUseCase(w, renderer)

			f, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("creating output file: %w", err)
			}
			defer f.Close()

			exportFormat := port.ExportMarkdown
			if format == "html" {
				exportFormat = port.ExportHTML
			}

			if err := exportUC.Export(exportFormat, f); err != nil {
				return fmt.Errorf("exporting: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Exported to %s\n", outputPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "markdown", "Export format: markdown or html")
	return cmd
}
