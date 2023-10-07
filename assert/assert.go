package assert

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

const (
	defaultSkip = 0
)

var (
	testMode = false
	testBuff *bytes.Buffer
)

// Note: Don't Use the test out of assert
func TestFunc(f func()) (ret string) {
	testMode = true
	testBuff = &bytes.Buffer{}
	defer func() {
		recover()
		testMode = false
	}()
	f()
	ret = testBuff.String()
	testBuff = nil
	return ret
}

func ExpectEQ(t testing.TB, expected, actual interface{}, msg ...string) {
	if ret := equal(defaultSkip+1, expected, actual, msg...); ret != "" {
		if testMode {
			testBuff.Write([]byte(ret))
		} else {
			fmt.Println(ret)
			t.Fail()
		}
	}
}

func CheckIsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func CheckIsNotNil(i interface{}) bool {
	return !CheckIsNil(i)
}

func ExpectNil(t testing.TB, actual interface{}, msg ...string) {
	if ret := test(defaultSkip+1,
		fmt.Sprintf("Expected: nil, actual:  '%v'", actual),
		func() bool {
			//fmt.Fprintf(os.Stderr, "actual: %v | is nil? %v", actual, CheckIsNil(actual))
			return CheckIsNil(actual)
		},
		msg...); ret != "" {
		if testMode {
			testBuff.Write([]byte(ret))
		} else {
			fmt.Println(ret)
			t.Fail()
		}
	}
}

func ExpectNotNil(t testing.TB, actual interface{}, msg ...string) {
	if ret := test(defaultSkip+1,
		fmt.Sprintf("Expected: not nil, actual:  '%v'", actual),
		func() bool {
			return CheckIsNotNil(actual)
		},
		msg...); ret != "" {
		if testMode {
			testBuff.Write([]byte(ret))
		} else {
			fmt.Println(ret)
			t.Fail()
		}
	}
}

// 期望长度
// TODO(yu): test cases
func ExpectLen(t testing.TB, length int, actual interface{}, msg ...string) {
	s := reflect.ValueOf(actual)

	for s.Kind() == reflect.Ptr {
		s = reflect.Indirect(s)
	}

	if !s.IsValid() {
		ret := test(defaultSkip+1,
			fmt.Sprintf("Expected: length = %v, invalid type: %v", length, s.Kind()),
			func() bool {
				return false
			},
			msg...)
		if testMode {
			testBuff.Write([]byte(ret))
		} else {
			fmt.Println(ret)
			t.Fail()
		}
	}

	switch s.Kind() {
	case reflect.Slice, reflect.Array:
		if ret := test(defaultSkip+1,
			fmt.Sprintf("Expected: length = %v, actual:  '%v'", length, s.Len()),
			func() bool {
				return s.Len() == length
			},
			msg...); ret != "" {
			if testMode {
				testBuff.Write([]byte(ret))
			} else {
				fmt.Println(ret)
				t.Fail()
			}
		}
	default:
		ret := test(defaultSkip+1,
			fmt.Sprintf("Expected: length = %v, unexpected type: %v", length, s.Kind()),
			func() bool {
				return false
			},
			msg...)
		if testMode {
			testBuff.Write([]byte(ret))
		} else {
			fmt.Println(ret)
			t.Fail()
		}
	}
}

// 期望出现 error
func ExpectErr(t testing.TB, err error, msg ...string) {
	if ret := test(defaultSkip+1,
		fmt.Sprintf("Expected: error, actual:  '%v'", err),
		func() bool {
			return err != nil
		},
		msg...); ret != "" {
		if testMode {
			testBuff.Write([]byte(ret))
		} else {
			fmt.Println(ret)
			t.Fail()
		}
	}
}

func ExpectNoErr(t testing.TB, err error, msg ...string) {
	if ret := test(defaultSkip+1,
		fmt.Sprintf("Expected: no error, actual:  '%v'", err),
		func() bool {
			return err == nil
		},
		msg...); ret != "" {
		if testMode {
			testBuff.Write([]byte(ret))
		} else {
			fmt.Println(ret)
			t.Fail()
		}
	}
}

func ExpectTrue(t testing.TB, condition bool, msg ...string) {
	if ret := equal(defaultSkip+1, true, condition, msg...); ret != "" {
		if testMode {
			testBuff.Write([]byte(ret))
		} else {
			fmt.Println(ret)
			t.Fail()
		}
	}
}

func ExpectFalse(t testing.TB, condition bool, msg ...string) {
	if ret := equal(defaultSkip+1, false, condition, msg...); ret != "" {
		if testMode {
			testBuff.Write([]byte(ret))
		} else {
			fmt.Println(ret)
			t.Fail()
		}
	}
}

func ExpectNE(t testing.TB, expected, actual interface{}, msg ...string) {
	if ret := notEqual(defaultSkip+1, expected, actual, msg...); ret != "" {
		if testMode {
			testBuff.Write([]byte(ret))
		} else {
			fmt.Println(ret)
			t.Fail()
		}
	}
}

// ExpectLT expects val1 less than val2
func ExpectLT(t testing.TB, val1, val2 interface{}, msg ...string) {
	if ret := compare(defaultSkip+1, val1, val2, "<", func() bool {
		if val1 == nil || val2 == nil {
			return false
		}
		//
		//s1, ok1 := val1.(sort.Interface)
		//s2, ok2 := val2.(sort.Interface)
		//if ok1 && ok2 {
		//	return s1.Less()
		//}
		// prime type
		v1 := reflect.ValueOf(val1)
		v2 := reflect.ValueOf(val2)
		if v1.Type() != v2.Type() {
			return false
		}
		return compareRealValues(val1, val2, v1.Kind(), compareLess)
	}, msg...); ret != "" {
		if testMode {
			testBuff.Write([]byte(ret))
		} else {
			fmt.Println(ret)
			t.Fail()
		}
	}
}

// ExpectGT expect val1 greater than val2
func ExpectGT(t testing.TB, val1, val2 interface{}, msg ...string) {
	if ret := compare(defaultSkip+1, val1, val2, ">", func() bool {
		if val1 == nil || val2 == nil {
			return false
		}
		// prime type
		v1 := reflect.ValueOf(val1)
		v2 := reflect.ValueOf(val2)
		if v1.Type() != v2.Type() {
			return false
		}
		return compareRealValues(val1, val2, v1.Kind(), compareGreater)
	}, msg...); ret != "" {
		if testMode {
			testBuff.Write([]byte(ret))
		} else {
			fmt.Println(ret)
			t.Fail()
		}
	}
}

// ExpectLE
func ExpectLE(t testing.TB, val1, val2 interface{}, msg ...string) {
	if ret := compare(defaultSkip+1, val1, val2, "<", func() bool {
		if val1 == nil || val2 == nil {
			return false
		}
		//
		//s1, ok1 := val1.(sort.Interface)
		//s2, ok2 := val2.(sort.Interface)
		//if ok1 && ok2 {
		//	return s1.Less()
		//}
		// prime type
		v1 := reflect.ValueOf(val1)
		v2 := reflect.ValueOf(val2)
		if v1.Type() != v2.Type() {
			return false
		}
		return compareRealValues(val1, val2, v1.Kind(), compareLess|compareEqual)
	}, msg...); ret != "" {
		if testMode {
			testBuff.Write([]byte(ret))
		} else {
			fmt.Println(ret)
			t.Fail()
		}
	}
}

// ExpectGE expect val1 greater than or equal to val2
func ExpectGE(t testing.TB, val1, val2 interface{}, msg ...string) {
	if ret := compare(defaultSkip+1, val1, val2, ">", func() bool {
		if val1 == nil || val2 == nil {
			return false
		}
		// prime type
		v1 := reflect.ValueOf(val1)
		v2 := reflect.ValueOf(val2)
		if v1.Type() != v2.Type() {
			return false
		}
		return compareRealValues(val1, val2, v1.Kind(), compareGreater|compareEqual)
	}, msg...); ret != "" {
		if testMode {
			testBuff.Write([]byte(ret))
		} else {
			fmt.Println(ret)
			t.Fail()
		}
	}
}

func compare(skip int, val1, val2 interface{}, operator string, cmp func() bool, msg ...string) string {
	return test(skip+1,
		fmt.Sprintf("Expected: '%v' %v '%v', actual:  '%v' (%v) vs. '%v' (%v)",
			val1, operator, val2,
			val1, reflect.TypeOf(val1), val2, reflect.TypeOf(val2)),
		cmp,
		msg...)
}

func equal(skip int, expected, actual interface{}, msg ...string) string {
	return test(skip+1,
		fmt.Sprintf("\tExpected: '%v'\n\tWhich is: %v\nTo be equal to: '%v'\n\tWhich is: %v",
			expected, reflect.TypeOf(expected), actual, reflect.TypeOf(actual)),
		func() bool {
			return reflect.DeepEqual(expected, actual)
		},
		msg...)
}

func notEqual(skip int, expected, actual interface{}, msg ...string) string {
	return test(skip+1,
		fmt.Sprintf("Expected: '%v' != '%v', actual:  '%v' (%v) vs. '%v' (%v)",
			expected, actual,
			expected, reflect.TypeOf(expected), actual, reflect.TypeOf(actual)),
		func() bool {
			return !reflect.DeepEqual(expected, actual)
		},
		msg...)
}

func test(skip int, concl string, cmp func() bool, msg ...string) string {
	if !cmp() {
		return fail(skip+1, "Failure:\n%s\nMessage: %s", concl, strings.Join(msg, " "))
	}
	return ""
}

func fail(skip int, format string, args ...interface{}) string {
	_, file, line, _ := runtime.Caller(skip + 1)
	return fmt.Sprintf("    %s:%d: %s\n", filepath.Base(file), line, fmt.Sprintf(format, args...))
}

func TestWrap(skip int, concl string, cmp func() bool, msg ...string) string {
	return test(skip+1, concl, cmp, msg...)
}
