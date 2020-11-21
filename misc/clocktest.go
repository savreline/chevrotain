package main

import (
	"fmt"

	"github.com/savreline/GoVector/govec/vclock"
)

func main() {
	map1 := map[string]uint64{
		"R1": 2,
		"R2": 2,
		"R3": 3,
	}
	map2 := map[string]uint64{
		"R1": 2,
		"R2": 2,
		"R3": 2,
	}
	clock := vclock.New()
	// fmt.Println("C", clock.CopyFromMap(map1).Compare(map2, vclock.Concurrent))
	// fmt.Println("A", clock.CopyFromMap(map1).Compare(map2, vclock.Ancestor))
	// fmt.Println("D", clock.CopyFromMap(map1).Compare(map2, vclock.Descendant))
	// fmt.Println("E", clock.CopyFromMap(map1).Compare(map2, vclock.Equal))

	fmt.Println(clock.CopyFromMap(map1).CompareClocks(map2))
}
