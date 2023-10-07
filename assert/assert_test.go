package assert_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/argcv/stork/assert"
)

func TestExpectEQ(t *testing.T) {
	assert.ExpectEQ(t, "aaa", "aaa", "message for test equal")
	buff := assert.TestFunc(func() {
		assert.ExpectEQ(t, "aaa", "bbb", "message for test equal")
	})
	assert.ExpectTrue(t, strings.Contains(buff, "Failure:"))
}

func TestExpectNE(t *testing.T) {
	assert.ExpectNE(t, "aaa", 1, "message for test not equal")
	buff := assert.TestFunc(func() {
		assert.ExpectNE(t, "aaa", "aaa", "message for test not equal")
	})
	assert.ExpectTrue(t, strings.Contains(buff, "Failure:"))
}

func TestExpectLE(t *testing.T) {
	assert.ExpectLE(t, "aaa", "aab", "message for test less equal")
	buff := assert.TestFunc(func() {
		assert.ExpectLE(t, "aab", "aaa", "message for test less equal")
	})
	assert.ExpectTrue(t, strings.Contains(buff, "Failure:"))
}

func TestExpectNil(t *testing.T) {
	assert.ExpectNil(t, nil, "message for test not nil")
	buff := assert.TestFunc(func() {
		value := 1
		assert.ExpectNil(t, &value, "message for test not nil")
	})
	assert.ExpectTrue(t, strings.Contains(buff, "Failure:"))
}

func TestExpectNotNil(t *testing.T) {
	value := 1
	assert.ExpectNotNil(t, &value, "message for test is nil")
	buff := assert.TestFunc(func() {
		assert.ExpectNotNil(t, nil, "message for test is nil")
	})
	assert.ExpectTrue(t, strings.Contains(buff, "Failure:"))
}

func TestExpectTrue(t *testing.T) {
	assert.ExpectTrue(t, true)
	buff := assert.TestFunc(func() {
		assert.ExpectTrue(t, false)
	})
	assert.ExpectTrue(t, strings.Contains(buff, "Failure:"))
}

func TestExpectFalse(t *testing.T) {
	assert.ExpectFalse(t, false)
	buff := assert.TestFunc(func() {
		assert.ExpectFalse(t, true)
	})
	assert.ExpectTrue(t, strings.Contains(buff, "Failure:"))
}

func TestExpectLT(t *testing.T) {
	assert.ExpectLT(t, "aaa", "aab")
	buff := assert.TestFunc(func() {
		assert.ExpectLT(t, "aab", "aaa")
	})
	assert.ExpectTrue(t, strings.Contains(buff, "Failure:"))
}

func TestExpectErr(t *testing.T) {
	assert.ExpectErr(t, errors.New("dummy error"))
	buff := assert.TestFunc(func() {
		assert.ExpectErr(t, nil)
	})
	assert.ExpectTrue(t, strings.Contains(buff, "Failure:"))
}

func TestExpectNoErr(t *testing.T) {
	assert.ExpectNoErr(t, nil)
	buff := assert.TestFunc(func() {
		assert.ExpectNoErr(t, errors.New("dummy error"))
	})
	assert.ExpectTrue(t, strings.Contains(buff, "Failure:"))
}

func TestTestWrap(t *testing.T) {
	buff := assert.TestWrap(1, "con", func() bool {
		return true
	})
	t.Logf("buff ok: %v", buff)
	assert.ExpectEQ(t, "", buff)
	buff = assert.TestWrap(1, "con", func() bool { return false }, "it may failed")
	t.Logf("buff error: %v", buff)
	assert.ExpectTrue(t, strings.Contains(buff, "Failure:"))
}
