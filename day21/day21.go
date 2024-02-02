package day21

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"math"
	"strings"
)

const NOT_REACHED = -1
const REACHED_ON_EVERY_STEP = math.MaxInt
const STEPS = 64

var printedSteps = make(map[int]int)

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/21/input", true)
	//taskLines := getTestLines()
	board, start := buildBoard(taskLines)
	//printBoard(board)

	startTask := task{start, 0}
	tasks := []task{startTask}
	for len(tasks) > 0 {
		tasks = processTask(board, tasks)
	}
	//printBoard(board, STEPS)
	res := 0
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			t := board[i][j]
			if t.blocked || t.stepReached == NOT_REACHED {
				continue
			}
			if t.stepReached == REACHED_ON_EVERY_STEP {
				res++
			}
			if (STEPS-t.stepReached)%2 == 0 {
				res++
			}
		}
	}

	fmt.Printf("res = %d\n", res)
}

func processTask(board [][]*tile, tasks []task) []task {
	tsk := tasks[0]
	if _, ok := printedSteps[tsk.step]; !ok {
		fmt.Printf("Board before step %d:\n", tsk.step)
		//printBoard(board, tsk.step)
		printedSteps[tsk.step] = tsk.step
	}
	tasks = tasks[1:]
	if tsk.step > STEPS {
		return tasks
	}
	tl := board[tsk.c.i][tsk.c.j]
	if tl.blocked || tl.stepReached == REACHED_ON_EVERY_STEP || (tl.stepReached >= 0 && (tsk.step-tl.stepReached)%2 == 0) {
		return tasks
	}
	if tl.stepReached == NOT_REACHED {
		tl.stepReached = tsk.step
	} else {
		tl.stepReached = REACHED_ON_EVERY_STEP
	}
	left := coordinate{tsk.c.i, tsk.c.j - 1}
	right := coordinate{tsk.c.i, tsk.c.j + 1}
	top := coordinate{tsk.c.i - 1, tsk.c.j}
	bot := coordinate{tsk.c.i + 1, tsk.c.j}
	nextStep := tsk.step + 1
	tasks = append(tasks, task{top, nextStep}, task{bot, nextStep}, task{left, nextStep}, task{right, nextStep})
	return tasks
}

func buildBoard(lines []string) (res [][]*tile, _ coordinate) {
	length := len(lines) + 2
	width := len(lines[0]) + 2
	res = make([][]*tile, length)
	var start coordinate
	res[0] = blockedLine(width, 0)
	res[length-1] = blockedLine(width, length-1)
	for lineIndex, line := range lines {
		i := lineIndex + 1
		res[i] = make([]*tile, width)
		res[i][0] = &tile{coordinate{i, 0}, true, NOT_REACHED}
		res[i][width-1] = &tile{coordinate{i, width - 1}, true, NOT_REACHED}
		for symbIndex, symb := range strings.Split(line, "") {
			j := symbIndex + 1
			res[i][j] = &tile{coordinate{i, j}, isBlocked(symb), NOT_REACHED}
			if isStart(symb) {
				start = coordinate{i, j}
			}
		}
	}
	return res, start
}

func blockedLine(width, i int) []*tile {
	res := make([]*tile, width)
	for j := 0; j < width; j++ {
		res[j] = &tile{coordinate{i, j}, true, NOT_REACHED}
	}
	return res
}

func isBlocked(symbol string) bool {
	if symbol == "#" {
		return true
	}
	return false
}
func isStart(symbol string) bool {
	if symbol == "S" {
		return true
	}
	return false
}

func printBoard(board [][]*tile, currentStep int) {
	fmt.Println()
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			t := board[i][j]
			if t.blocked {
				fmt.Print("#")
			} else {
				if t.stepReached == REACHED_ON_EVERY_STEP || (t.stepReached >= 0 && (currentStep-t.stepReached)%2 == 0) {
					fmt.Print("O")
				} else {
					fmt.Print(".")
				}
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

type tile struct {
	c           coordinate
	blocked     bool
	stepReached int
}

type coordinate struct {
	i, j int
}

type task struct {
	c    coordinate
	step int
}

func getTestLines() (taskLines []string) {
	test := "...........\n.....###.#.\n.###.##..#.\n..#.#...#..\n....#.#....\n.##..S####.\n.##..#...#.\n.......##..\n.##.#.####.\n.##..##.##.\n..........."
	return strings.Split(test, "\n")
}
func getTestLines2() (taskLines []string) {
	test := "...........\n.....###.#.\n.###.##..#.\n..#.#######\n....#......\n.##.#S####.\n.##..#...#.\n.......##..\n.##.#.####.\n.##..##.##.\n..........."
	return strings.Split(test, "\n")
}
func getTestLines3() (taskLines []string) {
	test := "..#.#...#..\n....#.#....\n.##..S####.\n.##..#...#."
	return strings.Split(test, "\n")
}
