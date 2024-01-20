package day14

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"hash/fnv"
	"strings"
)

const CYCLES = 1000000000

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/14/input", true)
	//taskLines := getTestLines()
	board := getBoard(taskLines)
	boardByCycleMap := make(map[uint64]int)

	cycleStart := -1
	cycleEnd := -1
	for k := 0; k < CYCLES; k++ {
		boardStr := toString(board)
		boardHash := hash(boardStr)
		if cycle, ok := boardByCycleMap[boardHash]; ok {
			fmt.Printf("Current cycle - %d, last same cycle - %d:\n", k, cycle)
			cycleStart = cycle
			cycleEnd = k
			break
		} else {
			boardByCycleMap[boardHash] = k
		}

		for m := 0; m < 4; m++ {
			tiltNorth(board)
			board = rotateClockwise(board)
		}

	}
	cyclesLeft := (CYCLES - cycleStart) % (cycleEnd - cycleStart)
	for k := 0; k < cyclesLeft; k++ {
		for m := 0; m < 4; m++ {
			tiltNorth(board)
			board = rotateClockwise(board)
		}
	}
	res := calculateNorth(board)
	fmt.Printf("res = %d\n", res)
}

func toString(board [][]string) string {
	res := ""
	for i := 0; i < len(board); i++ {
		line := strings.Join(board[i], "")
		res += line
	}
	return res
}

func hash(str string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(str))
	return h.Sum64()
}

func tiltNorth(board [][]string) {
	for j := 0; j < len(board[0]); j++ {
		currentEmptySpace := -1
		for i := 0; i < len(board); i++ {
			if board[i][j] == "#" {
				currentEmptySpace = -1
			} else if board[i][j] == "O" {
				if currentEmptySpace == -1 {
					continue
				}
				board[i][j] = "."
				board[currentEmptySpace][j] = "O"
				currentEmptySpace += 1
			} else if board[i][j] == "." {
				if currentEmptySpace == -1 {
					currentEmptySpace = i
				}
			}
		}
	}
}

func rotateClockwise(board [][]string) [][]string {
	res := make([][]string, len(board[0]))
	for j, newI := 0, 0; j < len(board[0]); j, newI = j+1, newI+1 {
		res[newI] = make([]string, len(board))
		for i, newJ := len(board)-1, 0; i >= 0; i, newJ = i-1, newJ+1 {
			res[newI][newJ] = board[i][j]
		}
	}
	return res
}

func calculateNorth(board [][]string) int {
	res := 0
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] == "O" {
				res += len(board) - i
			}
		}
	}
	return res
}

func getBoard(lines []string) [][]string {
	res := make([][]string, len(lines))
	for i, line := range lines {
		res[i] = strings.Split(line, "")
	}
	return res
}

func getTestLines() (taskLines []string) {
	test := "O....#....\nO.OO#....#\n.....##...\nOO.#O....O\n.O.....O#.\nO.#..O.#.#\n..O..#O..O\n.......O..\n#....###..\n#OO..#...."
	return strings.Split(test, "\n")
}
