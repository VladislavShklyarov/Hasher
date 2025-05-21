package logic

import (
	"business-service/gen"
	"fmt"
	"os"
	"strings"
)

// Проверка: участвует ли переменная как аргумент в других выражениях
func isUsedAsInput(varName string, ops []*gen.Operation) bool {
	for _, op := range ops {
		if op.Type == "calc" && (op.Left == varName || op.Right == varName) {
			return true
		}
	}
	return false
}

func ExportToDOT(ops []*gen.Operation, alive map[string]bool, graph map[string][]string, filename string) error {
	var sb strings.Builder
	sb.WriteString("digraph G {\n")
	sb.WriteString("  rankdir=LR;\n")
	sb.WriteString("  node [shape=box, style=filled];\n")

	// Узлы (только нечисловые)
	for _, op := range ops {
		color := "lightgrey"
		label := op.Var

		if op.Type == "print" {
			color = "lightgreen"
			label += "\\n[PRINT]"
		} else if alive[op.Var] {
			color = "lightblue"
		} else {
			color = "mistyrose"
		}

		sb.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\", fillcolor=%s];\n", op.Var, label, color))
	}

	// Рёбра (по уже готовому графу)
	for from, deps := range graph {
		for _, to := range deps {
			sb.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\";\n", to, from))
		}
	}

	sb.WriteString("}\n")
	return os.WriteFile(filename, []byte(sb.String()), 0644)
}
