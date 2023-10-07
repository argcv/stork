package assert

import (
	"reflect"
)

// inspired by https://github.com/stretchr/testify/blob/master/assert/assertion_compare.go

type CompareType int

const (
	compareLess CompareType = 1 << (iota + 1)
	compareEqual
	compareGreater
)

var (
	intType   = reflect.TypeOf(int(1))
	int8Type  = reflect.TypeOf(int8(1))
	int16Type = reflect.TypeOf(int16(1))
	int32Type = reflect.TypeOf(int32(1))
	int64Type = reflect.TypeOf(int64(1))

	uintType   = reflect.TypeOf(uint(1))
	uint8Type  = reflect.TypeOf(uint8(1))
	uint16Type = reflect.TypeOf(uint16(1))
	uint32Type = reflect.TypeOf(uint32(1))
	uint64Type = reflect.TypeOf(uint64(1))

	float32Type = reflect.TypeOf(float32(1))
	float64Type = reflect.TypeOf(float64(1))

	stringType = reflect.TypeOf("")
)

func compareRealValues(obj1, obj2 interface{}, kind reflect.Kind, expected CompareType) bool {
	obj1Value := reflect.ValueOf(obj1)
	obj2Value := reflect.ValueOf(obj2)

	testIsExpected := func(value CompareType) bool {
		return expected&value > 0
	}

	// throughout this switch we try and avoid calling .Convert() if possible,
	// as this has a pretty big performance impact
	switch kind {
	case reflect.Int:
		{
			intobj1, ok := obj1.(int)
			if !ok {
				intobj1 = obj1Value.Convert(intType).Interface().(int)
			}
			intobj2, ok := obj2.(int)
			if !ok {
				intobj2 = obj2Value.Convert(intType).Interface().(int)
			}
			if intobj1 > intobj2 {
				return testIsExpected(compareGreater)
			}
			if intobj1 == intobj2 {
				return testIsExpected(compareEqual)
			}
			if intobj1 < intobj2 {
				return testIsExpected(compareLess)
			}
		}
	case reflect.Int8:
		{
			int8obj1, ok := obj1.(int8)
			if !ok {
				int8obj1 = obj1Value.Convert(int8Type).Interface().(int8)
			}
			int8obj2, ok := obj2.(int8)
			if !ok {
				int8obj2 = obj2Value.Convert(int8Type).Interface().(int8)
			}
			if int8obj1 > int8obj2 {
				return testIsExpected(compareGreater)
			}
			if int8obj1 == int8obj2 {
				return testIsExpected(compareEqual)
			}
			if int8obj1 < int8obj2 {
				return testIsExpected(compareLess)
			}
		}
	case reflect.Int16:
		{
			int16obj1, ok := obj1.(int16)
			if !ok {
				int16obj1 = obj1Value.Convert(int16Type).Interface().(int16)
			}
			int16obj2, ok := obj2.(int16)
			if !ok {
				int16obj2 = obj2Value.Convert(int16Type).Interface().(int16)
			}
			if int16obj1 > int16obj2 {
				return testIsExpected(compareGreater)
			}
			if int16obj1 == int16obj2 {
				return testIsExpected(compareEqual)
			}
			if int16obj1 < int16obj2 {
				return testIsExpected(compareLess)
			}
		}
	case reflect.Int32:
		{
			int32obj1, ok := obj1.(int32)
			if !ok {
				int32obj1 = obj1Value.Convert(int32Type).Interface().(int32)
			}
			int32obj2, ok := obj2.(int32)
			if !ok {
				int32obj2 = obj2Value.Convert(int32Type).Interface().(int32)
			}
			if int32obj1 > int32obj2 {
				return testIsExpected(compareGreater)
			}
			if int32obj1 == int32obj2 {
				return testIsExpected(compareEqual)
			}
			if int32obj1 < int32obj2 {
				return testIsExpected(compareLess)
			}
		}
	case reflect.Int64:
		{
			int64obj1, ok := obj1.(int64)
			if !ok {
				int64obj1 = obj1Value.Convert(int64Type).Interface().(int64)
			}
			int64obj2, ok := obj2.(int64)
			if !ok {
				int64obj2 = obj2Value.Convert(int64Type).Interface().(int64)
			}
			if int64obj1 > int64obj2 {
				return testIsExpected(compareGreater)
			}
			if int64obj1 == int64obj2 {
				return testIsExpected(compareEqual)
			}
			if int64obj1 < int64obj2 {
				return testIsExpected(compareLess)
			}
		}
	case reflect.Uint:
		{
			uintobj1, ok := obj1.(uint)
			if !ok {
				uintobj1 = obj1Value.Convert(uintType).Interface().(uint)
			}
			uintobj2, ok := obj2.(uint)
			if !ok {
				uintobj2 = obj2Value.Convert(uintType).Interface().(uint)
			}
			if uintobj1 > uintobj2 {
				return testIsExpected(compareGreater)
			}
			if uintobj1 == uintobj2 {
				return testIsExpected(compareEqual)
			}
			if uintobj1 < uintobj2 {
				return testIsExpected(compareLess)
			}
		}
	case reflect.Uint8:
		{
			uint8obj1, ok := obj1.(uint8)
			if !ok {
				uint8obj1 = obj1Value.Convert(uint8Type).Interface().(uint8)
			}
			uint8obj2, ok := obj2.(uint8)
			if !ok {
				uint8obj2 = obj2Value.Convert(uint8Type).Interface().(uint8)
			}
			if uint8obj1 > uint8obj2 {
				return testIsExpected(compareGreater)
			}
			if uint8obj1 == uint8obj2 {
				return testIsExpected(compareEqual)
			}
			if uint8obj1 < uint8obj2 {
				return testIsExpected(compareLess)
			}
		}
	case reflect.Uint16:
		{
			uint16obj1, ok := obj1.(uint16)
			if !ok {
				uint16obj1 = obj1Value.Convert(uint16Type).Interface().(uint16)
			}
			uint16obj2, ok := obj2.(uint16)
			if !ok {
				uint16obj2 = obj2Value.Convert(uint16Type).Interface().(uint16)
			}
			if uint16obj1 > uint16obj2 {
				return testIsExpected(compareGreater)
			}
			if uint16obj1 == uint16obj2 {
				return testIsExpected(compareEqual)
			}
			if uint16obj1 < uint16obj2 {
				return testIsExpected(compareLess)
			}
		}
	case reflect.Uint32:
		{
			uint32obj1, ok := obj1.(uint32)
			if !ok {
				uint32obj1 = obj1Value.Convert(uint32Type).Interface().(uint32)
			}
			uint32obj2, ok := obj2.(uint32)
			if !ok {
				uint32obj2 = obj2Value.Convert(uint32Type).Interface().(uint32)
			}
			if uint32obj1 > uint32obj2 {
				return testIsExpected(compareGreater)
			}
			if uint32obj1 == uint32obj2 {
				return testIsExpected(compareEqual)
			}
			if uint32obj1 < uint32obj2 {
				return testIsExpected(compareLess)
			}
		}
	case reflect.Uint64:
		{
			uint64obj1, ok := obj1.(uint64)
			if !ok {
				uint64obj1 = obj1Value.Convert(uint64Type).Interface().(uint64)
			}
			uint64obj2, ok := obj2.(uint64)
			if !ok {
				uint64obj2 = obj2Value.Convert(uint64Type).Interface().(uint64)
			}
			if uint64obj1 > uint64obj2 {
				return testIsExpected(compareGreater)
			}
			if uint64obj1 == uint64obj2 {
				return testIsExpected(compareEqual)
			}
			if uint64obj1 < uint64obj2 {
				return testIsExpected(compareLess)
			}
		}
	case reflect.Float32:
		{
			float32obj1, ok := obj1.(float32)
			if !ok {
				float32obj1 = obj1Value.Convert(float32Type).Interface().(float32)
			}
			float32obj2, ok := obj2.(float32)
			if !ok {
				float32obj2 = obj2Value.Convert(float32Type).Interface().(float32)
			}
			if float32obj1 > float32obj2 {
				return testIsExpected(compareGreater)
			}
			if float32obj1 == float32obj2 {
				return testIsExpected(compareEqual)
			}
			if float32obj1 < float32obj2 {
				return testIsExpected(compareLess)
			}
		}
	case reflect.Float64:
		{
			float64obj1, ok := obj1.(float64)
			if !ok {
				float64obj1 = obj1Value.Convert(float64Type).Interface().(float64)
			}
			float64obj2, ok := obj2.(float64)
			if !ok {
				float64obj2 = obj2Value.Convert(float64Type).Interface().(float64)
			}
			if float64obj1 > float64obj2 {
				return testIsExpected(compareGreater)
			}
			if float64obj1 == float64obj2 {
				return testIsExpected(compareEqual)
			}
			if float64obj1 < float64obj2 {
				return testIsExpected(compareLess)
			}
		}
	case reflect.String:
		{
			stringobj1, ok := obj1.(string)
			if !ok {
				stringobj1 = obj1Value.Convert(stringType).Interface().(string)
			}
			stringobj2, ok := obj2.(string)
			if !ok {
				stringobj2 = obj2Value.Convert(stringType).Interface().(string)
			}
			if stringobj1 > stringobj2 {
				return testIsExpected(compareGreater)
			}
			if stringobj1 == stringobj2 {
				return testIsExpected(compareEqual)
			}
			if stringobj1 < stringobj2 {
				return testIsExpected(compareLess)
			}
		}
	}
	return false
}
