package day20

import (
	"adventofcode-2023/adventutils"
	"fmt"
	"slices"
	"strings"
)

const (
	LOW  pulse = "LOW"
	HIGH       = "HIGH"
)
const (
	ON  status = true
	OFF        = false
)

const INITIAL_MODULE = "broadcaster"
const RUNS = 1000

func Run() {
	taskLines := adventutils.GetFromUrl("https://adventofcode.com/2023/day/20/input", true)
	//taskLines := getTestLines2()
	modules := getModules(taskLines)
	lowTotal := 0
	highTotal := 0
	for i := 0; i < RUNS; i++ {
		low, high := run(modules)
		lowTotal += low
		highTotal += high
	}
	fmt.Printf("res = %d\n", lowTotal*highTotal)
}

func run(modules map[string]module) (low, high int) {
	low = 0
	high = 0
	tasks := []task{{"INPUT", INITIAL_MODULE, LOW}}
	for len(tasks) > 0 {
		t := tasks[0]
		if t.p == HIGH {
			high++
		} else {
			low++
		}
		tasks = tasks[1:]
		mod := modules[t.to]
		if mod == nil {
			continue
		}
		for to, p := range mod.process(t.from, t.p) {
			//fmt.Printf("%s -> %s (%s)\n", mod.getId(), to, p)
			tasks = append(tasks, task{mod.getId(), to, p})
		}
	}
	return low, high
}

func getModules(lines []string) map[string]module {
	res := make(map[string]module)
	for _, line := range lines {
		mod := getModule(line)
		res[mod.getId()] = mod
	}
	fillConjunctions(res)
	return res
}

func fillConjunctions(modules map[string]module) {
	inputs := make(map[string][]string)
	for _, mod := range modules {
		for _, destination := range mod.getDestinations() {
			inputsForDest := inputs[destination]
			if !slices.Contains(inputsForDest, mod.getId()) {
				inputs[destination] = append(inputs[destination], mod.getId())
			}
		}
	}
	for _, mod := range modules {
		if mod.getType() == "Conjunction" {
			conjInputs := inputs[mod.getId()]
			conj, _ := mod.(*conjunction)
			for _, input := range conjInputs {
				conj.state[input] = LOW
			}
		}
	}
}

func getModule(line string) module {
	modId := strings.Split(line, " -> ")[0]
	destinations := strings.Split(strings.Split(line, " -> ")[1], ", ")
	if modId[0] == '%' {
		return &flipFlop{modId[1:], destinations, OFF}
	} else if modId[0] == '&' {
		return &conjunction{modId[1:], destinations, make(map[string]pulse)}
	} else {
		return &broadcast{modId, destinations}
	}
}

type flipFlop struct {
	id           string
	destinations []string
	s            status
}

func (ff *flipFlop) getType() string           { return "FlipFlop" }
func (ff *flipFlop) getId() string             { return ff.id }
func (ff *flipFlop) getDestinations() []string { return ff.destinations }
func (ff *flipFlop) process(_ string, p pulse) map[string]pulse {
	if p == HIGH {
		return make(map[string]pulse)
	}
	if ff.s == ON {
		ff.s = OFF
		return sendPulse(LOW, ff.destinations)
	} else {
		ff.s = ON
		return sendPulse(HIGH, ff.destinations)
	}
}

type broadcast struct {
	id           string
	destinations []string
}

func (b *broadcast) getType() string           { return "Broadcast" }
func (b *broadcast) getId() string             { return b.id }
func (b *broadcast) getDestinations() []string { return b.destinations }
func (b *broadcast) process(_ string, p pulse) map[string]pulse {
	return sendPulse(p, b.destinations)
}

type conjunction struct {
	id           string
	destinations []string
	state        map[string]pulse
}

func (c *conjunction) getType() string           { return "Conjunction" }
func (c *conjunction) getId() string             { return c.id }
func (c *conjunction) getDestinations() []string { return c.destinations }
func (c *conjunction) process(from string, p pulse) map[string]pulse {
	c.state[from] = p
	for _, storedPulse := range c.state {
		if storedPulse == LOW {
			return sendPulse(HIGH, c.destinations)
		}
	}
	return sendPulse(LOW, c.destinations)
}

type module interface {
	getType() string
	getId() string
	getDestinations() []string
	process(from string, p pulse) map[string]pulse
}

type task struct {
	from, to string
	p        pulse
}

type pulse string
type status bool

func sendPulse(p pulse, destinations []string) map[string]pulse {
	res := make(map[string]pulse, len(destinations))
	for _, d := range destinations {
		res[d] = p
	}
	return res
}

func getTestLines() (taskLines []string) {
	test := "broadcaster -> a, b, c\n%a -> b\n%b -> c\n%c -> inv\n&inv -> a"
	return strings.Split(test, "\n")
}
func getTestLines2() (taskLines []string) {
	test := "broadcaster -> a\n%a -> inv, con\n&inv -> b\n%b -> con\n&con -> output"
	return strings.Split(test, "\n")
}
