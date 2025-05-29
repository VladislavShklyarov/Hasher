package logic

import (
	"business-service/gen"
	"testing"
)

func TestFindAliveVariables(t *testing.T) {
	tests := []struct {
		name       string
		operations []*gen.Operation
		wantAlive  map[string]bool
		wantGraph  map[string][]string
	}{
		{
			name: "basic case with dependencies",
			operations: []*gen.Operation{
				{Type: "calc", Op: "+", Var: "x", Left: "10", Right: "2"},
				{Type: "print", Var: "x"},
				{Type: "calc", Op: "-", Var: "y", Left: "x", Right: "3"},
				{Type: "calc", Op: "*", Var: "z", Left: "x", Right: "y"},
				{Type: "print", Var: "w"},
				{Type: "print", Var: "w1"},
				{Type: "print", Var: "w2"},
				{Type: "calc", Op: "*", Var: "w", Left: "z", Right: "y"},
				{Type: "calc", Op: "*", Var: "w1", Left: "w", Right: "x"},
				{Type: "calc", Op: "*", Var: "w2", Left: "z", Right: "3"},
			},
			wantAlive: map[string]bool{
				"x":  true,
				"y":  true,
				"z":  true,
				"w":  true,
				"w1": true,
				"w2": true,
			},
			wantGraph: map[string][]string{
				"x":  {},
				"y":  {"x"},
				"z":  {"x", "y"},
				"w":  {"z", "y"},
				"w1": {"w", "x"},
				"w2": {"z"},
			},
		},
		{
			name: "no print operations",
			operations: []*gen.Operation{
				{Type: "calc", Op: "+", Var: "a", Left: "1", Right: "2"},
				{Type: "calc", Op: "*", Var: "b", Left: "a", Right: "3"},
			},
			wantAlive: map[string]bool{},
			wantGraph: map[string][]string{
				"a": {},
				"b": {"a"},
			},
		},
		{
			name: "print variable with no dependencies",
			operations: []*gen.Operation{
				{Type: "print", Var: "c"},
			},
			wantAlive: map[string]bool{
				"c": true,
			},
			wantGraph: map[string][]string{},
		},
		{
			name: "cyclic dependency",
			operations: []*gen.Operation{
				{Type: "calc", Op: "+", Var: "p", Left: "q", Right: "1"},
				{Type: "calc", Op: "*", Var: "q", Left: "p", Right: "2"},
				{Type: "print", Var: "p"},
			},
			wantAlive: map[string]bool{
				"p": true,
				"q": true,
			},
			wantGraph: map[string][]string{
				"p": {"q"},
				"q": {"p"},
			},
		},
		{
			name: "mixed number and variable dependencies",
			operations: []*gen.Operation{
				{Type: "calc", Op: "-", Var: "m", Left: "10", Right: "n"},
				{Type: "calc", Op: "+", Var: "n", Left: "5", Right: "3"},
				{Type: "print", Var: "m"},
			},
			wantAlive: map[string]bool{
				"m": true,
				"n": true,
			},
			wantGraph: map[string][]string{
				"m": {"n"},
				"n": {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAlive, gotGraph := FindAliveVariables(tt.operations)

			// Проверка alive переменных
			if len(gotAlive) != len(tt.wantAlive) {
				t.Errorf("unexpected number of alive variables: got %d, want %d", len(gotAlive), len(tt.wantAlive))
			}
			for key := range tt.wantAlive {
				if !gotAlive[key] {
					t.Errorf("expected variable %q to be alive, but it was not", key)
				}
			}

			// Проверка графа
			for key, wantDeps := range tt.wantGraph {
				gotDeps := gotGraph[key]
				if len(gotDeps) != len(wantDeps) {
					t.Errorf("for variable %q: expected %v, got %v", key, wantDeps, gotDeps)
					continue
				}
				// простая проверка по значениям
				for i, wantDep := range wantDeps {
					if gotDeps[i] != wantDep {
						t.Errorf("for variable %q: at index %d expected %q, got %q", key, i, wantDep, gotDeps[i])
					}
				}
			}
		})
	}
}

func TestProcess(t *testing.T) {
	tests := []struct {
		name       string
		operations []*gen.Operation
		required   map[string]bool
		wantResult []*gen.VariableValue
		wantBroken []string
	}{
		{
			name: "basic_calc_and_print",
			operations: []*gen.Operation{
				// calc a = 1 + 2
				{
					Type:  "calc",
					Var:   "a",
					Left:  "1",
					Right: "2",
					Op:    "+",
				},
				// print a
				{
					Type: "print",
					Var:  "a",
				},
			},
			required: map[string]bool{"a": true},
			wantResult: []*gen.VariableValue{
				{Var: "a", Value: 3},
			},
			wantBroken: []string{},
		},
		{
			name: "print_missing_variable",
			operations: []*gen.Operation{
				{
					Type: "print",
					Var:  "missing",
				},
			},
			required: map[string]bool{},
			wantResult: []*gen.VariableValue{
				{Var: "missing", Value: 0},
			},
			wantBroken: []string{"missing"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotBroken := Process(tt.operations, tt.required)

			// Проверяем количество результатов
			if len(gotResult) != len(tt.wantResult) {
				t.Fatalf("expected %d results, got %d", len(tt.wantResult), len(gotResult))
			}
			// Проверяем каждый элемент results
			for i := range gotResult {
				if gotResult[i].Var != tt.wantResult[i].Var || gotResult[i].Value != tt.wantResult[i].Value {
					t.Errorf("result[%d] = %+v, want %+v", i, gotResult[i], tt.wantResult[i])
				}
			}

			// Проверяем количество сломанных переменных
			if len(gotBroken) != len(tt.wantBroken) {
				t.Fatalf("expected %d brokenVars, got %d", len(tt.wantBroken), len(gotBroken))
			}
			// Проверяем содержимое brokenVars (без учёта порядка)
			gotMap := make(map[string]bool)
			for _, v := range gotBroken {
				gotMap[v] = true
			}
			for _, wantV := range tt.wantBroken {
				if !gotMap[wantV] {
					t.Errorf("brokenVars missing %q", wantV)
				}
			}
		})
	}
}
