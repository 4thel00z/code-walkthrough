<p align="center">
  <img src="logo.png" alt="Code Walkthrough" width="400">
</p>

<h1 align="center">Code Walkthrough</h1>

<p align="center">
  <strong>Interactive, visual code exploration for the AI era</strong>
</p>

<p align="center">
  <code>AI-Powered</code> · <code>Terminal-Native</code> · <code>Diagram-Enhanced</code>
</p>

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go" alt="Go 1.25"></a>
  <a href="#license"><img src="https://img.shields.io/badge/license-MIT-blue?style=flat-square" alt="License"></a>
</p>

---

An AI agent skill and interactive TUI that generates step-by-step code walkthroughs enriched with diagrams, code snippets, and full-text search. Built for developers who need to understand codebases — fast.

## Features

- **AI-Driven Analysis** — A Claude Code skill analyzes your codebase and produces structured walkthrough JSON, tracing execution flows or mapping architecture
- **Interactive TUI** — Navigate walkthroughs in the terminal with keyboard-driven step browsing, table of contents, and bookmarks
- **Mermaid Diagrams** — Sequence, flowchart, and class diagrams rendered as ASCII art directly in your terminal
- **Full-Text Search** — Instantly search across titles, explanations, code snippets, and diagrams
- **Export** — Generate Markdown or HTML documentation from any walkthrough
- **Bookmarks** — Save and revisit important steps across sessions

## Installation

### From source

```bash
go install github.com/your-org/code-walkthrough/cmd/walkthrough@latest
```

### Build locally

```bash
git clone https://github.com/your-org/code-walkthrough.git
cd code-walkthrough
go build -o walkthrough ./cmd/walkthrough
```

## Quickstart

### 1. Install the Claude Code skill

```bash
walkthrough install
```

This installs the AI agent skill and JSON schema to `.claude/skills/code-walkthrough/`.

### 2. Generate a walkthrough

In Claude Code, invoke the skill:

```
/code-walkthrough
```

Choose a scope:

| Scope | Description |
|-------|-------------|
| `flow` | Trace a single request through the codebase end-to-end |
| `overview` | Map the high-level architecture and component relationships |

The agent analyzes the repository and writes a walkthrough to `.walkthrough/walkthrough.json`.

### 3. View interactively

```bash
walkthrough view .walkthrough/walkthrough.json
```

### 4. Export to documentation

```bash
walkthrough export walkthrough.json output.md --format=markdown
walkthrough export walkthrough.json output.html --format=html
```

## Keybindings

| Key | Action |
|-----|--------|
| `j` / `↓` | Next step |
| `k` / `↑` | Previous step |
| `g` | Table of contents |
| `/` | Search |
| `d` | Toggle diagram |
| `b` | Toggle bookmark |
| `e` | View bookmarks |
| `?` | Help |
| `q` | Quit |

## Architecture

The project follows hexagonal architecture with clear separation of concerns:

```
cmd/walkthrough/       Entry point
domain/                Core entities, navigation, search (zero dependencies)
application/           Use cases — navigate, search, bookmark, export, install
port/                  Inbound & outbound interfaces
adapter/               CLI (Cobra), TUI (Bubble Tea), filesystem, Mermaid renderer
skilldata/             Embedded AI agent skill & JSON schema
schema/                Walkthrough JSON schema
```

### Walkthrough Schema

Walkthroughs conform to a [JSON schema](schema/walkthrough.schema.json). The core structure:

```
Walkthrough
├── title, scope (flow | overview), repository
└── sections[]
    ├── id, title, description
    └── steps[]
        ├── id, title, explanation
        ├── codeSnippet? — filePath, language, lines, source
        └── diagram? — type (sequence | flowchart | classDiagram), mermaid DSL
```

## License

MIT
