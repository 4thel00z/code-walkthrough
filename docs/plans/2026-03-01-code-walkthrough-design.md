# Code Walkthrough — Design Document

**Date:** 2026-03-01
**Status:** Approved

## Overview

A two-part system for generating and rendering interactive code walkthroughs:

1. **AI agent skill** — analyzes a codebase and produces structured walkthrough JSON
2. **Go CLI (Charmbracelet Bubble Tea)** — consumes JSON and renders an interactive TUI

Architecture: hexagonal (ports & adapters), DDD, TDD.

## Domain Model

Single file `domain/model.go`. All value objects, entities, and the aggregate root.

**Aggregate Root:** `Walkthrough` — title, description, scope, repository, generatedAt, ordered Sections.

**Entities:**
- `Section` — id, title, description, ordered Steps
- `Step` — id, title, explanation, optional CodeSnippet, optional Diagram

**Value Objects:**
- `CodeSnippet` — filePath, language, startLine, endLine, source
- `Explanation` — text
- `Diagram` — type (sequence|flowchart|classDiagram|graph), mermaid source
- `Bookmark` — stepID, createdAt
- `StepID`, `SectionID` — typed string identifiers

**Domain Services:**
- `Navigator` (`domain/navigator.go`) — tracks position, next/prev/jump
- `SearchIndex` (`domain/search.go`) — indexes steps for full-text search

## Hexagonal Architecture

### Dependency Rule

```
adapter → application → domain
adapter → port (implements)
application → port (calls)
domain → nothing
```

### Inbound Ports (`port/inbound.go`)

| Port | Methods |
|------|---------|
| `WalkthroughLoader` | `Load(source string) (Walkthrough, error)` |
| `NavigationPort` | `Next()`, `Prev()`, `JumpTo(StepID)`, `Current() Step` |
| `SearchPort` | `Search(query string) []SearchResult` |
| `BookmarkPort` | `Add(StepID)`, `Remove(StepID)`, `List() []Bookmark` |
| `ExportPort` | `Export(format ExportFormat, writer io.Writer) error` |

### Outbound Ports (`port/outbound.go`)

| Port | Methods |
|------|---------|
| `SourceRepository` | `Scan(path string) ([]SourceFile, error)`, `ReadFile(path string) ([]byte, error)` |
| `WalkthroughRepository` | `Read(path string) ([]byte, error)`, `Write(path string, data []byte) error` |
| `DiagramRenderer` | `Render(diagram Diagram, width int) (string, error)` |
| `Presenter` | `RenderStep(Step)`, `RenderDiagram(string)`, `RenderTOC([]Section)` |
| `LanguageDetector` | `Detect(files []string) (map[string]string, error)` |
| `CallGraphBuilder` | `Build(entrypoint string, files []SourceFile) (CallGraph, error)` |
| `DependencyAnalyzer` | `Analyze(files []SourceFile) (DependencyGraph, error)` |
| `ExplanationGenerator` | `Generate(step Step) (string, error)` |
| `BookmarkStore` | `Save([]Bookmark) error`, `Load() ([]Bookmark, error)` |
| `SchemaValidator` | `Validate(data []byte) error` |

## Use Cases

### Generation (`application/generate.go`)

| Use Case | Description |
|----------|-------------|
| `AnalyzeCodebaseUseCase` | Scans target repo, identifies languages, entry points, module boundaries |
| `TraceFlowUseCase` | Traces execution flow from a starting point across files |
| `GenerateOverviewUseCase` | Produces high-level architecture walkthrough |
| `BuildDiagramUseCase` | Constructs Mermaid diagrams from analyzed relationships |
| `ComposeWalkthroughUseCase` | Assembles steps into sections with prose, code, diagrams |
| `SerializeWalkthroughUseCase` | Validates and serializes walkthrough to JSON |

### Navigation (`application/navigate.go`)

| Use Case | Description |
|----------|-------------|
| `LoadWalkthroughUseCase` | Reads JSON, validates schema, hydrates aggregate |
| `InitSessionUseCase` | Initializes navigation, search index, loads bookmarks |
| `StepForwardUseCase` | Advance to next step, render |
| `StepBackwardUseCase` | Go to previous step, render |
| `JumpToStepUseCase` | Jump to specific step by ID |
| `JumpToSectionUseCase` | Jump to first step of a section |
| `ViewTOCUseCase` | Render table of contents |
| `RenderDiagramUseCase` | Convert Mermaid to ASCII, display |
| `ToggleDiagramUseCase` | Show/hide diagram for current step |

### Search (`application/search.go`)

| Use Case | Description |
|----------|-------------|
| `SearchStepsUseCase` | Full-text search across all steps |
| `SelectSearchResultUseCase` | Jump to a step from search results |

### Bookmarks (`application/bookmark.go`)

| Use Case | Description |
|----------|-------------|
| `AddBookmarkUseCase` | Bookmark current step |
| `RemoveBookmarkUseCase` | Remove bookmark from current step |
| `ListBookmarksUseCase` | Show all bookmarks |

### Export (`application/export.go`)

| Use Case | Description |
|----------|-------------|
| `ExportMarkdownUseCase` | Export walkthrough to markdown |
| `ExportHTMLUseCase` | Export walkthrough to HTML |

## Adapters

### Inbound
- `adapter/cli.go` — Cobra CLI commands (`generate`, `view`, `export`)
- `adapter/tui.go` — Bubble Tea model/update/view
- `adapter/tui_keymap.go` — key bindings
- `adapter/tui_styles.go` — Lip Gloss styles
- `adapter/tui_step.go` — step view component
- `adapter/tui_diagram.go` — diagram view component
- `adapter/tui_toc.go` — TOC view component
- `adapter/tui_search.go` — search view component
- `adapter/tui_bookmarks.go` — bookmark view component

### Outbound
- `adapter/filesystem.go` — source + walkthrough file I/O
- `adapter/mermaid.go` — Mermaid to ASCII rendering
- `adapter/analysis.go` — language detection, call graph, dependency analysis
- `adapter/bookmarkstore.go` — JSON bookmark persistence
- `adapter/markdown_export.go` — markdown exporter
- `adapter/html_export.go` — HTML exporter

## JSON Schema

Contract between skill (producer) and TUI (consumer). Lives at `schema/walkthrough.schema.json`.

```json
{
  "title": "string",
  "description": "string",
  "scope": "flow | overview",
  "repository": "string",
  "generatedAt": "ISO 8601 datetime",
  "sections": [
    {
      "id": "string",
      "title": "string",
      "description": "string",
      "steps": [
        {
          "id": "string",
          "title": "string",
          "explanation": "string",
          "codeSnippet": {
            "filePath": "string",
            "language": "string",
            "startLine": "number",
            "endLine": "number",
            "source": "string"
          },
          "diagram": {
            "type": "sequence | flowchart | classDiagram | graph",
            "mermaid": "string"
          }
        }
      ]
    }
  ]
}
```

`codeSnippet` and `diagram` are optional per step.

## TUI Key Bindings

| Key | Action |
|-----|--------|
| `j` / `↓` | Next step |
| `k` / `↑` | Previous step |
| `g` | Jump to section (opens TOC) |
| `d` | Toggle diagram |
| `/` | Open search |
| `enter` | Select result / TOC entry |
| `b` | Toggle bookmark on current step |
| `B` | List all bookmarks |
| `e` | Export menu |
| `q` / `ctrl+c` | Quit |
| `?` | Help overlay |

## TUI Layout

```
┌─────────────────────────────────────────────┐
│ [Section 2/5] Request Entry Point    [⏐3/12]│
├─────────────────────────────────────────────┤
│                                             │
│ Step title                                  │
│                                             │
│ Explanation text                            │
│                                             │
│ ┌─ file.go:42-58 ────────────────────────┐  │
│ │ code snippet with syntax highlighting  │  │
│ └────────────────────────────────────────┘  │
│                                             │
│ ┌─ diagram type ─────────────────────────┐  │
│ │ ASCII-rendered Mermaid diagram         │  │
│ └────────────────────────────────────────┘  │
│                                             │
├─────────────────────────────────────────────┤
│ j/k:navigate  g:toc  /:search  d:diagram  ?│
└─────────────────────────────────────────────┘
```

## Testing Strategy (TDD)

| Layer | Test Type | Approach |
|-------|-----------|----------|
| `domain/` | Unit | Table-driven, pure logic, no mocks |
| `application/` | Unit | Mock ports via interfaces |
| `adapter/mermaid.go` | Unit | Fixture inputs, expected ASCII output |
| `adapter/filesystem.go` | Integration | Temp directories |
| `adapter/tui*.go` | Integration | `teatest` package |
| `adapter/*_export.go` | Golden file | Compare against committed expected output |
| `schema/` | Validation | Valid + invalid JSON fixtures |

TDD order: domain first → application with mocked ports → adapters.

## Project Structure

```
code-walkthrough/
├── cmd/walkthrough/main.go
├── domain/
│   ├── model.go
│   ├── navigator.go
│   └── search.go
├── application/
│   ├── generate.go
│   ├── navigate.go
│   ├── search.go
│   ├── bookmark.go
│   └── export.go
├── port/
│   ├── inbound.go
│   └── outbound.go
├── adapter/
│   ├── cli.go
│   ├── tui.go
│   ├── tui_keymap.go
│   ├── tui_styles.go
│   ├── tui_step.go
│   ├── tui_diagram.go
│   ├── tui_toc.go
│   ├── tui_search.go
│   ├── tui_bookmarks.go
│   ├── filesystem.go
│   ├── mermaid.go
│   ├── analysis.go
│   ├── bookmarkstore.go
│   ├── markdown_export.go
│   └── html_export.go
├── schema/walkthrough.schema.json
├── skill.md
├── go.mod
└── README.md
```

## AI Agent Skill

`skill.md` instructs the agent to:
1. Accept scope — target flow or "overview"
2. Analyze codebase using Glob, Grep, Read
3. Build walkthrough JSON conforming to `schema/walkthrough.schema.json`
4. Write to `.walkthrough/walkthrough.json`
5. Auto-launch TUI binary (or `--json-only` to skip)
