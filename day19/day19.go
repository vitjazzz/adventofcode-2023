package day19

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"strconv"
	"strings"
)

var workflows = make(map[string]workflow)

const MAX_RANGE = 4000

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/19/input", true)
	//taskLines := getTestLines()
	//parts := getParts(taskLines)
	workflows = getWorkflows(taskLines)
	initialRange := valRange{1, MAX_RANGE}
	initialGroup := partsGroup{map[string]valRange{"x": initialRange, "m": initialRange, "a": initialRange, "s": initialRange}}
	variants := calculateVariants(initialGroup, "in")
	fmt.Printf("res = %d\n", variants)
}

//167409079868000
//140809783868000

func calculateVariants(group partsGroup, workflowId string) int64 {
	wf := workflows[workflowId]
	res := int64(0)
	for _, p := range wf.processGroup(group) {
		key := p.workflowId
		val := p.group
		if key == "A" {
			res += getVariants(val)
			//printGroup(val, workflowId, getVariants(val))
		} else if key == "R" {
			continue
		} else {
			res += calculateVariants(val, key)
		}
	}
	return res
}

func printGroup(group partsGroup, wfId string, variants int64) {
	fmt.Printf("x: [%d, %d], m: [%d, %d], a: [%d, %d], s: [%d, %d], wfId - %s, variants - %d\n",
		group.attrs["x"].from, group.attrs["x"].to, group.attrs["m"].from, group.attrs["m"].to, group.attrs["a"].from, group.attrs["a"].to,
		group.attrs["s"].from, group.attrs["s"].to, wfId, variants)
}

func getVariants(group partsGroup) int64 {
	res := int64(1)
	for _, vr := range group.attrs {
		res *= int64((vr.to - vr.from) + 1)
	}
	return res
}

func getParts(lines []string) []part {
	var res []part
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if line[0] == '{' {
			attributes := strings.Split(line[1:len(line)-1], ",")
			x, _ := strconv.Atoi(attributes[0][2:])
			m, _ := strconv.Atoi(attributes[1][2:])
			a, _ := strconv.Atoi(attributes[2][2:])
			s, _ := strconv.Atoi(attributes[3][2:])
			attributesMap := map[string]int{"x": x, "m": m, "a": a, "s": s}
			res = append(res, part{attributesMap, x, m, a, s})
		}
	}
	return res
}

func getWorkflows(lines []string) map[string]workflow {
	res := make(map[string]workflow)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if line[0] != '{' {
			workflowId := strings.Split(line, "{")[0]
			var rules []rule
			rulesStr := strings.Split(line[len(workflowId)+1:len(line)-1], ",")
			for i, ruleStr := range rulesStr {
				if i == len(rulesStr)-1 {
					rules = append(rules, rule{"", "", -1, ruleStr})
					continue
				}
				attribute := string(ruleStr[0])
				operation := string(ruleStr[1])
				value, _ := strconv.Atoi(strings.Split(ruleStr[2:], ":")[0])
				resultWorkflowId := strings.Split(ruleStr[2:], ":")[1]
				rules = append(rules, rule{attribute, operation, value, resultWorkflowId})
			}
			res[workflowId] = workflow{workflowId, rules}
		}
	}
	return res
}

type workflow struct {
	id    string
	rules []rule
}

type rule struct {
	attribute        string
	operation        string
	value            int
	resultWorkflowId string
}

func (w workflow) process(p part) string {
	for i, r := range w.rules {
		if i == len(w.rules)-1 {
			return r.resultWorkflowId
		}
		switch r.operation {
		case ">":
			if p.attributes[r.attribute] > r.value {
				return r.resultWorkflowId
			}
		case "<":
			if p.attributes[r.attribute] < r.value {
				return r.resultWorkflowId
			}
		}
	}
	return "UNKNOWN"
}

func (w workflow) processGroup(pg partsGroup) []pair {
	currentGroup := pg
	var res []pair
	for i, r := range w.rules {
		if i == len(w.rules)-1 {
			res = append(res, pair{r.resultWorkflowId, currentGroup})
			break
		}
		switch r.operation {
		case ">":
			currentRange := currentGroup.attrs[r.attribute]
			if currentRange.to <= r.value {
				continue
			} else if currentRange.from > r.value {
				res = append(res, pair{r.resultWorkflowId, currentGroup})
				break
			} else {
				newRanges := make(map[string]valRange, len(currentGroup.attrs))
				for k, v := range currentGroup.attrs {
					newRanges[k] = v
				}
				newRange := valRange{r.value + 1, currentRange.to}
				newRanges[r.attribute] = newRange
				newGroup := partsGroup{newRanges}
				res = append(res, pair{r.resultWorkflowId, newGroup})
				currentRange = valRange{currentRange.from, r.value}
				currentRanges := currentGroup.attrs
				currentRanges[r.attribute] = currentRange
				currentGroup = partsGroup{currentRanges}
			}
		case "<":
			currentRange := currentGroup.attrs[r.attribute]
			if currentRange.from >= r.value {
				continue
			} else if currentRange.to < r.value {
				res = append(res, pair{r.resultWorkflowId, currentGroup})
				break
			} else {
				newRanges := make(map[string]valRange, len(currentGroup.attrs))
				for k, v := range currentGroup.attrs {
					newRanges[k] = v
				}
				newRange := valRange{currentRange.from, r.value - 1}
				newRanges[r.attribute] = newRange
				newGroup := partsGroup{newRanges}
				res = append(res, pair{r.resultWorkflowId, newGroup})
				currentRange = valRange{r.value, currentRange.to}
				currentRanges := currentGroup.attrs
				currentRanges[r.attribute] = currentRange
				currentGroup = partsGroup{currentRanges}
			}
		}
	}
	return res
}

type part struct {
	attributes map[string]int
	x, m, a, s int
}

type partsGroup struct {
	attrs map[string]valRange
}

type valRange struct {
	from, to int
}

type pair struct {
	workflowId string
	group      partsGroup
}

func getTestLines() (taskLines []string) {
	test := "px{a<2006:qkq,m>2090:A,rfg}\npv{a>1716:R,A}\nlnx{m>1548:A,A}\nrfg{s<537:gd,x>2440:R,A}\nqs{s>3448:A,lnx}\nqkq{x<1416:A,crn}\ncrn{x>2662:A,R}\nin{s<1351:px,qqz}\nqqz{s>2770:qs,m<1801:hdj,R}\ngd{a>3333:R,R}\nhdj{m>838:A,pv}\n\n{x=787,m=2655,a=1222,s=2876}\n{x=1679,m=44,a=2067,s=496}\n{x=2036,m=264,a=79,s=2244}\n{x=2461,m=1339,a=466,s=291}\n{x=2127,m=1623,a=2188,s=1013}"
	return strings.Split(test, "\n")
}
