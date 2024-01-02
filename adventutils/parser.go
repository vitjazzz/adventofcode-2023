package adventutils

import (
	"strconv"
	"strings"
)

func ParseNumbers(numbersStr, separator string) (numbers []int) {
	for _, numberStr := range strings.Split(numbersStr, separator) {
		if numberStr == "" {
			continue
		}
		number, _ := strconv.Atoi(numberStr)
		numbers = append(numbers, number)
	}
	return numbers
}
