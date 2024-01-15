package day6

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"strconv"
	"strings"
)

const timeKey = "Time:"
const distanceKey = "Distance:"

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/6/input", true)
	//taskLines := getTestLines()
	games := getGamesAdvanced(taskLines)
	res := 1
	for _, g := range games {
		waysToBeat := getWaysToBeat(g)
		res *= waysToBeat
		fmt.Printf("game - %s, ways to beat - %d\n", g, waysToBeat)
	}
	fmt.Printf("res - %d\n", res)
}

func getWaysToBeat(game game) int {
	lowestWinningTime := getLowestWinningTime(game)
	highestWinningTime := game.time - lowestWinningTime
	return (highestWinningTime - lowestWinningTime) + 1
}

func getLowestWinningTime(game game) int {
	highestLosingTime := 0
	lowestWinningTime := game.time
	for highestLosingTime+1 < lowestWinningTime {
		time := (highestLosingTime + lowestWinningTime) / 2
		if isWinning(game, time) {
			lowestWinningTime = time
		} else {
			highestLosingTime = time
		}
	}
	return lowestWinningTime
}

func isWinning(game game, holdingTime int) bool {
	distance := (game.time - holdingTime) * holdingTime
	return distance > game.distance
}

func getGamesAdvanced(lines []string) []game {
	var time int
	var distance int
	for _, line := range lines {
		if strings.Contains(line, timeKey) {
			timeStr := strings.ReplaceAll(line[len(timeKey)+1:], " ", "")
			time, _ = strconv.Atoi(timeStr)
		}
		if strings.Contains(line, distanceKey) {
			distanceStr := strings.ReplaceAll(line[len(distanceKey)+1:], " ", "")
			distance, _ = strconv.Atoi(distanceStr)
		}
	}
	return []game{{time: time, distance: distance}}
}

func getGames(lines []string) (games []game) {
	var times []int
	var distances []int
	for _, line := range lines {
		if strings.Contains(line, timeKey) {
			times = adventutils.ParseNumbers(line[len(timeKey)+1:], " ")
		}
		if strings.Contains(line, distanceKey) {
			distances = adventutils.ParseNumbers(line[len(distanceKey)+1:], " ")
		}
	}
	for i, _ := range times {
		games = append(games, game{time: times[i], distance: distances[i]})
	}
	return games
}

type game struct {
	time, distance int
}

func getTestLines() (taskLines []string) {
	test := "Time:      7  15   30\nDistance:  9  40  200"
	return strings.Split(test, "\n")
}
