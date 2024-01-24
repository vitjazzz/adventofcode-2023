package day17

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"math"
	"strings"
	"time"
)

const MAX_STEPS = 3

var totalTime int64
var totalIterations int64

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/17/input", true)
	//taskLines := getTestLines()
	board := buildBoard(taskLines)

	botTask := task{coordinate{0, 0}, coordinate{1, 0}, MAX_STEPS - 1, MAX_STEPS, 0}
	rightTask := task{coordinate{0, 0}, coordinate{0, 1}, MAX_STEPS, MAX_STEPS - 1, 0}
	tasks := []task{rightTask, botTask}
	for len(tasks) > 0 {
		tasks = calculateShortestPath(board, tasks)

		if totalIterations != 0 && totalIterations%100_000 == 0 {
			fmt.Printf("Avg time - %d\n", totalTime/totalIterations)
		}
	}

	fmt.Printf("res = %d\n", board[len(board)-1][len(board[0])-1].getShortestPathWeight())
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
	toTile := board[to.i][to.j]
	newPath := path{from, t.currentWeight + toTile.weight, t.horizontalStepsLeft, t.verticalStepsLeft}
	if !isNewPathRelevant(toTile.shortestPaths, newPath) {
		return tasks
	}
	addNewPath(toTile, newPath)
	//toTile.shortestPaths = append(toTile.shortestPaths, newPath)
	sameDirTask := task{to, dir, t.verticalStepsLeft - adventutils.Abs(dir.i),
		t.horizontalStepsLeft - adventutils.Abs(dir.j), newPath.weight}
	newDir1 := coordinate{dir.i - dir.i - dir.j, dir.j - dir.j - dir.i}
	newDirTask1 := task{to, newDir1, MAX_STEPS - adventutils.Abs(newDir1.i),
		MAX_STEPS - adventutils.Abs(newDir1.j), newPath.weight}
	newDir2 := coordinate{dir.i - dir.i + dir.j, dir.j - dir.j + dir.i}
	newDirTask2 := task{to, newDir2, MAX_STEPS - adventutils.Abs(newDir2.i),
		MAX_STEPS - adventutils.Abs(newDir2.j), newPath.weight}
	start := time.Now()
	tasks = append(tasks, sameDirTask, newDirTask1, newDirTask2)
	totalTime += time.Now().Sub(start).Nanoseconds()
	totalIterations++
	return tasks
}

func addNewPath(toTile *tile, newPath path) {
	var newShortestPaths []path
	for _, p := range toTile.shortestPaths {
		if p.horizontalStepsLeft > newPath.horizontalStepsLeft ||
			p.verticalStepsLeft > newPath.verticalStepsLeft ||
			p.weight < newPath.weight {
			newShortestPaths = append(newShortestPaths, p)
		}
	}
	newShortestPaths = append(newShortestPaths, newPath)
	toTile.shortestPaths = newShortestPaths
}

func isNewPathRelevant(paths []path, newPath path) bool {
	for _, p := range paths {
		if p.horizontalStepsLeft >= newPath.horizontalStepsLeft &&
			p.verticalStepsLeft >= newPath.verticalStepsLeft &&
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
	return &tile{c, weight, []path{}}
}

type tile struct {
	c             coordinate
	weight        int
	shortestPaths []path
}

type path struct {
	from                                   coordinate
	weight                                 int
	horizontalStepsLeft, verticalStepsLeft int
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
}

func getTestLines() (taskLines []string) {
	test := "2413432311323\n3215453535623\n3255245654254\n3446585845452\n4546657867536\n1438598798454\n4457876987766\n3637877979653\n4654967986887\n4564679986453\n1224686865563\n2546548887735\n4322674655533"
	return strings.Split(test, "\n")
}
