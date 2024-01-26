package day18

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/18/input", true)
	//taskLines := getTestLines()
	commands := getCommands(taskLines)
	digLines := buildLines(commands)
	board := buildBoard(digLines)
	res := calculateSurface(board)
	//printBoard(board)

	fmt.Printf("res = %d\n", res)
}

func calculateSurface(board [][]string) int {
	res := 0
	for i := 0; i < len(board); i++ {
		inside := false
		inDirection := 0
		outDirection := 0
		for j := 0; j < len(board[i]); j++ {
			symbol := board[i][j]
			if symbol == "#" {
				res++
			}
			if symbol == "." && inside {
				res++
			}
			c := coordinate{i, j}
			if symbol == "#" && !inside && isVertical(c, board) {
				inDirection = getLineDirection(c, board)
				if outDirection == 0 || inDirection == 0 || outDirection == inDirection {
					inside = true
				}
			} else if symbol == "#" && inside && isVertical(c, board) {
				outDirection = getLineDirection(c, board)
				if inDirection == outDirection || inDirection == 0 || outDirection == 0 {
					inside = false
				}
			}
		}
	}
	return res
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
		dir := lineParts[0]
		length, _ := strconv.Atoi(lineParts[1])
		colour := lineParts[2][1 : len(lineParts[2])-1]
		res[i] = cmd{dir, length, colour}
	}
	return res
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

func isVertical(c coordinate, board [][]string) bool {
	return board[c.i][c.j] == "#" && (board[c.i-1][c.j] == "#" || board[c.i+1][c.j] == "#")
}
func topDirection(c coordinate, board [][]string) int {
	if board[c.i-1][c.j] == "#" {
		return -1
	} else if board[c.i+1][c.j] == "#" {
		return 1
	} else {
		return 0
	}
}
func getLineDirection(c coordinate, board [][]string) int {
	if board[c.i-1][c.j] == "#" && board[c.i+1][c.j] == "#" {
		return 0
	} else if board[c.i+1][c.j] == "#" {
		return 1
	} else if board[c.i-1][c.j] == "#" {
		return -1
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
