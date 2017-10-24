// package lamportserver

// import (
// 	"fmt"
// )

// const (
// 	numSkiers = 40000
// 	numDays   = 2
// 	numLifts  = 40
// 	maxInt    = 1 << 31
// )

// // this function is not called in recent git commit.
// func mapSkierToDaysToLiftID(stats []*skierStat) ([][]int, [][]int) {
// 	var skierLiftTable = make([][]int, numSkiers+1)
// 	var skierVerticalTable = make([][]int, numSkiers+1)

// 	// allocate skierTable
// 	for i := range skierLiftTable {
// 		skierLiftTable[i] = make([]int, numDays+1)
// 		skierVerticalTable[i] = make([]int, numDays+1)
// 	}

// 	for _, stat := range stats {
// 		skierLiftTable[stat.skierID][stat.dayNum] += 1
// 		skierVerticalTable[stat.skierID][stat.dayNum] += liftIDToVertical(stat.liftID)
// 	}
// 	fmt.Println(skierLiftTable)
// 	fmt.Println(skierVerticalTable)
// 	return skierLiftTable, skierVerticalTable
// }
