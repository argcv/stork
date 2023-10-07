package cntr

// GetEditDistanceT returns the edit distance between two strings.
// The edit distance is the number of operations required to transform
// one string into the other, where an operation is defined as an
// insertion, deletion, or substitution of a single character, or a
// swap of two adjacent characters.
// The function will return -1 if the distance is greater than the threshold.
func GetEditDistanceT(r1, r2 string, threshold int) int {
	abs := func(a int) int {
		if a < 0 {
			return -a
		}
		return a
	}

	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	sig := func(a, b rune) int {
		if a == b {
			return 0
		}
		return 1
	}

	if threshold <= 0 {
		if r1 == r2 {
			return 0
		} else {
			return -1
		}
	}

	s1 := []rune(r1)
	s2 := []rune(r2)

	if abs(len(s1)-len(s2)) > threshold {
		return -1
	}

	d := make([][]int, len(s1)+1)
	for i := range d {
		d[i] = make([]int, len(s2)+1)
	}

	for i := 1; i <= len(s1); i++ {
		st := max(1, i-threshold)
		en := min(len(s2), i+threshold)
		flag := true
		for j := st; j <= en; j++ {
			d[i][j] = threshold + 1
			if j-i+1 <= threshold && d[i-1][j]+1 < d[i][j] {
				d[i][j] = d[i-1][j] + 1
			}
			if i-j+1 <= threshold && d[i][j-1]+1 < d[i][j] {
				d[i][j] = d[i][j-1] + 1
			}
			d[i][j] = min(d[i][j], d[i-1][j-1]+sig(s1[i-1], s2[j-1]))
			if d[i][j] <= threshold {
				flag = false
			}
		}
		if flag && i > threshold {
			return -1
		}
	}

	if d[len(s1)][len(s2)] > threshold {
		return -1
	}
	return d[len(s1)][len(s2)]
}
