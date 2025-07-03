package storage

import (
	"strconv"

	"go.uber.org/zap"
)

type Value struct {
	stringValue string
	intValue    int
	valueType   string
}

type Storage struct {
	inner  map[string]Value
	logger *zap.Logger
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
		logger: logger,
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
