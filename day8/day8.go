package day8

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"slices"
	"strings"
	"time"
)

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/8/input")
	//taskLines := getTestLines3()
	navigationRules, nodesRepresentations := getRules(taskLines)
	tree := buildTree(nodesRepresentations)
	//steps := findSteps(tree, navigationRules)
	steps := findStepsAdvanced(tree, navigationRules)
	fmt.Printf("Steps to find XXZ - %d\n", steps)
}

func findStepsAdvanced(tree map[string]*Node, navigationRules string) int {
	nodesTransitions := buildNodesTransitions(tree, navigationRules)
	steps := 0
	currentNodes := getStartingNodes(tree)
	var iterations int64
	var totalTimeSpent int64
	for ; ; steps += len(navigationRules) {
		iterationStartTime := time.Now()
		finalNodesIndexes := makeRange(0, len(navigationRules)-1)
		for _, node := range currentNodes {
			nodeTransitions := nodesTransitions[node.name]
			var newFinalNodesIndexes []int
			for _, finalNodeIndex := range nodeTransitions.finalNodesIndexes {
				if slices.Contains(finalNodesIndexes, finalNodeIndex) {
					newFinalNodesIndexes = append(newFinalNodesIndexes, finalNodeIndex)
				}
			}
			finalNodesIndexes = newFinalNodesIndexes
			if len(finalNodesIndexes) == 0 {
				break
			}
		}
		if len(finalNodesIndexes) > 0 {
			return steps + (finalNodesIndexes[0] + 1)
		}
		var newCurrentNodes []*Node
		for _, node := range currentNodes {
			nodeTransitions := nodesTransitions[node.name]
			newCurrentNode := nodeTransitions.transitionsTo[len(navigationRules)-1]
			newCurrentNodes = append(newCurrentNodes, newCurrentNode)
		}
		currentNodes = newCurrentNodes
		iterationEndTime := time.Now()
		iterationTime := iterationEndTime.Sub(iterationStartTime)
		totalTimeSpent += iterationTime.Nanoseconds()
		iterations++
		if iterations%100_000 == 0 {
			fmt.Printf("Avg ms for iteration - %d, current iteration - %v\n", totalTimeSpent/iterations, iterationTime)
		}
	}
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func findStepsAdvancedSlowVersion(tree map[string]*Node, navigationRules string) int {
	steps := 0
	startingNodes := getStartingNodes(tree)
	currentNodes := startingNodes
	finishedNodes := 0
	for ; finishedNodes != len(startingNodes); steps++ {
		finishedNodes = 0
		var newCurrentNodes []*Node
		navigationRule := navigationRules[steps%len(navigationRules)]
		for _, node := range currentNodes {
			var newNode *Node
			if navigationRule == 'R' {
				newNode = node.right
			} else {
				newNode = node.left
			}
			if newNode.name[2] == 'Z' {
				finishedNodes++
			}
			newCurrentNodes = append(newCurrentNodes, newNode)
		}
		currentNodes = newCurrentNodes
	}
	return steps
}

func getStartingNodes(tree map[string]*Node) (res []*Node) {
	for key, node := range tree {
		if key[2] == 'A' {
			res = append(res, node)
		}
	}
	return
}

func findSteps(tree map[string]*Node, navigationRules string) int {
	steps := 0
	for currentNode := tree["AAA"]; currentNode.name != "ZZZ"; steps++ {
		navigationRule := navigationRules[steps%len(navigationRules)]
		if navigationRule == 'R' {
			currentNode = currentNode.right
		} else {
			currentNode = currentNode.left
		}
	}
	return steps
}

func buildTree(nodesReps []NodeRepresentation) map[string]*Node {
	res := initMap(nodesReps)
	for _, nodeRep := range nodesReps {
		currentNode := res[nodeRep.name]
		currentNode.left = res[nodeRep.left]
		currentNode.right = res[nodeRep.right]
		res[nodeRep.name] = currentNode
	}
	return res
}

func buildNodesTransitions(tree map[string]*Node, navigationRules string) map[string]NodeTransitions {
	res := make(map[string]NodeTransitions)
	for _, node := range tree {
		currentNode := node
		var finalNodesIndexes []int
		var transitionsTo []*Node
		for i, navigationRule := range navigationRules {
			if navigationRule == 'R' {
				currentNode = currentNode.right
			} else {
				currentNode = currentNode.left
			}
			transitionsTo = append(transitionsTo, currentNode)
			if currentNode.name[2] == 'Z' {
				finalNodesIndexes = append(finalNodesIndexes, i)
			}
		}
		res[node.name] = NodeTransitions{node.name, finalNodesIndexes, transitionsTo}
	}
	return res
}

func initMap(nodesReps []NodeRepresentation) map[string]*Node {
	res := make(map[string]*Node)
	for _, node := range nodesReps {
		res[node.name] = &Node{name: node.name}
	}
	return res
}

func getRules(lines []string) (navigationRules string, nodes []NodeRepresentation) {
	for _, line := range lines {
		if line == "" {
			continue
		} else if strings.Contains(line, "=") {
			name := line[:3]
			left := line[7:10]
			right := line[12:15]
			nodes = append(nodes, NodeRepresentation{name, left, right})
		} else {
			navigationRules = line
		}
	}
	return
}

type NodeRepresentation struct {
	name, left, right string
}

type Node struct {
	name        string
	left, right *Node
}

type NodeTransitions struct {
	name              string
	finalNodesIndexes []int
	transitionsTo     []*Node
}

func getTestLines() (taskLines []string) {
	test := "RL\n\nAAA = (BBB, CCC)\nBBB = (DDD, EEE)\nCCC = (ZZZ, GGG)\nDDD = (DDD, DDD)\nEEE = (EEE, EEE)\nGGG = (GGG, GGG)\nZZZ = (ZZZ, ZZZ)\n"
	return strings.Split(test, "\n")
}
func getTestLines2() (taskLines []string) {
	test := "LLR\n\nAAA = (BBB, BBB)\nBBB = (AAA, ZZZ)\nZZZ = (ZZZ, ZZZ)\n"
	return strings.Split(test, "\n")
}
func getTestLines3() (taskLines []string) {
	test := "LR\n\n11A = (11B, XXX)\n11B = (XXX, 11Z)\n11Z = (11B, XXX)\n22A = (22B, XXX)\n22B = (22C, 22C)\n22C = (22Z, 22Z)\n22Z = (22B, 22B)\nXXX = (XXX, XXX)"
	return strings.Split(test, "\n")
}
