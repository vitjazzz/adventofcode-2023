package day5

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"math"
	"strings"
)

const seedsKey = "seeds: "
const mapKey = " map:"

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/5/input")
	//taskLines := getTestLines()
	seedsRanges := getSeedsAdvanced(taskLines)
	maxSeed := findMaxSeed(seedsRanges)
	allMappingRules := getMappingRules(taskLines)
	allTransitiveRules := calculateTransitiveRules(rangeStruct{0, maxSeed}, "seed", allMappingRules)
	//closestLocation := calculateClosestLocation(seedsRanges, allMappingRules)
	closestLocation := calculateClosestDestination(seedsRanges, allTransitiveRules)
	fmt.Printf("closestLocation - %d\n", closestLocation)
}

func calculateClosestDestination(sourceRanges []rangeStruct, allRules []mappingRule) int {
	closestDestination := math.MaxInt
	for _, sourceRange := range sourceRanges {
		for _, rule := range allRules {
			if (sourceRange.end() <= rule.sourceStart) || sourceRange.start >= rule.sourceEnd() {
				continue
			}
			diff := sourceRange.start - rule.sourceStart
			diff = adventutils.Max(0, diff)
			destination := rule.destinationStart + diff
			if closestDestination > destination {
				closestDestination = destination
			}
		}
	}
	return closestDestination
}

func findMaxSeed(seedsRanges []rangeStruct) int {
	maxNumber := 0
	for _, seedRange := range seedsRanges {
		seedRangeEnd := seedRange.start + seedRange.rangeSize
		if maxNumber < seedRangeEnd {
			maxNumber = seedRangeEnd
		}
	}
	return maxNumber
}

func calculateTransitiveRules(currentRange rangeStruct, source string, allMappingRules []mappingRules) []mappingRule {
	if source == "location" {
		return []mappingRule{{currentRange.start, currentRange.start, currentRange.rangeSize}}
	}
	currentRules := getMappingRule(source, allMappingRules)
	rulesInRange := getRulesInRange(currentRange, currentRules)
	//fmt.Printf("Source - %s, Rules in range - %v\n", source, rulesInRange)
	var res []mappingRule
	for _, sourceRule := range rulesInRange {
		destinationRules := calculateTransitiveRules(
			rangeStruct{sourceRule.destinationStart, sourceRule.mappingRange},
			currentRules.destination,
			allMappingRules)
		diff := sourceRule.sourceStart - sourceRule.destinationStart
		for _, destinationRule := range destinationRules {
			sourceStart := destinationRule.sourceStart + diff
			transitiveRule := mappingRule{sourceStart, destinationRule.destinationStart, destinationRule.mappingRange}
			res = append(res, transitiveRule)
		}
	}
	return res
}

func getRulesInRange(currentRange rangeStruct, rules mappingRules) (newRules []mappingRule) {
	rangeEnd := currentRange.start + currentRange.rangeSize
	for currentSourceStart := currentRange.start; currentSourceStart < rangeEnd; {
		nextRule := findNextRule(currentSourceStart, rules)
		nextRule.mappingRange = adventutils.Min(rangeEnd-nextRule.sourceStart, nextRule.mappingRange)
		newRules = append(newRules, nextRule)
		currentSourceStart += nextRule.mappingRange
	}
	return newRules
}

func findNextRule(sourceStart int, rules mappingRules) mappingRule {
	closestRuleStart := math.MaxInt
	for _, rule := range rules.rules {
		if sourceStart >= rule.sourceStart && sourceStart < rule.sourceEnd() {
			diff := sourceStart - rule.sourceStart
			return mappingRule{
				sourceStart:      sourceStart,
				destinationStart: rule.destinationStart + diff,
				mappingRange:     rule.mappingRange - diff,
			}
		}
		if rule.sourceStart > sourceStart && rule.sourceStart < closestRuleStart {
			closestRuleStart = rule.sourceStart
		}
	}
	mappingRange := closestRuleStart - sourceStart
	return mappingRule{sourceStart, sourceStart, mappingRange}
}

func getMappingRule(source string, rules []mappingRules) mappingRules {
	for _, rule := range rules {
		if rule.source == source {
			return rule
		}
	}
	fmt.Printf("Failed to find rule for %s source!!!\n", source)
	return mappingRules{}
}

func getSeedsRanges(lines []string) (seedsRanges []rangeStruct) {
	for _, line := range lines {
		if !strings.Contains(line, seedsKey) {
			continue
		}
		seedsStr := line[len(seedsKey):]
		seeds := adventutils.ParseNumbers(seedsStr, " ")
		for _, seed := range seeds {
			seedsRanges = append(seedsRanges, rangeStruct{start: seed, rangeSize: 1})
		}
	}
	return seedsRanges
}

func getSeedsAdvanced(lines []string) (seedsRanges []rangeStruct) {
	for _, line := range lines {
		if !strings.Contains(line, seedsKey) {
			continue
		}
		seedsStr := line[len(seedsKey):]
		seeds := adventutils.ParseNumbers(seedsStr, " ")
		for i := 0; i < len(seeds); i += 2 {
			seedsRanges = append(seedsRanges, rangeStruct{start: seeds[i], rangeSize: seeds[i+1]})

		}
	}
	return seedsRanges
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

type rangeStruct struct {
	start, rangeSize int
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
		if source >= rule.sourceStart && source < rule.sourceEnd() {
			diff := source - rule.sourceStart
			return rule.destinationStart + diff
		}
	}
	return source
}
func (rs rangeStruct) end() int {
	return rs.start + rs.rangeSize
}
func (mp mappingRule) sourceEnd() int {
	return mp.sourceStart + mp.mappingRange
}

func getTestLines() (taskLines []string) {
	test := "seeds: 79 14 55 13\n\nseed-to-soil map:\n50 98 2\n52 50 48\n\nsoil-to-fertilizer map:\n0 15 37\n37 52 2\n39 0 15\n\nfertilizer-to-water map:\n49 53 8\n0 11 42\n42 0 7\n57 7 4\n\nwater-to-light map:\n88 18 7\n18 25 70\n\nlight-to-temperature map:\n45 77 23\n81 45 19\n68 64 13\n\ntemperature-to-humidity map:\n0 69 1\n1 0 69\n\nhumidity-to-location map:\n60 56 37\n56 93 4"
	return strings.Split(test, "\n")
}
