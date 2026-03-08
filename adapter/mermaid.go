package adapter

import (
	"fmt"
	"strings"

	"github.com/4thel00z/code-walkthrough/domain"
)

// MermaidRenderer implements port.DiagramRenderer
type MermaidRenderer struct{}

func NewMermaidRenderer() *MermaidRenderer {
	return &MermaidRenderer{}
}

func (r *MermaidRenderer) Render(d domain.Diagram, width int) (string, error) {
	source := strings.TrimSpace(d.Mermaid)
	if source == "" {
		return "", nil
	}

	switch d.Type {
	case domain.DiagramSequence:
		return r.renderSequence(source, width), nil
	case domain.DiagramFlowchart:
		return r.renderFlowchart(source, width), nil
	default:
		return r.renderRaw(source, width), nil
	}
}

func (r *MermaidRenderer) renderSequence(source string, width int) string {
	lines := strings.Split(source, "\n")
	var participants []string
	var messages []struct{ from, to, label string }

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "sequenceDiagram" || line == "" {
			continue
		}
		for _, sep := range []string{"->>", "-->>", "->", "-->"} {
			if idx := strings.Index(line, sep); idx >= 0 {
				from := strings.TrimSpace(line[:idx])
				rest := line[idx+len(sep):]
				parts := strings.SplitN(rest, ":", 2)
				to := strings.TrimSpace(parts[0])
				label := ""
				if len(parts) > 1 {
					label = strings.TrimSpace(parts[1])
				}
				messages = append(messages, struct{ from, to, label string }{from, to, label})
				addUnique(&participants, from)
				addUnique(&participants, to)
				break
			}
		}
	}

	if len(participants) == 0 {
		return r.renderRaw(source, width)
	}

	var sb strings.Builder
	// Header
	sb.WriteString("  ")
	for i, p := range participants {
		if i > 0 {
			sb.WriteString("          ")
		}
		sb.WriteString(fmt.Sprintf("%-10s", p))
	}
	sb.WriteString("\n")

	// Separator
	sb.WriteString("  ")
	for i := range participants {
		if i > 0 {
			sb.WriteString("          ")
		}
		sb.WriteString("    |     ")
	}
	sb.WriteString("\n")

	// Messages
	for _, msg := range messages {
		fromIdx := indexOf(participants, msg.from)
		toIdx := indexOf(participants, msg.to)
		sb.WriteString("  ")
		for i := range participants {
			if i == fromIdx && fromIdx < toIdx {
				arrow := fmt.Sprintf("──▶ %s", msg.label)
				sb.WriteString(fmt.Sprintf("    |%s", padTo(arrow, 10*(toIdx-fromIdx)+5)))
			} else if i == toIdx && fromIdx < toIdx {
				sb.WriteString("|")
			} else if i == fromIdx && fromIdx > toIdx {
				sb.WriteString("|")
			} else if i == toIdx && fromIdx > toIdx {
				arrow := fmt.Sprintf("◀── %s", msg.label)
				sb.WriteString(fmt.Sprintf("    |%s", padTo(arrow, 10*(fromIdx-toIdx)+5)))
			} else if i > min(fromIdx, toIdx) && i < max(fromIdx, toIdx) {
				continue
			} else {
				sb.WriteString("    |     ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (r *MermaidRenderer) renderFlowchart(source string, width int) string {
	lines := strings.Split(source, "\n")
	var edges []struct{ from, to string }

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "flowchart") || strings.HasPrefix(line, "graph") || line == "" {
			continue
		}
		for _, sep := range []string{"-->", "---"} {
			if idx := strings.Index(line, sep); idx >= 0 {
				from := extractNodeLabel(strings.TrimSpace(line[:idx]))
				to := extractNodeLabel(strings.TrimSpace(line[idx+len(sep):]))
				edges = append(edges, struct{ from, to string }{from, to})
				break
			}
		}
	}

	if len(edges) == 0 {
		return r.renderRaw(source, width)
	}

	var sb strings.Builder
	for i, edge := range edges {
		fromBox := fmt.Sprintf("[ %s ]", edge.from)
		toBox := fmt.Sprintf("[ %s ]", edge.to)
		if i == 0 {
			sb.WriteString(fmt.Sprintf("  %s\n", fromBox))
		}
		sb.WriteString("      |\n")
		sb.WriteString("      v\n")
		sb.WriteString(fmt.Sprintf("  %s\n", toBox))
	}

	return sb.String()
}

func (r *MermaidRenderer) renderRaw(source string, width int) string {
	return source
}

func extractNodeLabel(s string) string {
	if idx := strings.Index(s, "["); idx >= 0 {
		end := strings.Index(s, "]")
		if end > idx {
			return s[idx+1 : end]
		}
	}
	if idx := strings.Index(s, "("); idx >= 0 {
		end := strings.Index(s, ")")
		if end > idx {
			return s[idx+1 : end]
		}
	}
	return s
}

func addUnique(slice *[]string, val string) {
	for _, v := range *slice {
		if v == val {
			return
		}
	}
	*slice = append(*slice, val)
}

func indexOf(slice []string, val string) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}

func padTo(s string, length int) string {
	if len(s) >= length {
		return s
	}
	return s + strings.Repeat(" ", length-len(s))
}
