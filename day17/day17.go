package day17

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"math"
	"strings"
)

const MAX_STEPS = 10
const MIN_STEPS = 4

var totalTime int64
var totalIterations int64

// 1067 is too high

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/17/input", true)
	//taskLines := getTestLines()
	board := buildBoard(taskLines)

	botTask := task{coordinate{0, 0}, coordinate{1, 0}, MAX_STEPS - 1, MAX_STEPS, 0, nil}
	rightTask := task{coordinate{0, 0}, coordinate{0, 1}, MAX_STEPS, MAX_STEPS - 1, 0, nil}
	tasks := []task{rightTask, botTask}
	for len(tasks) > 0 {
		tasks = calculateShortestPath(board, tasks)

		if totalIterations != 0 && totalIterations%100_000 == 0 {
			fmt.Printf("Avg time - %d\n", totalTime/totalIterations)
		}
	}
	restorePath(board)
	printBoard(board)

	fmt.Printf("res = %d\n", board[len(board)-1][len(board[0])-1].getShortestPathWeight())
}

func restorePath(board [][]*tile) {
	currentTile := board[len(board)-1][len(board[0])-1]
	shortestPathWeight := currentTile.getShortestPathWeight()
	var shortestPath *path
	for _, p := range currentTile.shortestPaths {
		if p.weight == shortestPathWeight {
			shortestPath = &p
			break
		}
	}
	currentTile.symbol = getSymbol(shortestPath.from, currentTile.c)
	for {
		currentTile = board[shortestPath.from.i][shortestPath.from.j]
		if currentTile.c.i == 0 && currentTile.c.j == 0 {
			return
		}
		shortestPath = shortestPath.prevPath
		if shortestPath == nil {
			break
		}
		currentTile.symbol = getSymbol(shortestPath.from, currentTile.c)
	}
}

func getSymbol(from, to coordinate) string {
	if to.i-from.i == 1 {
		return "v"
	} else if to.i-from.i == -1 {
		return "^"
	} else if to.j-from.j == 1 {
		return ">"
	} else {
		return "<"
	}
}

func printBoard(board [][]*tile) {
	fmt.Println()
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			t := board[i][j]
			if t.symbol == "" {
				fmt.Print(t.weight)
			} else {
				fmt.Print(t.symbol)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func calculateShortestPath(board [][]*tile, tasks []task) []task {
	t := tasks[0]
	tasks = tasks[1:]
	if t.verticalStepsLeft < 0 || t.horizontalStepsLeft < 0 {
		return tasks
	}
	from := t.c
	dir := t.direction
	to := coordinate{from.i + dir.i, from.j + dir.j}
	if to.i < 0 || to.i >= len(board) || to.j < 0 || to.j >= len(board[0]) {
		return tasks
	}
	if to.i == len(board)-1 && to.j == len(board[0])-1 &&
		t.horizontalStepsLeft > MAX_STEPS-MIN_STEPS && t.verticalStepsLeft > MAX_STEPS-MIN_STEPS {
		return tasks
	}
	toTile := board[to.i][to.j]
	newPath := path{from, t.currentWeight + toTile.weight, t.horizontalStepsLeft, t.verticalStepsLeft, t.prevPath, dir}
	if !isNewPathRelevant(toTile.shortestPaths, newPath) {
		return tasks
	}
	replaceNewPath(toTile, newPath)
	//addNewPath(toTile, newPath)
	sameDirTask := task{to, dir, t.verticalStepsLeft - adventutils.Abs(dir.i),
		t.horizontalStepsLeft - adventutils.Abs(dir.j), newPath.weight, &newPath}
	if t.horizontalStepsLeft > MAX_STEPS-MIN_STEPS && t.verticalStepsLeft > MAX_STEPS-MIN_STEPS {
		tasks = append(tasks, sameDirTask)
		return tasks
	}
	newDir1 := coordinate{dir.i - dir.i - dir.j, dir.j - dir.j - dir.i}
	newDirTask1 := task{to, newDir1, MAX_STEPS - adventutils.Abs(newDir1.i),
		MAX_STEPS - adventutils.Abs(newDir1.j), newPath.weight, &newPath}
	newDir2 := coordinate{dir.i - dir.i + dir.j, dir.j - dir.j + dir.i}
	newDirTask2 := task{to, newDir2, MAX_STEPS - adventutils.Abs(newDir2.i),
		MAX_STEPS - adventutils.Abs(newDir2.j), newPath.weight, &newPath}
	tasks = append(tasks, sameDirTask, newDirTask1, newDirTask2)
	return tasks
}

func replaceNewPath(toTile *tile, newPath path) {
	indexToReplace := -1
	for i, p := range toTile.shortestPaths {
		if p.horizontalStepsLeft == newPath.horizontalStepsLeft &&
			p.verticalStepsLeft == newPath.verticalStepsLeft &&
			p.direction.i == newPath.direction.i &&
			p.direction.j == newPath.direction.j &&
			p.weight >= newPath.weight {
			indexToReplace = i
		}
	}
	if indexToReplace == -1 {
		toTile.shortestPaths = append(toTile.shortestPaths, newPath)
	} else {
		toTile.shortestPaths[indexToReplace] = newPath
	}
}

func isNewPathRelevant(paths []path, newPath path) bool {
	for _, p := range paths {
		if p.horizontalStepsLeft == newPath.horizontalStepsLeft &&
			p.verticalStepsLeft == newPath.verticalStepsLeft &&
			p.direction.i == newPath.direction.i &&
			p.direction.j == newPath.direction.j &&
			p.weight <= newPath.weight {
			return false
		}
	}
	return true
}

func buildBoard(lines []string) [][]*tile {
	res := make([][]*tile, len(lines))
	for i, line := range lines {
		res[i] = make([]*tile, len(line))
		for j, weight := range adventutils.ParseNumbers(line, "") {
			res[i][j] = createTile(coordinate{i, j}, weight)
		}
	}
	return res
}

func createTile(c coordinate, weight int) *tile {
	return &tile{c, weight, []path{}, ""}
}

type tile struct {
	c             coordinate
	weight        int
	shortestPaths []path
	symbol        string
}

type path struct {
	from                                   coordinate
	weight                                 int
	horizontalStepsLeft, verticalStepsLeft int
	prevPath                               *path
	direction                              coordinate
}

func (t tile) getShortestPathWeight() int {
	if len(t.shortestPaths) == 0 {
		return t.weight
	}
	minWeight := math.MaxInt
	for _, p := range t.shortestPaths {
		if p.weight < minWeight {
			minWeight = p.weight
		}
	}
	return minWeight
}

type coordinate struct {
	i, j int
}

type task struct {
	c, direction                           coordinate
	verticalStepsLeft, horizontalStepsLeft int
	currentWeight                          int
	prevPath                               *path
}

func getTestLines() (taskLines []string) {
	test := "2413432311323\n3215453535623\n3255245654254\n3446585845452\n4546657867536\n1438598798454\n4457876987766\n3637877979653\n4654967986887\n4564679986453\n1224686865563\n2546548887735\n4322674655533"
	return strings.Split(test, "\n")
}

func getTestLines2() (taskLines []string) {
	test := "111111111111\n999999999991\n999999999991\n999999999991\n999999999991"
	return strings.Split(test, "\n")
}
