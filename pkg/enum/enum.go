package enum

import "fmt"

type Enum[K comparable, V Enumable[K]] struct {
	cases map[K]V
}

type Enumable[K any] interface {
	Value() K
}

func New[K comparable, V Enumable[K]]() *Enum[K, V] {
	return &Enum[K, V]{
		cases: make(map[K]V), // SuperAdmin: ["create_user", "create_parking"]
	}
}

func (enum *Enum[K, V]) Register(enumCase V) V {
	enum.cases[enumCase.Value()] = enumCase
	return enumCase
}

func (enum *Enum[K, V]) From(candidate K) (V, bool) {
	c, ok := enum.cases[candidate]
	return c, ok
}

func (enum *Enum[K, V]) FromMany(candidates []K) ([]V, bool) {
	var cases = make([]V, 0, len(candidates))
	for _, candidate := range candidates {
		_case, ok := enum.From(candidate)
		if !ok {
			return nil, false
		}
		cases = append(cases, _case)
	}
	return cases, false
}

func (enum *Enum[K, V]) Exists(candidate K) bool {
	_, ok := enum.cases[candidate]
	return ok
}

func (enum *Enum[K, V]) ExistValue(val V) bool {
	return enum.Exists(val.Value())
}

func (enum *Enum[K, V]) ExistAll(candidates []K) bool {
	for _, candidate := range candidates {
		if !enum.Exists(candidate) {
			return false
		}
	}
	return true
}

func (enum *Enum[K, V]) Validate(val V) error {
	if !enum.ExistValue(val) {
		return fmt.Errorf("unknown enum value \"%s\"", val.Value())
	}
	return nil
}

func (enum *Enum[K, V]) ValidateAll(values []V) error {
	for _, val := range values {
		if err := enum.Validate(val); err != nil {
			return err
		}
	}
	return nil
}

func (enum *Enum[K, V]) Cases() []V {
	var cases = make([]V, 0, len(enum.cases))
	for _, value := range enum.cases {
		cases = append(cases, value)
	}
	return cases
}

type Reader[K comparable, V Enumable[K]] interface {
	From(candidate K) (V, bool)
	FromMany(candidates []K) ([]V, bool)
	Exists(candidate K) bool
	ExistAll(candidates []K) bool
	ExistValue(val V) bool
	Validate(val V) error
	ValidateAll(values []V) error
	Cases() []V
}

func (enum *Enum[K, V]) Keys() []K {
	keys := make([]K, 0, len(enum.cases))
	for k := range enum.cases {
		keys = append(keys, k)
	}

	return keys
}
