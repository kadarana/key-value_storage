package storage

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"go.uber.org/zap"
)

type Value struct {
	stringValue string
	intValue    int
	valueType   string
}

type List struct {
	elem []string
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
		mu:     *make(sync.Mutex),
	}, nil
}

func (r Storage) Set(key, value string) {
	valueType := ckeckType(value)

	var val Value
	if valueType == "D" {
		intValue, _ := strconv.Atoi(value)
		val = Value{
			stringValue: value,
			intValue:    intValue,
			valueType:   valueType,
		}
	} else {
		val = Value{
			stringValue: value,
			valueType:   valueType,
		}
	}

	r.inner[key] = val
	r.logger.Info("value set", zap.String("key", key), zap.String("value", val.stringValue))
	defer r.logger.Sync()
}

func (r Storage) Get(key string) *string {
	result, ok := r.inner[key]
	if !ok {
		return nil
	}
	r.logger.Info("value get", zap.String("key", key), zap.String("value", result.stringValue))
	defer r.logger.Sync()
	return &result.stringValue
}

func (r Storage) GetType(key string) string {
	result, ok := r.inner[key]
	if !ok {
		return "No"
	}

	return result.valueType
}

func ckeckType(value string) string {
	_, err := strconv.Atoi(value)
	if err != nil {
		return "S"
	} else {
		return "D"
	}
}

func (s *Storage) LPUSH(key string, elements []string) error {

	if len(elements) == 0 {
		return errors.New("WrongArgs")
	}

	if _, exists := s.list[key]; !exists {
		s.list[key] = &List{}
	}

	list := s.list[key]
	for i := len(elements) - 1; i >= 0; i-- {
		list.elem = append([]string{elements[i]}, list.elem...)
	}

	s.logger.Info("LPUSH executed")
	return nil
}

func (s *Storage) RPUSH(key string, elements []string) error {
	if len(elements) == 0 {
		return errors.New("WrongArgs")
	}

	if _, exist := s.list[key]; !exist {
		s.list[key] = &List{}
	}

	list := s.list[key]
	for i := 0; i < len(elements); i++ {
		list.elem = append(list.elem, elements[i])
	}

	s.logger.Info("RPUSH executed")
	return nil
}

func (s *Storage) RADDTOSET(key string, elements []string) error {

	if len(elements) == 0 {
		return errors.New("WrongArgs")
	}

	if _, exist := s.list[key]; !exist {
		s.list[key] = &List{}
	}

	list := s.list[key]
	k := 0
	for i := 0; i < len(elements); i++ {

		res := strings.Contains(strings.Join(list.elem, ""), elements[i])

		if !res {
			list.elem = append(list.elem, elements[k])
			k++
		}
	}

	s.logger.Info("RADDTOSET executed")
	return nil
}

func (s *Storage) LPOP(key string, count ...int) ([]string, error) {

	list, exist := s.list[key]
	if !exist || len(list.elem) == 0 {
		return nil, errors.New("list is empty or does not exist")
	}

	if len(count) == 0 || len(count) < 2 {
		return nil, errors.New("WrongArgs")
	}

	if len(count) == 1 {
		start := count[0]
		if start > len(list.elem) {
			start = len(list.elem)
		}

		result := list.elem[:start]
		list.elem = list.elem[start:]

		return result, nil
	}

	start := count[0]
	end := count[1]
	start, end = convertIndex(start, len(list.elem)), convertIndex(end, len(list.elem))

	if start < 0 || end < 0 || start > end || start >= len(list.elem) {
		return nil, errors.New("invalid index range")
	}

	if end >= len(list.elem) {
		end = len(list.elem) - 1
	}

	result := list.elem[start : end+1]
	list.elem = append(list.elem[:start], list.elem[end+1:]...)

	return result, nil

}

func convertIndex(index int, length int) int {
	if index < 0 {
		index = length + index
	}

	return index
}

func (s *Storage) RPOP(key string, count ...int) ([]string, error) {
	list, exist := s.list[key]
	if !exist || len(list.elem) == 0 {
		return nil, errors.New("list is empty or does not exist")
	}

	if len(count) == 0 {
		lastIdx := len(list.elem) - 1
		popped := list.elem[lastIdx]
		list.elem = list.elem[:lastIdx]
		return []string{popped}, nil
	}

	if len(count) == 1 {
		n := count[0]
		if n <= 0 {
			return nil, errors.New("count must be positive")
		}
		if n > len(list.elem) {
			n = len(list.elem)
		}
		startIdx := len(list.elem) - n
		popped := list.elem[startIdx:]
		reverse(popped)
		list.elem = list.elem[:startIdx]
		return popped, nil
	}

	if len(count) == 2 {
		start, end := count[0], count[1]

		start = normalizeIndex(start, len(list.elem))
		end = normalizeIndex(end, len(list.elem))

		if start < 0 || end < 0 || start >= len(list.elem) || end >= len(list.elem) {
			return nil, errors.New("invalid index range")
		}

		if start > end {
			start, end = end, start
		}

		popped := list.elem[start : end+1]

		reverse(popped)
		list.elem = append(list.elem[:start], list.elem[end+1:]...)
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

func reverse(slice []string) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func (s *Storage) LSET(key string, index int, element string) (string, error) {
	list, exist := s.list[key]
	if !exist {
		return "", errors.New("WrongArgs")
	}

	if index < 0 {
		index += len(list.elem)
	}

	if index < 0 || index >= len(list.elem) {
		return "", errors.New("index out of range")
	}

	list.elem[index] = element
	s.logger.Info("LSET executed")

	return "OK", nil

}

func (s *Storage) LGET(key string, index int) (string, error) {
	list, exist := s.list[key]
	if !exist {
		return "", errors.New("WrongArgs")
	}

	if index < 0 {
		index += len(list.elem)
	}

	if index < 0 || index >= len(list.elem) {
		return "", errors.New("index out of range")
	}

	res := list.elem[index]
	s.logger.Info("LGET executed")

	return res, nil
}
