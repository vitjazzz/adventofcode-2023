package day11

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"math"
	"slices"
	"strings"
)

const multiplier = 1000000

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/11/input", true)
	//taskLines := getTestLines()
	universe := getUniverse(taskLines)
	expandedLines := getExpandedLines(universe)
	expandedColumns := getExpandedColumns(universe)
	galaxies := getGalaxies(universe)
	shortestPaths := calcShortestPaths(galaxies, expandedLines, expandedColumns)
	res := 0
	for _, l := range shortestPaths {
		res += l
	}
	fmt.Printf("res = %d\n", res)
}

func calcShortestPaths(galaxies []galaxy, expandedLines, expandedColumns []int) []int {
	calculated := make(map[int][]int)
	var res []int
	for i := 0; i < len(galaxies); i++ {
		for j := 0; j < len(galaxies); j++ {
			if i == j {
				continue
			}
			if slices.Contains(calculated[i], j) {
				continue
			}
			from := galaxies[i]
			to := galaxies[j]
			expandedLinesInPath := getExpansionsInRange(expandedLines, from.i, to.i)
			expandedColsInPath := getExpansionsInRange(expandedColumns, from.j, to.j)
			shortestPath := int(math.Abs(float64(from.i-to.i))+math.Abs(float64(from.j-to.j))) +
				(len(expandedLinesInPath)+len(expandedColsInPath))*(multiplier-1)
			res = append(res, shortestPath)
			calculated[i] = append(calculated[i], j)
			calculated[j] = append(calculated[j], i)
		}
	}
	return res
}

func getExpansionsInRange(expansions []int, a, b int) []int {
	var res []int
	var from, to int
	if a < b {
		from = a
		to = b
	} else {
		from = b
		to = a
	}
	for _, expansion := range expansions {
		if expansion > from && expansion < to {
			res = append(res, expansion)
		}
	}
	return res
}

func getGalaxies(universe [][]string) []galaxy {
	var res []galaxy
	currentId := 1
	for i := 0; i < len(universe); i++ {
		for j := 0; j < len(universe[i]); j++ {
			if universe[i][j] == "#" {
				res = append(res, galaxy{currentId, i, j})
				currentId++
			}
		}
	}
	return res
}

func getExpandedLines(universe [][]string) []int {
	var res []int
	for i := 0; i < len(universe); i++ {
		galaxyFound := false
		for j := 0; j < len(universe[i]); j++ {
			if universe[i][j] == "#" {
				galaxyFound = true
				break
			}
		}
		if !galaxyFound {
			res = append(res, i)
		}
	}
	return res
}

func getExpandedColumns(universe [][]string) []int {
	var res []int
	for j := 0; j < len(universe[0]); j++ {
		galaxyFound := false
		for i := 0; i < len(universe); i++ {
			if universe[i][j] == "#" {
				galaxyFound = true
				break
			}
		}
		if !galaxyFound {
			res = append(res, j)
		}
	}
	return res
}

func getUniverse(lines []string) [][]string {
	var res [][]string
	for _, line := range lines {
		res = append(res, strings.Split(line, ""))
	}
	return res
}

type galaxy struct {
	id   int
	i, j int
}

func getTestLines() (taskLines []string) {
	test := "...#......\n.......#..\n#.........\n..........\n......#...\n.#........\n.........#\n..........\n.......#..\n#...#....."
	return strings.Split(test, "\n")
}
