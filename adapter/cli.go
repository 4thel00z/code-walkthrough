package adapter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/4thel00z/code-walkthrough/application"
	"github.com/4thel00z/code-walkthrough/port"
)

func NewRootCmd(embeddedSkill, embeddedSchema []byte) *cobra.Command {
	root := &cobra.Command{
		Use:   "walkthrough",
		Short: "Interactive code walkthrough viewer",
	}

	root.AddCommand(newViewCmd())
	root.AddCommand(newExportCmd())
	root.AddCommand(newInstallCmd(embeddedSkill, embeddedSchema))

	return root
}

const defaultWalkthroughPath = ".walkthrough/walkthrough.json"

func newViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view [file.json]",
		Short: "Open an interactive walkthrough in the TUI",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := defaultWalkthroughPath
			if len(args) > 0 {
				path = args[0]
			}

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

func newInstallCmd(skill, schema []byte) *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the code-walkthrough skill for Claude Code",
		RunE: func(cmd *cobra.Command, args []string) error {
			installer := NewFileSkillInstaller()
			uc := application.NewInstallSkillUseCase(installer, skill, schema)

			if dir == "" {
				dir = uc.DefaultInstallDir()
			}

			if err := uc.Install(dir); err != nil {
				return fmt.Errorf("installing skill: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Skill installed to %s\n", dir)
			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", "", "Installation directory (default: .claude/skills/code-walkthrough)")
	return cmd
}
