package day1

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var digitMappings = map[string]int{"one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9}

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/1/input", true)
	//taskLines := getTestLinesAdvanced()
	var desiredSum int
	for _, line := range taskLines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fixedLine := fixLine(line)
		number := getDesiredNumber(fixedLine)
		desiredSum += number
	}
	fmt.Println(desiredSum)
}

func fixLine(line string) string {
	for key, value := range digitMappings {
		line = strings.ReplaceAll(line, key, key+strconv.Itoa(value)+key)
	}
	return line
}

func getDesiredNumber(line string) (desiredNumber int) {
	var numbers []int
	for _, char := range line {
		if unicode.IsDigit(char) {
			numbers = append(numbers, int(char-'0'))
		}
	}
	resultStr := strconv.Itoa(numbers[0]) + strconv.Itoa(numbers[len(numbers)-1])
	result, _ := strconv.Atoi(resultStr)
	return result
}

func getTestLines() (taskLines []string) {
	test := "1abc2\npqr3stu8vwx\na1b2c3d4e5f\ntreb7uchet"
	return strings.Split(test, "\n")
}

func getTestLinesAdvanced() (taskLines []string) {
	test := "two1nine\neightwothree\nabcone2threexyz\nxtwone3four\n4nineeightseven2\nzoneight234\n7pqrstsixteen"
	return strings.Split(test, "\n")
}
