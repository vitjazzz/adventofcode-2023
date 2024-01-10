package day12

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var emptyRecord = Record{[]string{}, []int{}}

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/12/input")
	//taskLines := getTestLines()
	res := 0
	c := make(chan int)
	for _, line := range taskLines {
		go func(l string) {
			record := getRecord(l)
			record = unfold(record)
			storedArrangements := make(map[string]int)
			arrangements := getPossibleArrangementsCount(record, storedArrangements)
			fmt.Printf("line = %s, arrangements = %d\n", l, arrangements)
			c <- arrangements
		}(line)
		//res += arrangements
	}
	for i := 0; i < len(taskLines); i++ {
		res += <-c
	}
	fmt.Printf("res = %d\n", res)
}

func getPossibleArrangementsCount(record Record, storedArrangements map[string]int) int {
	arrangementKey := buildKey(record)
	if res, ok := storedArrangements[arrangementKey]; ok {
		return res
	}
	if len(record.groups) == 0 {
		if slices.Contains(record.symbols, "#") {
			return 0
		} else {
			return 1
		}
	}
	res := 0
	for i := 0; i < len(record.symbols); i++ {
		if record.symbols[i] == "." {
			continue
		}
		newRecord, ok := tryRestoreGroup(record, i)
		if ok {
			res += getPossibleArrangementsCount(newRecord, storedArrangements)
		}
		if record.symbols[i] == "#" {
			break
		}
	}
	storedArrangements[arrangementKey] = res
	return res
}

func tryRestoreGroup(record Record, position int) (Record, bool) {
	groupSize := record.groups[0]
	if len(record.symbols) < position+groupSize {
		return emptyRecord, false
	}
	if len(record.symbols) > position+groupSize && record.symbols[position+groupSize] == "#" {
		return emptyRecord, false
	}

	group := record.symbols[position : position+groupSize]
	if slices.Contains(group, ".") {
		return emptyRecord, false
	}

	newGroups := record.groups[1:]
	var resSymbols []string
	if len(record.symbols) == position+groupSize || len(record.symbols) == position+groupSize+1 {
		resSymbols = []string{}
	} else {
		resSymbols = record.symbols[position+groupSize+1:]
	}
	return Record{resSymbols, newGroups}, true
}

func buildKey(record Record) string {
	symbolsKey := strings.Join(record.symbols, "")
	groupsKey := ""
	for i, group := range record.groups {
		groupsKey += strconv.Itoa(group)
		if i != len(record.groups)-1 {
			groupsKey += ","
		}
	}
	return symbolsKey + "-" + groupsKey
}

func getRecord(line string) Record {
	symbolsStr := strings.Split(line, " ")[0]
	symbols := strings.Split(symbolsStr, "")
	groupsStr := strings.Split(line, " ")[1]
	groups := adventutils.ParseNumbers(groupsStr, ",")
	return Record{symbols, groups}
}

func unfold(record Record) Record {
	newGroups := record.groups
	newSymbols := record.symbols
	for i := 0; i < 4; i++ {
		newGroups = append(newGroups, record.groups...)
		newSymbols = append(newSymbols, "?")
		newSymbols = append(newSymbols, record.symbols...)
	}
	return Record{newSymbols, newGroups}
}

type Record struct {
	symbols []string
	groups  []int
}

func getTestLines() (taskLines []string) {
	test := "???.### 1,1,3\n.??..??...?##. 1,1,3\n?#?#?#?#?#?#?#? 1,3,1,6\n????.#...#... 4,1,1\n????.######..#####. 1,6,5\n?###???????? 3,2,1"
	return strings.Split(test, "\n")
}
