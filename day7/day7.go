package day7

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"
)

var cardNominations = map[string]int{"A": 12, "K": 11, "Q": 10, "T": 9, "9": 8, "8": 7, "7": 6, "6": 5, "5": 4, "4": 3, "3": 2, "2": 1, "J": 0}

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/7/input")
	//taskLines := getTestLines()
	hands := getHands(taskLines)
	sortedHands := hands[:]
	sort.Slice(sortedHands, func(i, j int) bool {
		return sortedHands[i].power < sortedHands[j].power
	})
	fmt.Printf("sortedHands - %v\n", sortedHands)
	res := 0
	for i, h := range sortedHands {
		winning := h.bid * (i + 1)
		res += winning
		fmt.Printf("hand - %s, winning - %d\n", h.cards, winning)
	}
	fmt.Printf("result - %d\n", res)
}

func getHands(lines []string) (hands []hand) {
	for _, line := range lines {
		if line == "" {
			continue
		}
		values := strings.Split(line, " ")
		cards := values[0]
		power := getPower(cards)
		bid, _ := strconv.Atoi(values[1])
		hands = append(hands, hand{cards, power, bid})
	}
	return hands
}

func getPower(cards string) int {
	cardStrength := getCardStrength(cards)
	combinationName, combinationRank := getCombination(cards)
	combinationStrength := combinationRank * powInt(len(cardNominations), len(cards))
	fmt.Printf("Cards %s have combinationName - %s, combinationStrength - %d, cardStrength - %d\n",
		cards, combinationName, combinationStrength, cardStrength)
	return combinationStrength + cardStrength
}

func getCombination(cards string) (name string, rank int) {
	cardsArray := strings.Split(cards, "")
	uniqueCardsMap := make(map[string]int)
	jokers := 0
	for _, card := range cardsArray {
		if card == "J" {
			jokers++
		} else {
			uniqueCardsMap[card]++
		}
	}
	uniqueCardsMap = useJokers(cardsArray, uniqueCardsMap, jokers)
	uniqueCards := len(uniqueCardsMap)
	if uniqueCards == 1 {
		return "five of a kind", 6
	} else if uniqueCards == 2 && mapContainsValue(uniqueCardsMap, 4) {
		return "four of a kind", 5
	} else if uniqueCards == 2 && mapContainsValue(uniqueCardsMap, 3) {
		return "full house", 4
	} else if uniqueCards == 3 && mapContainsValue(uniqueCardsMap, 3) {
		return "three of a kind", 3
	} else if uniqueCards == 3 && mapContainsValue(uniqueCardsMap, 2) {
		return "two pair", 2
	} else if uniqueCards == 4 && mapContainsValue(uniqueCardsMap, 2) {
		return "one pair", 1
	} else {
		return "high card", 0
	}
}

func useJokers(cardsArray []string, uniqueCardsMap map[string]int, jokers int) map[string]int {
	cardToDuplicate := cardsArray[0]
	for card, count := range uniqueCardsMap {
		if uniqueCardsMap[cardToDuplicate] < count {
			cardToDuplicate = card
		}
	}
	uniqueCardsMap[cardToDuplicate] += jokers
	return uniqueCardsMap
}

func getCardStrength(cards string) int {
	reversedCards := strings.Split(cards, "")
	slices.Reverse(reversedCards)
	res := 0
	for i, card := range reversedCards {
		cardPower := cardNominations[card] * powInt(len(cardNominations), i)
		res += cardPower
	}
	return res
}

type hand struct {
	cards      string
	power, bid int
}

func powInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func mapContainsValue(targetMap map[string]int, targetValue int) bool {
	for _, val := range targetMap {
		if val == targetValue {
			return true
		}
	}
	return false
}

func getTestLines() (taskLines []string) {
	test := "32T3K 765\nT55J5 684\nKK677 28\nKTJJT 220\nQQQJA 483"
	return strings.Split(test, "\n")
}
