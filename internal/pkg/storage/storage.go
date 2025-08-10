package storage

import (
	"errors"
	"math"
	"sync"

	"go.uber.org/zap"
)

type Kind string

const (
	kindInt      Kind = "D"
	kindString   Kind = "S"
	kindUndefind Kind = "UND"
)

type Value struct {
	Val       any  `json:"val"`
	ValueType Kind `json:"valueType"`
}

func newValue(val any) (Value, error) {
	valueType := getType(val)
	if valueType != kindUndefind {
		return Value{
			Val:       val,
			ValueType: valueType,
		}, nil
	}
	return Value{}, errors.New("Undefined ValueType")
}

func getType(val any) Kind {
	switch val.(type) {
	case int:
		return kindInt
	case float64:
		if isFloatInt(val) {
			return kindInt
		}
		return kindUndefind
	case string:
		return kindString
	default:
		return kindUndefind
	}
}

func isFloatInt(num any) bool {
	return num.(float64) == math.Trunc(num.(float64))
}

type List struct {
	Elem []any `json:"elements"`
}

type Storage struct {
	inner  map[string]Value
	list   map[string]*List
	logger *zap.Logger
	mu     *sync.Mutex
}

func NewStorage() (Storage, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return Storage{}, err
	}
	defer logger.Sync()
	logger.Info("storage created")

	return Storage{
		inner:  make(map[string]Value),
		list:   make(map[string]*List),
		logger: logger,
		mu:     new(sync.Mutex),
	}, nil
}

func (r Storage) Set(key string, value any) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	var val Value

	val, err := newValue(value)
	if err != nil {
		return err
	}

	r.inner[key] = val
	r.logger.Info("value set",
		zap.String("key", key),
		zap.String("value", string(val.ValueType)))

	defer r.logger.Sync()

	return nil
}

func (r Storage) Get(key string) *any {
	r.mu.Lock()
	defer r.mu.Unlock()

	result, ok := r.inner[key]
	if !ok {
		return nil
	}

	r.logger.Info("value get",
		zap.String("key", key),
		zap.String("value", string(result.ValueType)))

	defer r.logger.Sync()
	return &result.Val
}

func (r Storage) GetType(key string) any {
	r.mu.Lock()
	defer r.mu.Unlock()
	result, ok := r.inner[key]
	if !ok {
		return "No"
	}

	return result.ValueType
}

func (s *Storage) LPUSH(key string, elements []any) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(elements) == 0 {
		return errors.New("WrongArgs")
	}

	if _, exists := s.list[key]; !exists {
		s.list[key] = &List{}
	}

	list := s.list[key]
	for i := len(elements) - 1; i >= 0; i-- {
		list.Elem = append([]any{elements[i]}, list.Elem...)
	}

	s.logger.Info("LPUSH executed")
	return nil
}

func (s *Storage) RPUSH(key string, elements []any) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(elements) == 0 {
		return errors.New("WrongArgs")
	}

	if _, exist := s.list[key]; !exist {
		s.list[key] = &List{}
	}

	list := s.list[key]
	for i := 0; i < len(elements); i++ {
		list.Elem = append(list.Elem, elements[i])
	}

	s.logger.Info("RPUSH executed")
	return nil
}

func (s *Storage) RADDTOSET(key string, elements []any) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	if len(elements) == 0 {
		return errors.New("WrongArgs")
	}

	if _, exist := s.list[key]; !exist {
		s.list[key] = &List{Elem: make([]any, 0)}
	}

	list := s.list[key]
	existing := make(map[any]bool)

	for _, elem := range list.Elem {
		existing[elem] = true
	}

	for _, elem := range elements {
		if !existing[elem] {
			list.Elem = append(list.Elem, elem)
			existing[elem] = true
		}
	}

	s.logger.Info("RADDTOSET executed")
	return nil
}

func (s *Storage) LPOP(key string, count ...int) ([]any, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	list, exist := s.list[key]
	if !exist || len(list.Elem) == 0 {
		return nil, errors.New("list is empty or does not exist")
	}

	if len(count) > 2 {
		return nil, errors.New("WrongArgs")
	}

	if len(count) == 0 {
		result := list.Elem[:1]
		list.Elem = list.Elem[1:]
		return result, nil
	}

	if len(count) == 1 {
		start := count[0]
		if start > len(list.Elem) {
			start = len(list.Elem)
		}

		result := list.Elem[:start]
		list.Elem = list.Elem[start:]

		return result, nil
	}

	start := count[0]
	end := count[1]
	start, end = convertIndex(start, len(list.Elem)), convertIndex(end, len(list.Elem))

	if start < 0 || end < 0 || start > end || start >= len(list.Elem) {
		return nil, errors.New("invalid index range")
	}

	if end >= len(list.Elem) {
		end = len(list.Elem) - 1
	}

	result := list.Elem[start : end+1]
	list.Elem = append(list.Elem[:start], list.Elem[end+1:]...)

	return result, nil

}

func convertIndex(index int, length int) int {
	if index < 0 {
		index = length + index
	}

	return index
}

func (s *Storage) RPOP(key string, count ...int) ([]any, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	list, exist := s.list[key]
	if !exist || len(list.Elem) == 0 {
		return nil, errors.New("list is empty or does not exist")
	}

	if len(count) == 0 {
		lastIdx := len(list.Elem) - 1
		popped := list.Elem[lastIdx]
		list.Elem = list.Elem[:lastIdx]
		return []any{popped}, nil
	}

	if len(count) == 1 {
		n := count[0]
		if n <= 0 {
			return nil, errors.New("count must be positive")
		}
		if n > len(list.Elem) {
			n = len(list.Elem)
		}
		startIdx := len(list.Elem) - n
		popped := list.Elem[startIdx:]
		reverse(popped)
		list.Elem = list.Elem[:startIdx]
		return popped, nil
	}

	if len(count) == 2 {
		start, end := count[0], count[1]

		start = normalizeIndex(start, len(list.Elem))
		end = normalizeIndex(end, len(list.Elem))

		if start < 0 || end < 0 || start >= len(list.Elem) || end >= len(list.Elem) {
			return nil, errors.New("invalid index range")
		}

		if start > end {
			start, end = end, start
		}

		popped := list.Elem[start : end+1]

		reverse(popped)
		list.Elem = append(list.Elem[:start], list.Elem[end+1:]...)
		return popped, nil
	}

	return nil, errors.New("wrong number of arguments")
}

func normalizeIndex(index, length int) int {
	if index < 0 {
		index = length + index
		if index < 0 {
			index = 0
		}
	}
	return index
}

func reverse(slice []any) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func (s *Storage) LSET(key string, index int, element any) (any, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	list, exist := s.list[key]
	if !exist {
		return "", errors.New("WrongArgs")
	}

	if index < 0 {
		index += len(list.Elem)
	}

	if index < 0 || index >= len(list.Elem) {
		return "", errors.New("index out of range")
	}

	list.Elem[index] = element
	s.logger.Info("LSET executed")

	return "OK", nil

}

func (s *Storage) LGET(key string, index int) (any, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	list, exist := s.list[key]
	if !exist {
		return "", errors.New("WrongArgs")
	}

	if index < 0 {
		index += len(list.Elem)
	}

	if index < 0 || index >= len(list.Elem) {
		return "", errors.New("index out of range")
	}

	res := list.Elem[index]
	s.logger.Info("LGET executed")

	return res, nil
}
