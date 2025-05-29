# Бизнес-логика: обработка операций

# Общая идея

Исходная структура входных данных, а также рекомендованые требования к работе программы определяют некоторые сложности в реализации алгоритма:

1. Поля `left` и `right` могут содержать как число в строковом представленнии ("1", "94" и тд.) так и имя переменной. 
2. Некоторые переменные могут быть "фиктивными": они не выводятся на печать и не используются при других расчетах.

Алгоритм должен учитывать эти ограничения: нужно корреткно парсить входящие значения и производить рассчеты таким образом, чтобы не тратить лишнее время на фиктивные операции.

---

## Хранение
Прежде, чем приступать к непосредственному описанию логики работы алгоритма, разберемся с тем, как мы будем хранить данные.

Для хранения переменных используется потокобезопасная структура `VarStore`, основанная на `sync.RWMutex`. Это позволяет:

- Безопасно читать/записывать значения переменных из разных горутин
- Исключить повторную запись уже рассчитанных переменных

```go
type VarStore struct {
	data map[string]int
	mu   sync.RWMutex
}
```

Создается структура при помощи конструктора:
```go
func NewVarStore() *VarStore {
	return &VarStore{
		data: make(map[string]int, 20),
	}
}
```
Значение capacity задано условно, может быть изменено на большее (во избежания излишней аллокации памяти при переполнениии мапы)

У структуры есть два метода: getter и setter. Они позволяют потокобезопасно доставать и складывать значения в мапку, используя `sync.RWMutex`.

```go
func (s *VarStore) Set(name string, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[name] = value
}

func (s *VarStore) Get(name string) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[name]
	return val, ok
}
```

## Алгоритм
В этом разделе описан принцип работы алгоритма и функций, которые его составляют
### `parseOperand()` 

Вспомогательная функция `parseOperand()` предназначена для валидации входящих значений и решения проблемы под пунктом 1.

```go
func parseOperand(op string, vars *VarStore) (int, error) {
	val, err := strconv.Atoi(op)
	if err == nil { // Если ошибки нет, значит операнд - число.
		return val, nil
	}

	if v, ok := vars.Get(op); ok {
		return v, nil
	}
	return 0, fmt.Errorf("неизвестная переменная: %s", op)
}
```
Принцип работы: если при конвертации строки в число появляется ошибка: значит идем в мапку и берем оттуда его числовое значение. Если нет - конвертируем и возвращаем число. 

### `DoCalc()`

Функция DoCalc() составляет основу работы программы: именно она занимается всеми вычислениями.

```go
func DoCalc(vars *VarStore, variable, left, right, op string) bool {
	time.Sleep(1 * time.Second)

	if _, ok := vars.Get(variable); ok {
		return false
	}

	leftVal, err1 := parseOperand(left, vars)
	rightVal, err2 := parseOperand(right, vars)

	if err1 != nil || err2 != nil {
		return false // невозможно выполнить (еще нет нужных переменных)
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
```

Интересные аспекты:
- Вызывает функцию `parseOperand()` для получения значений о левой и правой части математической операции.
- Ставит временную задержку, эмулируя выполнение ресурсоемкой операции
- Вызывает методы `get` и `set` структуры `VarStore` для чтения значения и записи результата вычислений.
- Не имеет собственного возвращаемого значения (поскольку мы работаем с мапой напрямую), лишь булевый флаг: с его помощью функции более высокого порядка получают информация о возможности/невозможности осуществелния вычисления на данном этапе цикла. В успешном сценарии возвращает `true`

### `ProcessLocal()`

`ProcessLocal()` это основная исполняемая функция, внутри которой реализован весь алгоритм.

```go

func ProcessLocal(operations []*gen.Operation) {
	vars := NewVarStore()

	requiredVars := make(map[string]bool, 20)
	for _, op := range operations {
		switch op.GetType() {
		case "print":
			requiredVars[op.GetVar()] = true
		case "calc":
			{
				if _, err := strconv.Atoi(op.GetLeft()); err != nil {
					requiredVars[op.GetLeft()] = true
				}
				if _, err := strconv.Atoi(op.GetRight()); err != nil {
					requiredVars[op.GetRight()] = true
				}
			}
		default:
			fmt.Println("Wrong operation: calc or print expected")
		}
	}

	remaining := make([]*gen.Operation, len(operations))
	copy(remaining, operations)

	for len(remaining) > 0 {
		progress := false
		next := []*gen.Operation{}
		for _, op := range remaining {
			switch op.GetType() {
			case "calc":
				if requiredVars[op.GetVar()] {
					if DoCalc(vars, op.GetVar(), op.GetLeft(), op.GetRight(), op.GetOp()) {
						progress = true
					} else {
						next = append(next, op)
					}
				}
			case "print":
				DoPrint(vars, op.GetVar())
			default:
				fmt.Println("Wrong operation: calc or print expected")
			}
		}

		if !progress {
			break
		}
		remaining = next
	}

}
```
