package common

import (
	"errors"
	"fmt"
	"strconv"
)

type ParseMapHelper struct{}

const (
	TypeString    = "string"
	TypeInt       = "int"
	TypeFloat64   = "float64"
	TypeBool      = "bool"
	TypeMapString = "map[string]"
	TypeSlice     = "[]"
)

func (nvh *ParseMapHelper) GetTypedNestedValue(data map[string]interface{}, expectedType string, keys ...string) (interface{}, error) {
	var current interface{} = data

	// 遍历 keys 逐层解析
	for i, key := range keys {
		switch typedCurrent := current.(type) {
		case map[string]interface{}:
			// 当前是 map[string]interface{}，查找下一层
			if value, exists := typedCurrent[key]; exists {
				current = value
			} else {
				return nil, fmt.Errorf("key not found: %s", key)
			}
		case []interface{}:
			// 当前是数组，尝试将 key 转换为索引访问
			index, err := nvh.ParseIndex(key)
			if err != nil {
				return nil, fmt.Errorf("key '%s' is not interface valid index in array", key)
			}
			if index >= 0 && index < len(typedCurrent) {
				current = typedCurrent[index]
			} else {
				return nil, fmt.Errorf("index out of bounds: %d", index)
			}
		default:
			return nil, fmt.Errorf("unexpected type %T at key: %s", current, key)
		}

		// 如果已经到达最后一个 key，检查类型
		if i == len(keys)-1 {
			return nvh.CheckType(current, expectedType)
		}
	}

	return nil, errors.New("unexpected error during key traversal")
}

// 检查类型是否匹配并返回值
func (nvh *ParseMapHelper) CheckType(value interface{}, expectedType string) (interface{}, error) {
	switch expectedType {
	case TypeString:
		if strValue, ok := value.(string); ok {
			return strValue, nil
		}
	case TypeInt:
		if intValue, ok := value.(int); ok {
			return intValue, nil
		}
	case TypeFloat64:
		if floatValue, ok := value.(float64); ok {
			return floatValue, nil
		}
	case TypeBool:
		if boolValue, ok := value.(bool); ok {
			return boolValue, nil
		}
	case TypeMapString:
		if mapValue, ok := value.(map[string]interface{}); ok {
			return mapValue, nil
		}
	case TypeSlice:
		if sliceValue, ok := value.([]interface{}); ok {
			return sliceValue, nil
		}
	default:
		return nil, fmt.Errorf("unsupported expected type: %s", expectedType)
	}
	return nil, fmt.Errorf("value found but type mismatch, expected: %s, got: %T", expectedType, value)
}

// 尝试将 key 转换为索引
func (nvh *ParseMapHelper) ParseIndex(key string) (int, error) {
	return strconv.Atoi(key)
}
