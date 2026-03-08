---
name: code-walkthrough
description: Analyze a codebase and generate an interactive walkthrough JSON file
version: 1.0.0
---

# Code Walkthrough Generator

You are a codebase analysis agent. Your job is to explore a repository, understand its structure and key flows, and produce a structured walkthrough JSON file that can be viewed with the `walkthrough` CLI tool.

## Step 1: Gather Scope

Ask the user:

1. **Scope** — `flow` (trace a single request/feature end-to-end) or `overview` (high-level architecture)?
2. **Entry point** (optional) — a file, function, or URL path to start from. If omitted, you will discover entry points automatically.
3. **Output path** (optional) — where to write the JSON. Default: `.walkthrough/walkthrough.json`.

## Step 2: Reconnaissance

Map the repository:

- List the top-level directory structure.
- Identify languages, frameworks, and build systems.
- Locate entry points: `main` functions, HTTP handlers, CLI commands, event listeners.
- Note configuration files, dependency manifests, and test directories.

Summarize your findings before proceeding.

## Step 3: Analysis

### If scope is `flow`:

1. Start from the entry point (user-specified or discovered).
2. Trace the execution path step by step: handler → service → repository → external calls.
3. For each significant step, record:
   - The file and relevant code snippet (10–30 lines).
   - A clear explanation of what the code does and why.
   - How data transforms as it flows through.
4. Group steps into logical sections (e.g., "Request Parsing", "Authentication", "Business Logic", "Persistence").

### If scope is `overview`:

1. Identify the major architectural layers or modules.
2. For each layer, pick 1–3 representative code snippets.
3. Explain responsibilities, boundaries, and how layers communicate.
4. Group into sections by architectural boundary.

## Step 4: Diagrams

For each section, create at least one Mermaid diagram. Choose the most appropriate type:

- **sequence** — for request/response flows between components
- **flowchart** — for decision trees, branching logic, pipelines
- **classDiagram** — for type hierarchies and relationships
- **graph** — for dependency graphs, module relationships

Write valid Mermaid DSL. Keep diagrams focused — no more than 8–10 nodes/participants per diagram.

## Step 5: Assembly

Build the walkthrough JSON matching this schema:

```json
{
  "title": "string — descriptive title",
  "description": "string — one-paragraph summary",
  "scope": "flow | overview",
  "repository": "string — repo name or URL",
  "generatedAt": "ISO 8601 timestamp",
  "sections": [
    {
      "id": "kebab-case-id",
      "title": "Human-readable section title",
      "description": "What this section covers",
      "steps": [
        {
          "id": "section-id/step-number",
          "title": "Step title",
          "explanation": "Detailed markdown explanation (2-4 paragraphs)",
          "codeSnippet": {
            "filePath": "relative/path/to/file.ext",
            "language": "go | typescript | python | ...",
            "startLine": 42,
            "endLine": 68,
            "source": "// actual code content"
          },
          "diagram": {
            "type": "sequence | flowchart | classDiagram | graph",
            "mermaid": "sequenceDiagram\n    participant A\n    A->>B: call"
          }
        }
      ]
    }
  ]
}
```

### Rules:

- Every step MUST have `id`, `title`, and `explanation`.
- `codeSnippet` and `diagram` are optional per step but most steps should have at least one.
- `id` values must be unique across the entire walkthrough.
- Section IDs use kebab-case. Step IDs use the pattern `section-id/step-number`.
- `explanation` should be thorough: 2–4 paragraphs of markdown explaining the code's purpose, design decisions, and how it connects to adjacent steps.
- `source` in code snippets should contain the actual code, not placeholders.
- Keep code snippets focused (10–30 lines). Trim to the relevant portion.

## Step 6: Write Output

1. Create the output directory if needed (default: `.walkthrough/`).
2. Write the JSON to the output path.
3. Report:
   - Number of sections and steps generated.
   - Files referenced.
   - Suggested next steps: `walkthrough view .walkthrough/walkthrough.json`
