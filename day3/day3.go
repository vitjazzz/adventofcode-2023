package day3

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/3/input")
	//taskLines := getTestLines()
	taskArray := toArray(taskLines)
	totalSum := calculateGearsSum(taskArray)
	fmt.Printf("Total sum is %d", totalSum)
}

func calculateGearsSum(taskArray [][]rune) int {
	numbersArray := computeNumbersArray(taskArray)
	totalSum := calculateSum(taskArray, numbersArray)

	return totalSum
}

func calculateSum(taskArray [][]rune, numbersArray [][]int) int {
	totalSum := 0
	for i, line := range taskArray {
		for j, char := range line {
			if char != '*' {
				continue
			}
			var adjacentNumbers []int
			for iTmp := adventutils.Max(0, i-1); iTmp <= adventutils.Min(len(numbersArray)-1, i+1); iTmp++ {
				for jTmp := adventutils.Max(0, j-1); jTmp <= adventutils.Min(len(line)-1, j+1); jTmp++ {
					adjacentNumber := numbersArray[iTmp][jTmp]
					if adjacentNumber == 0 {
						continue
					}
					if jTmp == adventutils.Max(0, j-1) || numbersArray[iTmp][jTmp-1] == 0 {
						adjacentNumbers = append(adjacentNumbers, adjacentNumber)
					}
				}
			}
			if len(adjacentNumbers) == 2 {
				totalSum += adjacentNumbers[0] * adjacentNumbers[1]
			}
		}
	}
	return totalSum
}

func computeNumbersArray(taskArray [][]rune) [][]int {
	numbersArray := make([][]int, len(taskArray))
	for i, line := range taskArray {
		numbersArray[i] = make([]int, len(line))
		currentNumberStr := ""
		for j, char := range line {
			if !unicode.IsDigit(char) {
				continue
			}
			currentNumberStr += string(char)
			if currentNumberStr != "" && isEndOfNumber(j, line) {
				currentNumber, _ := strconv.Atoi(currentNumberStr)
				for jTmp := j + 1 - len(currentNumberStr); jTmp <= j; jTmp++ {
					numbersArray[i][jTmp] = currentNumber
				}
				currentNumberStr = ""
			}
		}
	}
	return numbersArray
}

func isEndOfNumber(j int, line []rune) bool {
	return j+1 >= len(line) || !unicode.IsDigit(line[j+1])
}

func toArray(lines []string) [][]rune {
	lines = filterEmptyLines(lines)
	result := make([][]rune, len(lines))
	for i, line := range lines {
		result[i] = []rune(line)
	}
	return result
}

func filterEmptyLines(lines []string) []string {
	var filteredLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			filteredLines = append(filteredLines, line)
		}
	}
	return filteredLines
}

func getTestLines() (taskLines []string) {
	test := "467..114..\n...*......\n..35..633.\n......#...\n617*......\n.....+.58.\n..592.....\n......755.\n...$.*....\n.664.598.."
	return strings.Split(test, "\n")
}
