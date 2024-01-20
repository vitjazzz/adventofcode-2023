package day14

import (
	"fmt"
	"strings"
)

const CYCLES = 100

func Run() {
	//taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/14/input", true)
	taskLines := getTestLines()
	board := getBoard(taskLines)
	for k := 0; k < CYCLES; k++ {
		for m := 0; m < 4; m++ {
			tiltNorth(board)
			board = rotateClockwise(board)
		}
		res := calculateNorth(board)
		fmt.Printf("After %d cycle, res - %d:\n", k+1, res)
		printBoard(board)
	}
	res := calculateNorth(board)
	fmt.Printf("res = %d\n", res)
}

func printBoard(board [][]string) {
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			fmt.Print(board[i][j])
		}
		fmt.Println()
	}
	fmt.Println()
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
