package day9

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"strings"
)

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/9/input")
	//taskLines := getTestLines2()
	res := 0
	for _, line := range taskLines {
		if line == "" {
			continue
		}
		sequence := adventutils.ParseNumbers(line, " ")
		predictionMap := buildPredictionMap(sequence)
		//nextValue := predictNextValue(predictionMap)
		prevValue := predictPrevValue(predictionMap)
		fmt.Printf("for line %s prev value is %d\n", line, prevValue)
		res += prevValue
	}
	fmt.Printf("result = %d\n", res)
}

func predictNextValue(predictionMap map[int][]int) int {
	res := 0
	for i := len(predictionMap) - 1; i >= 0; i-- {
		currentSequence := predictionMap[i]
		res += currentSequence[len(currentSequence)-1]
	}
	return res
}

func predictPrevValue(predictionMap map[int][]int) int {
	res := 0
	for i := len(predictionMap) - 1; i >= 0; i-- {
		currentSequence := predictionMap[i]
		res = currentSequence[0] - res
	}
	return res
}

func buildPredictionMap(sequence []int) map[int][]int {
	res := make(map[int][]int)
	res[0] = sequence
	prevSequence := sequence
	for i := 1; !allZeroes(prevSequence); i++ {
		var currentSequence []int
		for j := 1; j < len(prevSequence); j++ {
			newValue := prevSequence[j] - prevSequence[j-1]
			currentSequence = append(currentSequence, newValue)
		}
		res[i] = currentSequence
		prevSequence = currentSequence
	}
	return res
}

func allZeroes(sequence []int) bool {
	for _, val := range sequence {
		if val != 0 {
			return false
		}
	}
	return true
}

func getTestLines() (taskLines []string) {
	test := "0 3 6 9 12 15\n1 3 6 10 15 21\n10 13 16 21 30 45"
	return strings.Split(test, "\n")
}
func getTestLines2() (taskLines []string) {
	test := "10  13  16  21  30  45"
	return strings.Split(test, "\n")
}
