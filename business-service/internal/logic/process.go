package logic

import (
	"business-service/gen"
	"fmt"
	"strconv"
	"time"
)

func Process(operations []*gen.Operation, required map[string]bool, graph map[string][]string) ([]*gen.VariableValue, []string) {

	vars := NewVarStore()
	var result []*gen.VariableValue

	brokenVars := make([]string, 0, 10)

	remaining := make([]*gen.Operation, len(operations))
	copy(remaining, operations)

	for len(remaining) > 0 {
		progress := false
		next := []*gen.Operation{}

		for _, op := range remaining {
			switch op.GetType() {
			case "calc":
				if required[op.GetVar()] {
					if DoCalc(vars, op.GetVar(), op.GetLeft(), op.GetRight(), op.GetOp()) {
						progress = true
					} else {
						next = append(next, op)
					}
				}
			case "print":
				DoPrint(vars, &result, op.GetVar(), &brokenVars)
				continue
			}
		}

		if !progress {
			break
		}
		remaining = next
	}

	for k, v := range graph {
		fmt.Printf("%s -> %v\n", k, v)
	}
	fmt.Println("\n==========================\n")
	for _, variable := range result {
		fmt.Println(variable)
	}
	fmt.Println("\n==========================\n")
	return result, brokenVars

}

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func FindAliveVariables(operations []*gen.Operation) (map[string]bool, map[string][]string) {
	graph := map[string][]string{}

	// Построим граф зависимостей
	for _, op := range operations {
		if op.GetType() == "calc" {
			if !isNumber(op.GetLeft()) {
				graph[op.GetVar()] = append(graph[op.GetVar()], op.GetLeft())
			}
			if !isNumber(op.GetRight()) {
				graph[op.GetVar()] = append(graph[op.GetVar()], op.GetRight())
			}
		}
	}

	required := map[string]bool{}
	queue := []string{}

	for _, op := range operations {
		if op.GetType() == "print" {
			varName := op.GetVar()
			required[varName] = true
			queue = append(queue, varName)
		}
	}

	fmt.Println("Очередь:", queue)

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, dep := range graph[curr] {
			if !required[dep] {
				required[dep] = true
				queue = append(queue, dep)
			}
		}
	}

	return required, graph
}

func DoCalc(vars *VarStore, variable, left, right, op string) bool {
	time.Sleep(100 * time.Millisecond) // симуляция задержки

	if _, ok := vars.Get(variable); ok {
		return false
	}

	leftVal, err1 := parseOperand(left, vars)
	rightVal, err2 := parseOperand(right, vars)
	if err1 != nil || err2 != nil {
		return false
	}

	var result int
	switch op {
	case "+":
		result = leftVal + rightVal
	case "-":
		result = leftVal - rightVal
	case "*":
		result = leftVal * rightVal
	default:
		return false
	}

	vars.Set(variable, result)
	return true
}

func parseOperand(op string, vars *VarStore) (int, error) {
	val, err := strconv.Atoi(op)
	if err == nil {
		return val, nil
	}

	if v, ok := vars.Get(op); ok {
		return v, nil
	}
	return 0, fmt.Errorf("неизвестная переменная: %s", op)
}

func DoPrint(vars *VarStore, results *[]*gen.VariableValue, variable string, brokenVars *[]string) {
	if val, ok := vars.Get(variable); ok {
		*results = append(*results, &gen.VariableValue{
			Var:   variable,
			Value: int64(val),
		})
	} else {
		*results = append(*results, &gen.VariableValue{
			Var:   variable,
			Value: 0,
		})
		*brokenVars = append(*brokenVars, variable)
	}
}
