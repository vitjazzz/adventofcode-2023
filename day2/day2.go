package day2

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"strconv"
	"strings"
)

var cubesCount = map[string]int{"red": 12, "green": 13, "blue": 14}

// 3280 is too high
// 3003 is too high

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/2/input")
	//taskLines := getTestLines()
	var desiredSum int
	for _, line := range taskLines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		power := getPower(line)
		desiredSum += power
		fmt.Printf("Adding power %d, new sum is %d\n", power, desiredSum)
	}
	fmt.Println(desiredSum)
}

func checkIfPossible(line string) bool {
	colonIndex := strings.Index(line, ":")
	gameLog := line[colonIndex+2:]
	for _, playLog := range strings.Split(gameLog, "; ") {
		for _, cubesStr := range strings.Split(playLog, ", ") {
			strs := strings.Split(cubesStr, " ")
			count, _ := strconv.Atoi(strs[0])
			colour := strs[1]
			if cubesCount[colour] < count {
				fmt.Printf("%s is not possible because of %s game\n", line, cubesStr)
				return false
			}
		}
	}
	return true
}

func getPower(line string) int {
	maxColour := map[string]int{"red": 0, "green": 0, "blue": 0}
	colonIndex := strings.Index(line, ":")
	gameLog := line[colonIndex+2:]
	for _, playLog := range strings.Split(gameLog, "; ") {
		for _, cubesStr := range strings.Split(playLog, ", ") {
			strs := strings.Split(cubesStr, " ")
			count, _ := strconv.Atoi(strs[0])
			colour := strs[1]
			if maxColour[colour] < count {
				maxColour[colour] = count
			}
		}
	}
	power := 1
	for _, count := range maxColour {
		power *= count
	}
	return power
}

func getGameId(line string) (gameId int) {
	colonIndex := strings.Index(line, ":")
	gameIdStr := line[5:colonIndex]
	gameId, _ = strconv.Atoi(gameIdStr)
	return gameId
}

func getTestLines() (taskLines []string) {
	test := "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green\nGame 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue\nGame 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red\nGame 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red\nGame 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green"
	return strings.Split(test, "\n")
}
