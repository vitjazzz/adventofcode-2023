package day23

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"math"
	"slices"
	"strings"
)

// 5046 is too low
func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/23/input", true)
	//taskLines := getTestLines()
	board := buildBoard(taskLines)
	printBoard(board, path{[]coordinate{}, 0})

	startPath := path{[]coordinate{}, 0}
	startTask := task{coordinate{0, 1}, startPath}
	tasks := []task{startTask}
	for len(tasks) > 0 {
		tasks = processTask(board, tasks)
	}

	res := board[len(board)-1][len(board[0])-2].bestPath
	printBoard(board, res)
	fmt.Printf("res = %d\n", res.length-1)
}

func processTask(board [][]*tile, tasks []task) []task {
	tsk := tasks[0]
	tasks = tasks[1:]
	if tsk.c.i < 0 || tsk.c.i == len(board) {
		return tasks
	}
	tl := board[tsk.c.i][tsk.c.j]
	if tl.symbol == "#" ||
		tl.bestPath.length > tsk.currentPath.length ||
		slices.Contains(tsk.currentPath.coordinates, tsk.c) {
		return tasks
	}
	newBestPath := make([]coordinate, len(tsk.currentPath.coordinates)+1)
	for i, c := range tsk.currentPath.coordinates {
		newBestPath[i] = c
	}
	newBestPath[len(tsk.currentPath.coordinates)] = tsk.c
	tl.bestPath = path{newBestPath, tsk.currentPath.length + 1}
	left := coordinate{tsk.c.i, tsk.c.j - 1}
	right := coordinate{tsk.c.i, tsk.c.j + 1}
	top := coordinate{tsk.c.i - 1, tsk.c.j}
	bot := coordinate{tsk.c.i + 1, tsk.c.j}
	//if tl.symbol == ">" {
	//	tasks = append(tasks, task{right, tl.bestPath})
	//} else if tl.symbol == "<" {
	//	tasks = append(tasks, task{left, tl.bestPath})
	//} else if tl.symbol == "^" {
	//	tasks = append(tasks, task{top, tl.bestPath})
	//} else if tl.symbol == "v" {
	//	tasks = append(tasks, task{bot, tl.bestPath})
	//} else if tl.symbol == "." {
	//	tasks = append(tasks, task{top, tl.bestPath}, task{bot, tl.bestPath}, task{left, tl.bestPath}, task{right, tl.bestPath})
	//} else {
	//	fmt.Printf("NOT VALID!!!!!!")
	//}
	tasks = append(tasks, task{top, tl.bestPath}, task{bot, tl.bestPath}, task{left, tl.bestPath}, task{right, tl.bestPath})

	return tasks
}

func buildBoard(lines []string) (res [][]*tile) {
	res = make([][]*tile, len(lines))
	for i, line := range lines {
		res[i] = make([]*tile, len(lines[i]))
		for j, symbol := range strings.Split(line, "") {
			res[i][j] = &tile{coordinate{i, j}, symbol, path{[]coordinate{}, math.MinInt}}
		}
	}
	return res
}

func printBoard(board [][]*tile, p path) {
	fmt.Println()
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			t := board[i][j]
			if slices.Contains(p.coordinates, coordinate{i, j}) {
				fmt.Print("O")
			} else {
				fmt.Print(t.symbol)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

type tile struct {
	c        coordinate
	symbol   string
	bestPath path
}

type path struct {
	coordinates []coordinate
	length      int
}

type coordinate struct {
	i, j int
}

type task struct {
	c           coordinate
	currentPath path
}

func getTestLines() (taskLines []string) {
	test := "#.#####################\n#.......#########...###\n#######.#########.#.###\n###.....#.>.>.###.#.###\n###v#####.#v#.###.#.###\n###.>...#.#.#.....#...#\n###v###.#.#.#########.#\n###...#.#.#.......#...#\n#####.#.#.#######.#.###\n#.....#.#.#.......#...#\n#.#####.#.#.#########v#\n#.#...#...#...###...>.#\n#.#.#v#######v###.###v#\n#...#.>.#...>.>.#.###.#\n#####v#.#.###v#.#.###.#\n#.....#...#...#.#.#...#\n#.#########.###.#.#.###\n#...###...#...#...#.###\n###.###.#.###v#####v###\n#...#...#.#.>.>.#.>.###\n#.###.###.#.###.#.#v###\n#.....###...###...#...#\n#####################.#"
	return strings.Split(test, "\n")
}
