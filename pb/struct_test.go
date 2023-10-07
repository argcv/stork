package pb

import (
	"testing"

	st "google.golang.org/protobuf/types/known/structpb"

	"github.com/argcv/stork/assert"
)

func TestToStruct(t *testing.T) {
	type cst struct {
		foo int
	}

	baseMap := map[string]interface{}{
		"nil":     nil,
		"bool":    true,
		"int":     1,
		"int8":    int8(8),
		"int16":   int8(16),
		"int32":   int8(32),
		"int64":   int8(64),
		"uint8":   uint8(8),
		"uint16":  uint8(16),
		"uint32":  uint8(32),
		"uint64":  uint8(64),
		"float":   2.01,
		"float32": float32(32.0),
		"float64": float32(64.0),
		"string":  "hello",
		"cst": cst{
			foo: 1,
		},
	}

	foo := map[string]interface{}{}
	for k, v := range baseMap {
		foo[k] = v
	}
	foo["sub"] = baseMap

	val := ToStruct(foo)
	_, ok := val.Fields["empty"]
	assert.ExpectFalse(t, ok, "no field empty")
	nilVal, ok := val.Fields["nil"]
	assert.ExpectTrue(t, ok, "found field nil")
	assert.ExpectEQ(t, st.NullValue_NULL_VALUE, nilVal.GetNullValue())
	boolVal, ok := val.Fields["bool"]
	assert.ExpectTrue(t, ok, "found field bool")
	assert.ExpectEQ(t, true, boolVal.GetBoolValue())
	intVal, ok := val.Fields["int"]
	assert.ExpectTrue(t, ok, "found field int")
	assert.ExpectEQ(t, 1.0, intVal.GetNumberValue())
	floatVal, ok := val.Fields["float"]
	assert.ExpectTrue(t, ok, "found field float")
	assert.ExpectEQ(t, 2.01, floatVal.GetNumberValue())
	stringVal, ok := val.Fields["string"]
	assert.ExpectTrue(t, ok, "found field string")
	assert.ExpectEQ(t, "hello", stringVal.GetStringValue())
	cstVal, ok := val.Fields["cst"]
	assert.ExpectTrue(t, ok, "found field cst")
	assert.ExpectNE(t, nil, cstVal.GetStructValue())
}
