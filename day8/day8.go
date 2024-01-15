package day8

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"slices"
	"strings"
)

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/8/input", true)
	//taskLines := getTestLines3()
	navigationRules, nodesRepresentations := getRules(taskLines)
	tree := buildTree(nodesRepresentations)
	//steps := findSteps(tree, navigationRules)
	steps := findStepsAdvanced(tree, navigationRules)
	fmt.Printf("Steps to find XXZ - %d\n", steps)
}

func findStepsAdvanced(tree map[string]*Node, navigationRules string) int {
	nodesTransitions := buildNodesTransitions(tree, navigationRules)
	loopsMap := calculateLoops(nodesTransitions, navigationRules)

	// turns out I don't need it
	//calculatePossibleCombinations(loopsMap, nodesTransitions)

	var loopsSizes []int
	for _, loop := range loopsMap {
		loopsSizes = append(loopsSizes, len(loop.nodeNames)-1)
	}
	// Less common multiple
	currentLcm := LCM(loopsSizes)

	return currentLcm * len(navigationRules)
}

func calculateLoops(nodesTransitions map[string]*NodeTransitions, navigationRules string) map[int]Loop {
	loopsMap := make(map[int]Loop)
	currentNodesTransitions := getStartingNodesTransitions(nodesTransitions)
	for !areLoopsCompleted(loopsMap) {
		var newNodesTransitions []*NodeTransitions
		for i, nodeTransitions := range currentNodesTransitions {
			loop := loopsMap[i]
			if !slices.Contains(loop.nodeNames, nodeTransitions.name) {
				loop.nodeNames = append(loop.nodeNames, nodeTransitions.name)
				loopsMap[i] = loop
			} else {
				loop.completed = true
				loopsMap[i] = loop
			}
			newCurrentNodeTransition := nodeTransitions.transitionsTo[len(navigationRules)-1]
			newNodesTransitions = append(newNodesTransitions, newCurrentNodeTransition)
		}
		currentNodesTransitions = newNodesTransitions
	}
	return loopsMap
}

func calculatePossibleCombinations(loopsMap map[int]Loop, nodesTransitions map[string]*NodeTransitions) []Combination {
	var possibleCombinations []Combination
	for i := 0; i < len(loopsMap); i++ {
		var newPossibleCombinations []Combination
		loopNodeNames := loopsMap[i].nodeNames[1:]
		for _, loopNodeName := range loopNodeNames {
			if i == 0 {
				newPossibleCombinations = append(newPossibleCombinations, Combination{[]string{loopNodeName}})
				continue
			}
			for _, possibleCombination := range possibleCombinations {
				successfulIntersection := true
				for _, name := range possibleCombination.groupedNames {
					if !findIntersection(nodesTransitions[loopNodeName], nodesTransitions[name]) {
						successfulIntersection = false
					}
				}
				if successfulIntersection {
					newPossibleCombination := possibleCombination
					newPossibleCombination.groupedNames = append(newPossibleCombination.groupedNames, loopNodeName)
					newPossibleCombinations = append(newPossibleCombinations, newPossibleCombination)
				}
			}
		}
		possibleCombinations = newPossibleCombinations
	}
	var uniquePossibleCombinations []Combination
	var uniqueCombinationNames []string
	for _, combination := range possibleCombinations {
		groupedName := ""
		for _, name := range combination.groupedNames {
			groupedName += name
		}
		if !slices.Contains(uniqueCombinationNames, groupedName) {
			uniqueCombinationNames = append(uniqueCombinationNames, groupedName)
			uniquePossibleCombinations = append(uniquePossibleCombinations, combination)
		}
	}
	return uniquePossibleCombinations
}

func LCM(loopsSizes []int) int {
	currentGcd := gcd(loopsSizes[0], loopsSizes[1])
	for _, loopSize := range loopsSizes {
		currentGcd = gcd(currentGcd, loopSize)
	}

	currentLcm := 1
	for _, loopSize := range loopsSizes {
		currentLcm *= loopSize
	}
	currentLcm = currentLcm / currentGcd
	return currentLcm
}

func findIntersection(a, b *NodeTransitions) bool {
	for _, index := range a.finalNodesIndexes {
		if slices.Contains(b.finalNodesIndexes, index) {
			return true
		}
	}
	return false
}

func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func areLoopsCompleted(loops map[int]Loop) bool {
	if len(loops) == 0 {
		return false
	}
	for _, loop := range loops {
		if !loop.completed {
			return false
		}
	}
	return true
}
func getStartingNodesTransitions(nodesTransitions map[string]*NodeTransitions) (res []*NodeTransitions) {
	for key, node := range nodesTransitions {
		if key[2] == 'A' {
			res = append(res, node)
		}
	}
	return
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

func buildNodesTransitions(tree map[string]*Node, navigationRules string) map[string]*NodeTransitions {
	res := make(map[string]*NodeTransitions)
	for _, node := range tree {
		res[node.name] = &NodeTransitions{name: node.name}
	}
	for _, node := range tree {
		currentNode := node
		var finalNodesIndexes []int
		var transitionsTo []*NodeTransitions
		for i, navigationRule := range navigationRules {
			if navigationRule == 'R' {
				currentNode = currentNode.right
			} else {
				currentNode = currentNode.left
			}
			transitionsTo = append(transitionsTo, res[currentNode.name])
			if currentNode.name[2] == 'Z' {
				finalNodesIndexes = append(finalNodesIndexes, i)
			}
		}
		res[node.name].finalNodesIndexes = finalNodesIndexes
		res[node.name].transitionsTo = transitionsTo
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

type Loop struct {
	completed bool
	nodeNames []string
}

type NodeTransitions struct {
	name              string
	finalNodesIndexes []int
	transitionsTo     []*NodeTransitions
}

type Combination struct {
	groupedNames []string
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
