package logic

import (
	"business-service/gen"
	"container/list"
	"fmt"
	"strconv"
	"sync"
	"time"
)

func Process(operations []*gen.Operation, required map[string]bool) ([]*gen.VariableValue, []string) {
	vars := NewVarStore()
	var result []*gen.VariableValue
	brokenVars := make([]string, 0, 10)

	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	pending := append([]*gen.Operation{}, operations...)

	for {

		var (
			progress  bool
			remaining []*gen.Operation
			readyOps  []*gen.Operation
		)

		remaining, readyOps = devideOperations(pending, required, vars)

		if len(readyOps) == 0 {
			break
		}

		for _, op := range readyOps {
			wg.Add(1)
			go func(op *gen.Operation) {
				defer wg.Done()
				if doCalc(vars, op.GetVar(), op.GetLeft(), op.GetRight(), op.GetOp()) {
					mu.Lock()
					progress = true // Дает возможность добавить доп. условия
					mu.Unlock()
				}
			}(op)
		}

		wg.Wait()

		if !progress {
			break
		}

		pending = remaining
	}

	processPrint(vars, &result, operations, brokenVars)
	fmt.Println(result)
	return result, brokenVars
}

func processPrint(vars *VarStore, result *[]*gen.VariableValue, operations []*gen.Operation, brokenVars []string) {
	for _, op := range operations {
		if op.GetType() == "print" {
			doPrint(vars, result, op.GetVar(), &brokenVars)
		}
	}
}

func devideOperations(pending []*gen.Operation, required map[string]bool, vars *VarStore) (remaining []*gen.Operation, readyOps []*gen.Operation) {
	for _, op := range pending {
		if op.GetType() != "calc" || !required[op.GetVar()] {
			continue
		}

		if doCalcReady(vars, op.GetLeft(), op.GetRight()) {
			readyOps = append(readyOps, op)
		} else {
			remaining = append(remaining, op)
		}
	}
	return remaining, readyOps
}

func FindAliveVariables(operations []*gen.Operation) (map[string]bool, map[string][]string) {
	graph := map[string][]string{}
	required := map[string]bool{}
	queue := list.New()

	for _, op := range operations {
		fmt.Println(op)
		// Сначала строим граф зависимостей. Если в расчете один из элементов - переменная, то var зависит от нее
		if op.GetType() == "calc" {
			if !isNumber(op.GetLeft()) {
				graph[op.GetVar()] = append(graph[op.GetVar()], op.GetLeft())
			}
			if !isNumber(op.GetRight()) {
				graph[op.GetVar()] = append(graph[op.GetVar()], op.GetRight())
			}
		} else if op.GetType() == "print" {
			varName := op.GetVar()
			required[varName] = true
			queue.PushBack(varName)
			fmt.Println("		Добавили в очередь:", varName)
		}
	}

	fmt.Println("Итоговая Очередь:", queue)
	fmt.Println("Начинаем обход графа...")

	for queue.Len() > 0 {
		curr := queue.Front().Value.(string)
		queue.Remove(queue.Front())
		for _, dep := range graph[curr] {
			fmt.Printf("	Зависимость %s", dep)
			if !required[dep] {
				required[dep] = true
				queue.PushBack(dep)
			}
		}
	}

	fmt.Println(graph)
	return required, graph
}

func doCalc(vars *VarStore, variable, left, right, op string) bool {
	time.Sleep(50 * time.Millisecond) // симуляция задержки

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

func doPrint(vars *VarStore, results *[]*gen.VariableValue, variable string, brokenVars *[]string) {
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

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func doCalcReady(vars *VarStore, left, right string) bool {
	isReady := func(s string) bool {
		if isNumber(s) {
			return true
		}
		_, ok := vars.Get(s)
		return ok
	}

	return isReady(left) && isReady(right)
}
