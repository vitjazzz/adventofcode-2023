package day16

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"strings"
)

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/16/input", false)
	//taskLines := getTestLines()
	bestRes := 0
	for i := 0; i < len(taskLines); i++ {
		board := buildBoard(taskLines)
		tasks := []beam{{i: i, j: -1, iDirection: 0, jDirection: 1}}
		for len(tasks) > 0 {
			tasks = processTask(board, tasks)
		}
		power := calculatePower(board)
		if power > bestRes {
			bestRes = power
		}
	}
	for i := 0; i < len(taskLines); i++ {
		board := buildBoard(taskLines)
		tasks := []beam{{i: i, j: len(board[i]), iDirection: 0, jDirection: -1}}
		for len(tasks) > 0 {
			tasks = processTask(board, tasks)
		}
		power := calculatePower(board)
		if power > bestRes {
			bestRes = power
		}
	}
	for j := 0; j < len(taskLines[0]); j++ {
		board := buildBoard(taskLines)
		tasks := []beam{{i: -1, j: j, iDirection: 1, jDirection: 0}}
		for len(tasks) > 0 {
			tasks = processTask(board, tasks)
		}
		power := calculatePower(board)
		if power > bestRes {
			bestRes = power
		}
	}
	for j := 0; j < len(taskLines[0]); j++ {
		board := buildBoard(taskLines)
		tasks := []beam{{i: len(board), j: j, iDirection: -1, jDirection: 0}}
		for len(tasks) > 0 {
			tasks = processTask(board, tasks)
		}
		power := calculatePower(board)
		if power > bestRes {
			bestRes = power
		}
	}

	fmt.Printf("res = %d\n", bestRes)
}

func calculatePower(board [][]*tile) int {
	res := 0
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j].activated {
				res++
			}
		}
	}
	return res
}

func processTask(board [][]*tile, tasks []beam) []beam {
	task := tasks[len(tasks)-1]
	tasks = tasks[:len(tasks)-1]
	for i, j := task.i+task.iDirection, task.j+task.jDirection; i >= 0 && j >= 0 && i < len(board) && j < len(board[i]); i, j = i+task.iDirection, j+task.jDirection {
		currentTile := board[i][j]
		if (currentTile.horizontalPass && task.jDirection != 0) ||
			(currentTile.verticalPass && task.iDirection != 0) {
			return tasks
		}
		currentTile.activated = true
		if currentTile.symbol == "/" {
			var newTask beam
			if task.iDirection == -1 {
				newTask = beam{i, j, 0, 1}
			} else if task.iDirection == 1 {
				newTask = beam{i, j, 0, -1}
			} else if task.jDirection == 1 {
				newTask = beam{i, j, -1, 0}
			} else if task.jDirection == -1 {
				newTask = beam{i, j, 1, 0}
			}
			tasks = append(tasks, newTask)
			return tasks
		} else if currentTile.symbol == "\\" {
			var newTask beam
			if task.iDirection == -1 {
				newTask = beam{i, j, 0, -1}
			} else if task.iDirection == 1 {
				newTask = beam{i, j, 0, 1}
			} else if task.jDirection == 1 {
				newTask = beam{i, j, 1, 0}
			} else if task.jDirection == -1 {
				newTask = beam{i, j, -1, 0}
			}
			tasks = append(tasks, newTask)
			return tasks
		}
		if task.jDirection != 0 {
			currentTile.horizontalPass = true
		}
		if task.iDirection != 0 {
			currentTile.verticalPass = true
		}
		if currentTile.symbol == "-" && task.iDirection != 0 {
			leftTask := beam{i, j, 0, -1}
			tasks = append(tasks, leftTask)
			rightTask := beam{i, j, 0, 1}
			tasks = append(tasks, rightTask)
			return tasks
		}
		if currentTile.symbol == "|" && task.jDirection != 0 {
			topTask := beam{i, j, -1, 0}
			tasks = append(tasks, topTask)
			botTask := beam{i, j, 1, 0}
			tasks = append(tasks, botTask)
			return tasks
		}
	}
	return tasks
}

func buildBoard(lines []string) [][]*tile {
	res := make([][]*tile, len(lines))
	for i, line := range lines {
		res[i] = make([]*tile, len(line))
		for j, symbol := range strings.Split(line, "") {
			res[i][j] = createTile(symbol)
		}
	}
	return res
}

func createTile(symbol string) *tile {
	return &tile{false, false, symbol, false}
}

type tile struct {
	horizontalPass, verticalPass bool
	symbol                       string
	activated                    bool
}

type beam struct {
	i, j                   int
	iDirection, jDirection int
}

func getTestLines() (taskLines []string) {
	test := ".|...\\....\n|.-.\\.....\n.....|-...\n........|.\n..........\n.........\\\n..../.\\\\..\n.-.-/..|..\n.|....-|.\\\n..//.|...."
	return strings.Split(test, "\n")
}
