package viperextra

import (
	"reflect"

	"go.uber.org/zap"
)

func ZapAtomicLevelDecodeHookFunc(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	var level zap.AtomicLevel
	if from.Kind() != reflect.String || to != reflect.TypeOf(level) {
		return data, nil
	}
	if err := level.UnmarshalText([]byte(data.(string))); err != nil {
		return nil, err
	}
	return level, nil
}
