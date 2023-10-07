package cntr

import (
	"testing"
)

func TestDistinctStrings(t *testing.T) {
	k := []string{"aa", "aa", "bb", "aa"}
	k = DistinctStrings(k...)
	if len(k) != 2 {
		t.Errorf("Distinct Failed!!!")
	}
}

func TestDistinctStringsCaseInsensitive(t *testing.T) {
	k := []string{"aa", "AA", "bb", "aa"}
	k = DistinctStringsCaseInsensitive(k...)
	if len(k) != 2 {
		t.Errorf("Distinct Failed!!!")
	}
}
