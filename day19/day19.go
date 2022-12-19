package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	part1()
	part2()
}

type PerType[T any] struct {
	ore, clay, obsidian, geode T
}

type State struct {
	cost                  PerType[PerType[int]]
	robotCount, inventory PerType[int]
}

var lineRe = regexp.MustCompile(`Blueprint \d*: Each ore robot costs (?P<orecost>\d+) ore. Each clay robot costs (?P<clayorecost>\d+) ore. Each obsidian robot costs (?P<obsidianorecost>\d+) ore and (?P<obsidianclaycost>\d+) clay. Each geode robot costs (?P<geodeorecost>\d+) ore and (?P<geodeobsidiancost>\d+) obsidian.`)

func reSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}
	return subMatchMap
}

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	total := 0
	blueprints := []PerType[PerType[int]]{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		m := reSubMatchMap(lineRe, line)
		bp := PerType[PerType[int]]{}
		bp.ore.ore, _ = strconv.Atoi(m["orecost"])
		bp.clay.ore, _ = strconv.Atoi(m["clayorecost"])
		bp.obsidian.ore, _ = strconv.Atoi(m["obsidianorecost"])
		bp.obsidian.clay, _ = strconv.Atoi(m["obsidianclaycost"])
		bp.geode.ore, _ = strconv.Atoi(m["geodeorecost"])
		bp.geode.obsidian, _ = strconv.Atoi(m["geodeobsidiancost"])
		blueprints = append(blueprints, bp)
	}
	for i, bp := range blueprints {
		geodes := tickItt(State{
			cost: bp,
			robotCount: PerType[int]{
				ore: 1,
			},
		}, 24)
		total += (i + 1) * geodes
	}
	fmt.Printf("Part 1: %d\n", total)
}

func tickItt(start State, time int) int {
	states := []map[State]bool{
		{
			start: true,
		},
	}
	maxGeodes := 0
	for t := 0; t < time; t++ {
		states = append(states, make(map[State]bool, len(states[t])*5))
		for prevState := range states[t] {
			if prevState.inventory.geode > maxGeodes {
				maxGeodes = prevState.inventory.geode
			}
			if prevState.inventory.geode+prevState.robotCount.geode < maxGeodes {
				continue
			}
			states[t+1][State{
				cost:       prevState.cost,
				robotCount: prevState.robotCount,
				inventory: PerType[int]{
					ore:      prevState.inventory.ore + prevState.robotCount.ore,
					clay:     prevState.inventory.clay + prevState.robotCount.clay,
					obsidian: prevState.inventory.obsidian + prevState.robotCount.obsidian,
					geode:    prevState.inventory.geode + prevState.robotCount.geode,
				},
			}] = true
			if prevState.inventory.ore >= prevState.cost.geode.ore && prevState.inventory.obsidian >= prevState.cost.geode.obsidian {
				states[t+1][State{
					cost: prevState.cost,
					robotCount: PerType[int]{
						ore:      prevState.robotCount.ore,
						clay:     prevState.robotCount.clay,
						obsidian: prevState.robotCount.obsidian,
						geode:    prevState.robotCount.geode + 1,
					},
					inventory: PerType[int]{
						ore:      prevState.robotCount.ore + prevState.inventory.ore - prevState.cost.geode.ore,
						clay:     prevState.robotCount.clay + prevState.inventory.clay,
						obsidian: prevState.robotCount.obsidian + prevState.inventory.obsidian - prevState.cost.geode.obsidian,
						geode:    prevState.robotCount.geode + prevState.inventory.geode,
					},
				}] = true
			}
			if prevState.inventory.ore >= prevState.cost.obsidian.ore && prevState.inventory.clay >= prevState.cost.obsidian.clay {
				states[t+1][State{
					cost: prevState.cost,
					robotCount: PerType[int]{
						ore:      prevState.robotCount.ore,
						clay:     prevState.robotCount.clay,
						obsidian: prevState.robotCount.obsidian + 1,
						geode:    prevState.robotCount.geode,
					},
					inventory: PerType[int]{
						ore:      prevState.robotCount.ore + prevState.inventory.ore - prevState.cost.obsidian.ore,
						clay:     prevState.robotCount.clay + prevState.inventory.clay - prevState.cost.obsidian.clay,
						obsidian: prevState.robotCount.obsidian + prevState.inventory.obsidian,
						geode:    prevState.robotCount.geode + prevState.inventory.geode,
					},
				}] = true
			}
			if prevState.inventory.ore >= prevState.cost.clay.ore {
				states[t+1][State{
					cost: prevState.cost,
					robotCount: PerType[int]{
						ore:      prevState.robotCount.ore,
						clay:     prevState.robotCount.clay + 1,
						obsidian: prevState.robotCount.obsidian,
						geode:    prevState.robotCount.geode,
					},
					inventory: PerType[int]{
						ore:      prevState.robotCount.ore + prevState.inventory.ore - prevState.cost.clay.ore,
						clay:     prevState.robotCount.clay + prevState.inventory.clay,
						obsidian: prevState.robotCount.obsidian + prevState.inventory.obsidian,
						geode:    prevState.robotCount.geode + prevState.inventory.geode,
					},
				}] = true
			}
			if prevState.inventory.ore >= prevState.cost.ore.ore {
				states[t+1][State{
					cost: prevState.cost,
					robotCount: PerType[int]{
						ore:      prevState.robotCount.ore + 1,
						clay:     prevState.robotCount.clay,
						obsidian: prevState.robotCount.obsidian,
						geode:    prevState.robotCount.geode,
					},
					inventory: PerType[int]{
						ore:      prevState.robotCount.ore + prevState.inventory.ore - prevState.cost.ore.ore,
						clay:     prevState.robotCount.clay + prevState.inventory.clay,
						obsidian: prevState.robotCount.obsidian + prevState.inventory.obsidian,
						geode:    prevState.robotCount.geode + prevState.inventory.geode,
					},
				}] = true
			}
		}
	}
	max := 0
	for s, _ := range states[time] {
		if s.inventory.geode > max {
			max = s.inventory.geode
		}
	}
	return max
}

func part2() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	blueprints := []PerType[PerType[int]]{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		m := reSubMatchMap(lineRe, line)
		bp := PerType[PerType[int]]{}
		bp.ore.ore, _ = strconv.Atoi(m["orecost"])
		bp.clay.ore, _ = strconv.Atoi(m["clayorecost"])
		bp.obsidian.ore, _ = strconv.Atoi(m["obsidianorecost"])
		bp.obsidian.clay, _ = strconv.Atoi(m["obsidianclaycost"])
		bp.geode.ore, _ = strconv.Atoi(m["geodeorecost"])
		bp.geode.obsidian, _ = strconv.Atoi(m["geodeobsidiancost"])
		blueprints = append(blueprints, bp)
	}
	ch1, ch2, ch3 := make(chan int), make(chan int), make(chan int)
	go checkBpAsync(blueprints[0], ch1)
	go checkBpAsync(blueprints[1], ch2)
	go checkBpAsync(blueprints[2], ch3)
	geode0, geode1, geode2 := <-ch1, <-ch2, <-ch3
	fmt.Printf("Part 2: %d\n", geode0*geode1*geode2)
}

func checkBpAsync(bp PerType[PerType[int]], out chan int) {
	out <- tickItt(State{
		cost: bp,
		robotCount: PerType[int]{
			ore: 1,
		},
	}, 32)
}
