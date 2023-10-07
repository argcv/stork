package cntr

import (
	"encoding/json"
	"testing"

	"github.com/argcv/stork/assert"
)

func TestM_CleanZeroBytes(t *testing.T) {
	rawStr := `{
	"options": {
		"f1": "abc\u0000def",
		"f2": {
			"f3": "abc\u0000def"
		},
		"f4": null,
		"f5": {},
		"f6": 1.0
	},
	"k1": "abc\u0000def"
}`
	rawM := M{}
	err := json.Unmarshal([]byte(rawStr), &rawM)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("rawM: %+v", rawM)
	rawM.CleanZeroBytes()
	t.Logf("cleanedM: %+v", rawM)
	assert.ExpectEQ(t, rawM.GetString("k1"), "abc def")
	assert.ExpectEQ(t, rawM.GetM("options").GetString("f1"), "abc def")
	assert.ExpectEQ(t, rawM.GetM("options").GetM("f2").GetString("f3"), "abc def")
}
