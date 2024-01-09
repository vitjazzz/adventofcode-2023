package day10

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"math"
	"strings"
)

var emptyPoint = Point{-1, -1}

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/10/input")
	//taskLines := getTestLines5()
	strMatrix := getStringMatrix(taskLines)
	grid := initializeGrid(strMatrix)
	grid, groups := markLineGroups(grid)
	start := getStart(grid)
	start = markStartNode(groups, start)
	startGroup := groups[start.groupId]
	fmt.Printf("part 1 res %d\n", (len(startGroup)/2)+1)
	nodes := findNodesWithinLoop(grid, start.groupId)
	fmt.Printf("part 2 res %d\n", len(nodes))
}

func findNodesWithinLoop(grid [][]*Node, loopGroupId int) []*Node {
	var res []*Node
	for _, nodes := range grid {
		inside := false
		inDirection := 0
		outDirection := 0
		for _, node := range nodes {
			if node.groupId != loopGroupId && inside {
				res = append(res, node)
			} else if node.groupId == loopGroupId && !inside && node.isVertical() {
				firstDirection := node.coordinate.i - node.first.i
				secondDirection := node.coordinate.i - node.second.i
				inDirection = firstDirection + secondDirection
				if outDirection == 0 || inDirection == 0 || outDirection == inDirection {
					inside = true
				}
			} else if node.groupId == loopGroupId && inside && node.isVertical() {
				firstDirection := node.coordinate.i - node.first.i
				secondDirection := node.coordinate.i - node.second.i
				outDirection = firstDirection + secondDirection
				if inDirection == outDirection || inDirection == 0 || outDirection == 0 {
					inside = false
				}
			}
		}
	}
	return res
}

func markStartNode(groups map[int][]*Node, start *Node) *Node {
	for id, group := range groups {
		var linksToStart []*Node
		for _, node := range group {
			if node.first == start.coordinate || node.second == start.coordinate {
				linksToStart = append(linksToStart, node)
			}
		}
		if len(linksToStart) == 2 {
			start.groupId = id
			start.first = linksToStart[0].coordinate
			start.second = linksToStart[1].coordinate
			break
		}
	}
	return start
}

func getStart(grid [][]*Node) *Node {
	for _, nodes := range grid {
		for _, node := range nodes {
			if node.start {
				return node
			}
		}
	}
	return grid[math.MaxInt][math.MaxInt]
}

func markLineGroups(grid [][]*Node) ([][]*Node, map[int][]*Node) {
	length := len(grid)
	width := len(grid[0])
	groups := make(map[int][]*Node)
	for i := 0; i < length; i++ {
		for j := 0; j < width; j++ {
			node := grid[i][j]
			if node.first == emptyPoint || node.second == emptyPoint {
				node.groupId = 0
				continue
			}
			firstNode := grid[node.first.i][node.first.j]
			if connected(node, firstNode) {
				groups = mergeGroup(groups, node, firstNode)
			}
			secondNode := grid[node.second.i][node.second.j]
			if connected(node, secondNode) {
				groups = mergeGroup(groups, node, secondNode)
			}
		}
	}
	return grid, groups
}

func mergeGroup(group map[int][]*Node, first, second *Node) map[int][]*Node {
	if first.groupId <= 0 && second.groupId <= 0 {
		newGroupId := first.coordinate.i*1000 + first.coordinate.j
		first.groupId = newGroupId
		group[newGroupId] = append(group[newGroupId], first)
		second.groupId = newGroupId
		group[newGroupId] = append(group[newGroupId], second)
	} else if first.groupId == second.groupId {
		//	do nothing
	} else if first.groupId <= 0 && second.groupId > 0 {
		groupId := second.groupId
		first.groupId = groupId
		group[groupId] = append(group[groupId], first)
	} else if second.groupId <= 0 && first.groupId > 0 {
		groupId := first.groupId
		second.groupId = groupId
		group[groupId] = append(group[groupId], second)
	} else if first.groupId > 0 && second.groupId > 0 {
		leftGroupId := first.groupId
		deleteGroupId := second.groupId
		for _, node := range group[deleteGroupId] {
			node.groupId = leftGroupId
			group[leftGroupId] = append(group[leftGroupId], node)
		}
		delete(group, deleteGroupId)
	}
	return group
}

func connected(a, b *Node) bool {
	if (a.first == b.coordinate || a.second == b.coordinate) && (b.first == a.coordinate || b.second == a.coordinate) {
		return true
	}
	return false
}

func initializeGrid(matrix [][]string) [][]*Node {
	originalLength := len(matrix)
	originalWidth := len(matrix[0])
	var res = emptyGrid(matrix)
	for i := 0; i < len(res); i++ {
		for j := 0; j < len(res[0]); j++ {
			currentPoint := Point{i, j}
			if i == 0 || j == 0 || i == originalLength+1 || j == originalWidth+1 {
				res[i][j] = &Node{"x", currentPoint, emptyPoint, emptyPoint, -1, false}
				continue
			}
			pipeSymbol := matrix[i-1][j-1]
			switch pipeSymbol {
			case "S":
				res[i][j] = &Node{pipeSymbol, currentPoint, emptyPoint, emptyPoint, -1, true}
			case ".":
				res[i][j] = &Node{pipeSymbol, currentPoint, emptyPoint, emptyPoint, -1, false}
			case "-":
				res[i][j] = &Node{pipeSymbol, currentPoint, Point{i, j - 1}, Point{i, j + 1}, -1, false}
			case "|":
				res[i][j] = &Node{pipeSymbol, currentPoint, Point{i - 1, j}, Point{i + 1, j}, -1, false}
			case "L":
				res[i][j] = &Node{pipeSymbol, currentPoint, Point{i - 1, j}, Point{i, j + 1}, -1, false}
			case "J":
				res[i][j] = &Node{pipeSymbol, currentPoint, Point{i - 1, j}, Point{i, j - 1}, -1, false}
			case "7":
				res[i][j] = &Node{pipeSymbol, currentPoint, Point{i + 1, j}, Point{i, j - 1}, -1, false}
			case "F":
				res[i][j] = &Node{pipeSymbol, currentPoint, Point{i + 1, j}, Point{i, j + 1}, -1, false}
			}
		}
	}
	return res
}

func emptyGrid(matrix [][]string) [][]*Node {
	res := make([][]*Node, len(matrix)+2)
	for i := range res {
		res[i] = make([]*Node, len(matrix[0])+2)
	}
	return res
}

func getStringMatrix(lines []string) [][]string {
	var res [][]string
	for _, line := range lines {
		res = append(res, strings.Split(line, ""))
	}
	return res
}

type Node struct {
	symbol        string
	coordinate    Point
	first, second Point
	groupId       int
	start         bool
}

type Point struct {
	i, j int
}

func (node Node) isOpeningNode() bool {
	if node.first.i == -1 || node.second.i == -1 {
		return false
	}
	if node.first.i == node.second.i {
		return false
	}
	currentJ := node.coordinate.j
	if node.first.j > currentJ || node.second.j > currentJ || node.first.j == node.second.j {
		return true
	}
	return false
}

func (node Node) isClosingNode() bool {
	if node.first.i == -1 || node.second.i == -1 {
		return false
	}
	if node.first.i == node.second.i {
		return false
	}
	currentJ := node.coordinate.j
	if node.first.j < currentJ || node.second.j < currentJ || node.first.j == node.second.j {
		return true
	}
	return false
}

func (node Node) isVertical() bool {
	if node.first.i == -1 || node.second.i == -1 {
		return false
	}
	currentI := node.coordinate.i
	if math.Abs(float64(node.first.i-currentI)) == 1 || math.Abs(float64(node.second.i-currentI)) == 1 {
		return true
	}
	return false
}

func getTestLines() (taskLines []string) {
	test := "-L|F7\n7S-7|\nL|7||\n-L-J|\nL|-JF"
	return strings.Split(test, "\n")
}
func getTestLines2() (taskLines []string) {
	test := "7-F7-\n.FJ|7\nSJLL7\n|F--J\nLJ.LJ"
	return strings.Split(test, "\n")
}
func getTestLines3() (taskLines []string) {
	test := "...........\n.S-------7.\n.|F-----7|.\n.||.....||.\n.||.....||.\n.|L-7.F-J|.\n.|..|.|..|.\n.L--J.L--J.\n..........."
	return strings.Split(test, "\n")
}
func getTestLines4() (taskLines []string) {
	test := "..........\n.S------7.\n.|F----7|.\n.||....||.\n.||....||.\n.|L-7F-J|.\n.|..||..|.\n.L--JL--J.\n.........."
	return strings.Split(test, "\n")
}
func getTestLines5() (taskLines []string) {
	test := "FF7FSF7F7F7F7F7F---7\nL|LJ||||||||||||F--J\nFL-7LJLJ||||||LJL-77\nF--JF--7||LJLJ7F7FJ-\nL---JF-JLJ.||-FJLJJ7\n|F|F-JF---7F7-L7L|7|\n|FFJF7L7F-JF7|JL---7\n7-L-JL7||F7|L7F-7F7|\nL.L7LFJ|||||FJL7||LJ\nL7JLJL-JLJLJL--JLJ.L"
	return strings.Split(test, "\n")
}
