package day15

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"strings"
)

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/15/input", false)
	//taskLines := getTestLines()
	steps := getSteps(taskLines[0])
	res := 0
	for _, step := range steps {
		hashValue := hash(step)
		//fmt.Printf("step = %v, hashValue = %d\n", step, hashValue)
		res += hashValue
	}
	fmt.Printf("res = %d\n", res)
}

func hash(str string) int {
	res := 0
	symbols := []byte(str)
	for _, symbol := range symbols {
		res += int(symbol)
		res *= 17
		res %= 256
	}
	return res
}

func getSteps(line string) []string {
	line = strings.ReplaceAll(line, "\n", "")
	return strings.Split(line, ",")
}

func getTestLines() (taskLines []string) {
	test := "rn=1,cm-,qp=3,cm=2,qp-,pc=4,ot=9,ab=5,pc-,pc=6,ot=7"
	return strings.Split(test, "\n")
}
