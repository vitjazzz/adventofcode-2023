package day4

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"strconv"
	"strings"
)

// 42364 is too high

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/4/input")
	//taskLines := getTestLines()
	matchedNumbersMap, maxCardId := buildMatchedNumbersMap(taskLines)
	fmt.Printf("matchedNumbersMap - %v, maxCardId - %d\n", matchedNumbersMap, maxCardId)
	wonCardsMap := calculateWonCardsMap(matchedNumbersMap, maxCardId)
	fmt.Printf("wonCardsMap - %v\n", wonCardsMap)
	totalSum := 0
	for cardId, wonCards := range wonCardsMap {
		fmt.Printf("Card %d has %d won cards\n", cardId, wonCards)
		totalSum += wonCards
	}
	fmt.Printf("Total sum is %d", totalSum)
}

func buildMatchedNumbersMap(lines []string) (matchedNumbersMap map[int]int, maxCardId int) {
	matchedNumbersMap = make(map[int]int)
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		cardId := getCardId(line)
		userNumbers := getUserNumbers(line)
		winningNumbers := getWinningNumbers(line)
		matchedNumbers := getMatchedNumbers(userNumbers, winningNumbers)
		matchedNumbersMap[cardId] = len(matchedNumbers)
		maxCardId = cardId
	}
	return matchedNumbersMap, maxCardId
}

func calculateWonCardsMap(matchedNumbersMap map[int]int, maxCardId int) (wonCardsMap map[int]int) {
	wonCardsMap = make(map[int]int)
	for cardId, _ := range matchedNumbersMap {
		wonCardsMap[cardId] = 1
	}
	for currentCardId := 1; currentCardId <= maxCardId; currentCardId++ {
		wonCopies := matchedNumbersMap[currentCardId]
		for cardId := currentCardId + 1; cardId <= currentCardId+wonCopies; cardId++ {
			wonCardsMap[cardId] += wonCardsMap[currentCardId]
		}
	}
	return wonCardsMap
}

func getWinningNumbers(line string) (winningNumbers []int) {
	colonIndex := strings.Index(line, ":")
	dividerIndex := strings.Index(line, "|")
	winningNumbersStr := strings.TrimSpace(line[colonIndex+1 : dividerIndex])
	return adventutils.ParseNumbers(winningNumbersStr, " ")
}

func getUserNumbers(line string) (userNumbers []int) {
	dividerIndex := strings.Index(line, "|")
	userNumbersStr := strings.TrimSpace(line[dividerIndex+1:])
	return adventutils.ParseNumbers(userNumbersStr, " ")
}

func getCardId(line string) (cardId int) {
	colonIndex := strings.Index(line, ":")
	cardIdStr := strings.TrimSpace(line[5:colonIndex])
	cardId, _ = strconv.Atoi(cardIdStr)
	return cardId
}

func getMatchedNumbers(userNumbers, winningNumbers []int) []int {
	userNumbersMap := toMap(userNumbers)
	var matchedNumbers []int
	for _, winningNumber := range winningNumbers {
		if userNumbersMap[winningNumber] {
			matchedNumbers = append(matchedNumbers, winningNumber)
		}
	}
	return matchedNumbers
}

func calculatePoints(matchedNumbers []int) int {
	if len(matchedNumbers) == 0 {
		return 0
	}
	points := 1
	for i := 1; i < len(matchedNumbers); i++ {
		points *= 2
	}
	return points
}

func toMap(arr []int) map[int]bool {
	result := make(map[int]bool)
	for _, val := range arr {
		result[val] = true
	}
	return result
}

func getTestLines() (taskLines []string) {
	test := "Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53\nCard 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19\nCard 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1\nCard 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83\nCard 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36\nCard 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11"
	return strings.Split(test, "\n")
}
