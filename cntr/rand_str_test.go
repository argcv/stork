package cntr

import (
	"testing"

	"github.com/argcv/stork/assert"
)

func TestRandomString(t *testing.T) {
	length := 10
	randStr := RandomString(length)
	assert.ExpectEQ(t, length, len(randStr))
}

func TestRandomStringWithCharset(t *testing.T) {
	length := 10
	charset := "a"
	randStr := RandomStringWithCharset(length, charset)
	assert.ExpectEQ(t, length, len(randStr))
	for i := 0; i < length; i++ {
		assert.ExpectEQ(t, charset[0], randStr[i])
	}
	assert.ExpectEQ(t, "", RandomStringWithCharset(10, ""))
	assert.ExpectEQ(t, "", RandomStringWithCharset(-1, ""))
}
