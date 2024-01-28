package day18

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"math"
	"math/big"
	"sort"
	"strings"
)

func Run() {
	//taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/18/input", true)
	taskLines := getTestLines()
	commands := getCommands(taskLines)
	digLines := buildLines(commands)
	linesLinks := rebuildLines(digLines)
	//linesLinks[len(linesLinks)-1].dir = coordinate{linesLinks[len(linesLinks)-1].dir.i - 1, linesLinks[len(linesLinks)-1].dir.j}
	res := calculateSurface(linesLinks)
	//printBoard(board)

	fmt.Printf("res = %d\n", res)
}

func calculateSurface(lines []*line) int {
	maxI := math.MinInt
	for _, l := range lines {
		maxI = adventutils.Max(maxI, l.c.i)
		maxI = adventutils.Max(maxI, l.c.i+l.dir.i)
	}
	res := 0
	for _, l := range lines {
		res += adventutils.Abs(l.dir.i + l.dir.j)
	}
	for i := 0; i < maxI; i++ {
		if i%1_000_000 == 0 {
			fmt.Printf("Current - %d, all - %d", i, maxI)
		}
		j := 0
		prevVerticalLine := nextVerticalLine(i, j, lines)
		if prevVerticalLine == nil {
			continue
		}
		j = prevVerticalLine.c.j + 1
		currentVerticalLine := nextVerticalLine(i, j, lines)
		inside := true
		for currentVerticalLine != nil {
			prevLineDirection := getLineDirection(i, prevVerticalLine)
			currentLineDirection := getLineDirection(i, currentVerticalLine)
			if inside && (currentLineDirection == 0 || prevLineDirection == 0) {
				res += currentVerticalLine.c.j - prevVerticalLine.c.j - 1
			}
			if !inside {
				if currentLineDirection == 0 || prevLineDirection == 0 || prevLineDirection == currentLineDirection {
					inside = true
				}
			} else if inside {
				if currentLineDirection == 0 || prevLineDirection == 0 || prevLineDirection == currentLineDirection {
					inside = false
				}
			}

			prevVerticalLine = currentVerticalLine
			j = currentVerticalLine.c.j + 1
			currentVerticalLine = nextVerticalLine(i, j, lines)
		}
	}
	return res
}

func orderedVerticalLines(i int, lines []*line) []*line {
	var possibleVerticalLines []*line
	for _, l := range lines {
		if l.dir.i == 0 {
			continue
		}
		if l.dir.i > 0 {
			if i >= l.c.i && i <= l.c.i+l.dir.i {
				possibleVerticalLines = append(possibleVerticalLines, l)
			}
		} else if l.dir.i < 0 {
			if i <= l.c.i && i >= l.c.i+l.dir.i {
				possibleVerticalLines = append(possibleVerticalLines, l)
			}
		}
	}
	sort.Slice(possibleVerticalLines, func(i, j int) bool {
		return possibleVerticalLines[i].c.j < possibleVerticalLines[j].c.j
	})
	return possibleVerticalLines
}

func nextVerticalLine(i, j int, lines []*line) *line {
	var possibleVerticalLines []*line
	for _, l := range lines {
		if l.dir.i == 0 {
			continue
		}
		if l.c.j < j {
			continue
		}
		if l.dir.i > 0 {
			if i >= l.c.i && i <= l.c.i+l.dir.i {
				possibleVerticalLines = append(possibleVerticalLines, l)
			}
		} else if l.dir.i < 0 {
			if i <= l.c.i && i >= l.c.i+l.dir.i {
				possibleVerticalLines = append(possibleVerticalLines, l)
			}
		}
	}
	if len(possibleVerticalLines) == 0 {
		return nil
	}
	res := possibleVerticalLines[0]
	for _, vertLine := range possibleVerticalLines {
		if vertLine.c.j < res.c.j {
			res = vertLine
		}
	}
	return res
}

func onLine(i, j int, l line) bool {
	if l.dir.i == 0 && i != l.c.i ||
		l.dir.j == 0 && j != l.c.j {
		return false
	}
	if l.dir.i > 0 {
		return i >= l.c.i && i < l.c.i+l.dir.i
	} else if l.dir.i < 0 {
		return i <= l.c.i && i > l.c.i+l.dir.i
	} else if l.dir.j > 0 {
		return j >= l.c.j && j < l.c.j+l.dir.j
	} else if l.dir.j < 0 {
		return j <= l.c.j && j > l.c.j+l.dir.j
	}
	return false
}

func buildBoard(lines []line) [][]string {
	minI, maxI, minJ, maxJ := math.MaxInt, math.MinInt, math.MaxInt, math.MinInt
	for _, l := range lines {
		minI = adventutils.Min(minI, l.c.i)
		minI = adventutils.Min(minI, l.c.i+l.dir.i)
		maxI = adventutils.Max(maxI, l.c.i)
		maxI = adventutils.Max(maxI, l.c.i+l.dir.i)
		minJ = adventutils.Min(minJ, l.c.j)
		minJ = adventutils.Min(minJ, l.c.j+l.dir.j)
		maxJ = adventutils.Max(maxJ, l.c.j)
		maxJ = adventutils.Max(maxJ, l.c.j+l.dir.j)
	}
	length := maxI - minI + 3
	width := maxJ - minJ + 3
	res := make([][]string, length)
	for i := 0; i < length; i++ {
		res[i] = make([]string, width)
		for j := 0; j < width; j++ {
			res[i][j] = "."
		}
	}
	for _, l := range lines {
		iStart := l.c.i - minI + 1
		jStart := l.c.j - minJ + 1
		iInc, jInc := getIncrements(l.dir)
		for i, j := iStart, jStart; isLineEnded(i, j, iStart, jStart, l.dir); i, j = i+iInc, j+jInc {
			res[i][j] = "#"
		}
	}
	return res
}

func rebuildLines(lines []line) []*line {
	minI, minJ := math.MaxInt, math.MaxInt
	for _, l := range lines {
		minI = adventutils.Min(minI, l.c.i)
		minI = adventutils.Min(minI, l.c.i+l.dir.i)
		minJ = adventutils.Min(minJ, l.c.j)
		minJ = adventutils.Min(minJ, l.c.j+l.dir.j)
	}
	var res []*line
	for _, l := range lines {
		newCoordinate := coordinate{l.c.i - minI, l.c.j - minJ}
		res = append(res, &line{newCoordinate, l.dir, l.colour})
	}
	return res
}

func isLineEnded(i, j int, iStart, jStart int, dir coordinate) bool {
	if dir.i > 0 {
		return i < iStart+dir.i
	} else if dir.i < 0 {
		return i > iStart+dir.i
	} else if dir.j > 0 {
		return j < jStart+dir.j
	} else {
		return j > jStart+dir.j
	}
}

func getIncrements(dir coordinate) (iInc, jInc int) {
	if dir.j == 0 {
		return dir.i / adventutils.Abs(dir.i), 0
	} else {
		return 0, dir.j / adventutils.Abs(dir.j)
	}
}

func buildLines(commands []cmd) []line {
	var res []line
	i := 0
	j := 0
	for _, cmd := range commands {
		c := coordinate{i, j}
		var dir coordinate
		switch cmd.dir {
		case "R":
			dir = coordinate{0, cmd.length}
		case "L":
			dir = coordinate{0, -cmd.length}
		case "U":
			dir = coordinate{-cmd.length, 0}
		case "D":
			dir = coordinate{cmd.length, 0}
		}
		res = append(res, line{c, dir, cmd.colour})
		i += dir.i
		j += dir.j
	}
	return res
}

func getCommands(lines []string) []cmd {
	res := make([]cmd, len(lines))
	for i, line := range lines {
		lineParts := strings.Split(line, " ")
		colour := lineParts[2][1 : len(lineParts[2])-1]
		dir := getDir(string(colour[len(colour)-1]))
		length := getLength(colour[1 : len(colour)-1])
		res[i] = cmd{dir, length, ""}
	}
	return res
}

func getDir(dirCoded string) string {
	switch dirCoded {
	case "0":
		return "R"
	case "1":
		return "D"
	case "2":
		return "L"
	case "3":
		return "U"
	}
	return "UNKNOWN"
}

func getLength(lengthCoded string) int {
	n := new(big.Int)
	n.SetString(lengthCoded, 16)
	return int(n.Int64())
}

func printBoard(board [][]string) {
	fmt.Println()
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			fmt.Print(board[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
}

func getLineDirection(i int, l *line) int {
	minI := adventutils.Min(l.c.i, l.c.i+l.dir.i)
	maxI := adventutils.Max(l.c.i, l.c.i+l.dir.i)
	if i == minI {
		return -1
	} else if i == maxI {
		return 1
	} else {
		return 0
	}
}

type cmd struct {
	dir    string
	length int
	colour string
}

type line struct {
	c      coordinate
	dir    coordinate
	colour string
}

type coordinate struct {
	i, j int
}

func getTestLines() (taskLines []string) {
	test := "R 6 (#70c710)\nD 5 (#0dc571)\nL 2 (#5713f0)\nD 2 (#d2c081)\nR 2 (#59c680)\nD 2 (#411b91)\nL 5 (#8ceee2)\nU 2 (#caa173)\nL 1 (#1b58a2)\nU 2 (#caa171)\nR 2 (#7807d2)\nU 3 (#a77fa3)\nL 2 (#015232)\nU 2 (#7a21e3)"
	return strings.Split(test, "\n")
}
