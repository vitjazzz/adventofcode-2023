package day13

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"slices"
	"strings"
)

const (
	VERTICAL MirrorType = iota
	HORIZONTAL
)
const EXPECTED_SMUDGES = 0

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/13/input", false)
	//taskLines := getTestLines()
	patterns := getPatterns(taskLines)
	res := 0
	for _, pattern := range patterns {
		mirrorType, position := calculateMirror(pattern)
		fmt.Printf("mirrorType = %d, position = %d\n", mirrorType, position)
		switch mirrorType {
		case VERTICAL:
			res += position
		case HORIZONTAL:
			res += position * 100
		}
	}
	fmt.Printf("res = %d\n", res)
}

func calculateMirror(pattern Pattern) (MirrorType, int) {
	verticalPosition := calculateMirrorPosition(pattern.symbols)
	if verticalPosition > 0 {
		return VERTICAL, verticalPosition
	}
	invertedSymbols := invertSymbols(pattern.symbols)
	return HORIZONTAL, calculateMirrorPosition(invertedSymbols)
}

func calculateMirrorPosition(symbols [][]string) int {
	possibleMirrorsPerRow := make(map[int][]int)
	for i := 0; i < len(symbols); i++ {
		var possibleMirrors []Mirror
		for j := 0; j < len(symbols[i]); j++ {
			var newPossibleMirrors []Mirror
			for _, possibleMirror := range possibleMirrors {
				mirroredPosition := len(possibleMirror.leftSymbols) - (j - possibleMirror.position) - 1

				if mirroredPosition < 0 || possibleMirror.leftSymbols[mirroredPosition] == symbols[i][j] {
					newPossibleMirrors = append(newPossibleMirrors, possibleMirror)
				}
			}
			possibleMirrors = newPossibleMirrors
			if j != len(symbols[i])-1 {
				possibleMirrors = append(possibleMirrors, Mirror{j + 1, symbols[i][:j+1]})
			}
		}
		var possibleMirrorsPositions []int
		for _, possibleMirror := range possibleMirrors {
			possibleMirrorsPositions = append(possibleMirrorsPositions, possibleMirror.position)
		}
		possibleMirrorsPerRow[i] = possibleMirrorsPositions
	}
	var mirrors []int
	for i := 0; i < len(symbols); i++ {
		possibleMirrors := possibleMirrorsPerRow[i]
		if i == 0 {
			mirrors = possibleMirrors
		} else {
			var newMirrors []int
			for _, mirror := range mirrors {
				if slices.Contains(possibleMirrors, mirror) {
					newMirrors = append(newMirrors, mirror)
				}
			}
			mirrors = newMirrors
		}
	}
	if len(mirrors) > 1 {
		fmt.Printf("Unexpected mirrors - %v", mirrors)
	} else if len(mirrors) == 0 {
		return 0
	}
	return mirrors[0]
}

func invertSymbols(symbols [][]string) [][]string {
	res := make([][]string, len(symbols[0]))
	for j := 0; j < len(symbols[0]); j++ {
		res[j] = make([]string, len(symbols))
		for i := 0; i < len(symbols); i++ {
			res[j][i] = symbols[i][j]
		}
	}
	return res
}

func getPatterns(lines []string) (res []Pattern) {
	currentPattern := Pattern{}
	for _, line := range lines {
		if line == "" && len(currentPattern.symbols) == 0 {
			continue
		} else if line == "" && len(currentPattern.symbols) > 0 {
			res = append(res, currentPattern)
			currentPattern = Pattern{}
		} else {
			symbols := strings.Split(line, "")
			currentPattern.symbols = append(currentPattern.symbols, symbols)
		}
	}
	if len(currentPattern.symbols) > 0 {
		res = append(res, currentPattern)
	}
	return
}

type Pattern struct {
	symbols [][]string
}

type Mirror struct {
	position    int
	leftSymbols []string
}

type MirrorType int

func getTestLines() (taskLines []string) {
	test := "#.##..##.\n..#.##.#.\n##......#\n##......#\n..#.##.#.\n..##..##.\n#.#.##.#.\n\n#...##..#\n#....#..#\n..##..###\n#####.##.\n#####.##.\n..##..###\n#....#..#"
	return strings.Split(test, "\n")
}
