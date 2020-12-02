package util

import (
	"fmt"
	"sort"
)

// UpdateCurTick updates the current "safe" tick
func UpdateCurTick(ticks [][]int, noReplicas int, curTick int) (int, string) {
	if len(ticks[0]) == 0 {
		return curTick, ""
	}
	var str string

	/* Determine the latest timestamp per replica */
	b := make([]int, noReplicas)
	for i := 0; i < noReplicas; i++ {
		b[i] = max(ticks, i, curTick)
	}

	/* Determine the earliest timestamp for all replicas */
	res := min(b)
	if res == 0 {
		res = 1
	}
	for i := 0; i < noReplicas; i++ {
		if len(ticks[i]) > 0 {
			str = str + fmt.Sprintln(ticks[i])
		}
	}

	str = str + fmt.Sprintln(b) + "=======\n"
	str = str + ":" + fmt.Sprintln(res)
	return res, str
}

// determine the latest timestamp per replica
func max(a [][]int, i int, curTick int) int {
	sort.Ints(a[i])
	for j := 0; j < len(a[i])-1; j++ {
		if a[i][j+1] > curTick+1 && a[i][j+1] > a[i][j]+1 {
			return a[i][j]
		}
	}
	return a[i][len(a[i])-1]
}

// determine the earliest timestamp for all replicas
func min(b []int) int {
	res := b[0]
	for j := 1; j < len(b); j++ {
		if b[j] < res {
			res = b[j]
		}
	}
	return res
}
