package day5

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"math"
	"strings"
)

const seedsKey = "seeds: "
const mapKey = " map:"

var mappingOrder = []string{"seed", "soil", "fertilizer", "water", "light", "temperature", "humidity", "location"}

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/5/input")
	//taskLines := getTestLines()
	seeds := getSeeds(taskLines)
	allMappingRules := getMappingRules(taskLines)
	//fmt.Printf("seeds - %v, allMappingRules - %v\n", seeds, allMappingRules)
	locations := calculateLocations(seeds, allMappingRules)
	closestLocation := math.MaxInt
	for _, location := range locations {
		if closestLocation > location {
			closestLocation = location
		}
	}
	fmt.Printf("Closest location - %d\n", closestLocation)
}

func calculateLocations(seeds []int, rules []mappingRules) (locations []int) {
	for _, seed := range seeds {
		locations = append(locations, calculateLocation(seed, rules))
	}
	return locations
}

func calculateLocation(seed int, rules []mappingRules) int {
	currentSource := 0
	currentDestination := seed
	for i := 0; i < len(mappingOrder)-1; i++ {
		currentSource = currentDestination
		from := mappingOrder[i]
		to := mappingOrder[i+1]
		rule := getMappingRule(from, to, rules)
		currentDestination = rule.Map(currentSource)
	}
	return currentDestination
}

func getMappingRule(source, destination string, rules []mappingRules) mappingRules {
	for _, rule := range rules {
		if rule.source == source && rule.destination == destination {
			return rule
		}
	}
	fmt.Printf("Failed to find rule for %s source and %s destination!!!\n", source, destination)
	return mappingRules{}
}

func getSeeds(lines []string) []int {
	for _, line := range lines {
		if strings.Contains(line, seedsKey) {
			seedsStr := line[len(seedsKey):]
			return adventutils.ParseNumbers(seedsStr, " ")
		}
	}
	return []int{}
}

func getMappingRules(lines []string) (allMappingRules []mappingRules) {
	var currentRules mappingRules
	for i, line := range lines {
		if strings.Contains(line, seedsKey) {
			continue
		}
		if strings.Contains(line, mapKey) {
			mapKeyIndex := strings.Index(line, mapKey)
			sourceToDest := strings.Split(line[:mapKeyIndex], "-to-")
			currentRules = mappingRules{source: sourceToDest[0], destination: sourceToDest[1]}
			continue
		}
		if line != "" {
			numbers := adventutils.ParseNumbers(line, " ")
			rule := mappingRule{destinationStart: numbers[0], sourceStart: numbers[1], mappingRange: numbers[2]}
			currentRules.rules = append(currentRules.rules, rule)
		}
		if i+1 >= len(lines) || lines[i+1] == "" {
			allMappingRules = append(allMappingRules, currentRules)
		}
	}
	return allMappingRules
}

type mappingRules struct {
	source, destination string
	rules               []mappingRule
}

type mappingRule struct {
	sourceStart, destinationStart, mappingRange int
}

type Mapper interface {
	Map(source int) (destination int)
}

func (mp mappingRules) Map(source int) (destination int) {
	for _, rule := range mp.rules {
		if source >= rule.sourceStart && source < rule.sourceStart+rule.mappingRange {
			diff := source - rule.sourceStart
			return rule.destinationStart + diff
		}
	}
	return source
}

func getTestLines() (taskLines []string) {
	test := "seeds: 79 14 55 13\n\nseed-to-soil map:\n50 98 2\n52 50 48\n\nsoil-to-fertilizer map:\n0 15 37\n37 52 2\n39 0 15\n\nfertilizer-to-water map:\n49 53 8\n0 11 42\n42 0 7\n57 7 4\n\nwater-to-light map:\n88 18 7\n18 25 70\n\nlight-to-temperature map:\n45 77 23\n81 45 19\n68 64 13\n\ntemperature-to-humidity map:\n0 69 1\n1 0 69\n\nhumidity-to-location map:\n60 56 37\n56 93 4"
	return strings.Split(test, "\n")
}
