package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {
	part1()
	part2()
}

type Elf struct {
	x, y int
}

type Dir int

const (
	North Dir = iota
	South
	East
	West
)

func (e Elf) adj() []Elf {
	return []Elf{
		{e.x - 1, e.y},
		{e.x + 1, e.y},
		{e.x, e.y - 1},
		{e.x, e.y + 1},
		{e.x - 1, e.y - 1},
		{e.x - 1, e.y + 1},
		{e.x + 1, e.y - 1},
		{e.x + 1, e.y + 1},
	}
}
func (e Elf) move(d Dir) (Elf, []Elf) {
	var p Elf
	var check []Elf
	switch d {
	case North:
		check = []Elf{
			{e.x - 1, e.y - 1},
			{e.x, e.y - 1},
			{e.x + 1, e.y - 1},
		}
		p = Elf{e.x, e.y - 1}
	case South:
		check = []Elf{
			{e.x - 1, e.y + 1},
			{e.x, e.y + 1},
			{e.x + 1, e.y + 1},
		}
		p = Elf{e.x, e.y + 1}
	case East:
		check = []Elf{
			{e.x + 1, e.y + 1},
			{e.x + 1, e.y},
			{e.x + 1, e.y - 1},
		}
		p = Elf{e.x + 1, e.y}
	case West:
		check = []Elf{
			{e.x - 1, e.y + 1},
			{e.x - 1, e.y},
			{e.x - 1, e.y - 1},
		}
		p = Elf{e.x - 1, e.y}
	}
	return p, check
}

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	elves := map[Elf]bool{}
	y := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		for x, r := range []rune(line) {
			if r == '#' {
				elves[Elf{x, y}] = true
			}
		}
		y++
	}
	moves := []Dir{North, South, West, East}
	for round := 0; round < 10; round++ {
		proposedLocs := make(map[Elf][]Elf, len(elves))
		for elf := range elves {
			if noneInMap(elves, elf.adj()) {
				continue
			}
			for i := 0; i < len(moves); i++ {
				d := moves[(i+round)%len(moves)]
				newLoc, checks := elf.move(d)
				if noneInMap(elves, checks) {
					proposedLocs[newLoc] = append(proposedLocs[newLoc], elf)
					break
				}
			}
		}
		for newLoc, pelfves := range proposedLocs {
			if len(pelfves) == 1 {
				delete(elves, pelfves[0])
				elves[newLoc] = true
			}
		}
	}
	fmt.Printf("Part 1: %d\n", emptyInRect(elves))
}

func emptyInRect(m map[Elf]bool) int {
	xs := make([]int, 0, len(m))
	ys := make([]int, 0, len(m))
	for e := range m {
		xs = append(xs, e.x)
		ys = append(ys, e.y)
	}
	sort.Ints(xs)
	sort.Ints(ys)
	total := 0
	for y := ys[0]; y <= ys[len(ys)-1]; y++ {
		for x := xs[0]; x <= xs[len(xs)-1]; x++ {
			if !m[Elf{x, y}] {
				total++
			}
		}
	}
	return total
}

func printElves(m map[Elf]bool) {
	xs := make([]int, 0, len(m))
	ys := make([]int, 0, len(m))
	for e := range m {
		xs = append(xs, e.x)
		ys = append(ys, e.y)
	}
	sort.Ints(xs)
	sort.Ints(ys)
	for y := ys[0]; y <= ys[len(ys)-1]; y++ {
		for x := xs[0]; x <= xs[len(xs)-1]; x++ {
			if m[Elf{x, y}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func noneInMap[M ~map[K]V, K comparable, V bool](m M, keys []K) bool {
	for _, k := range keys {
		if m[k] {
			return false
		}
	}
	return true
}

func part2() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	elves := map[Elf]bool{}
	y := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		for x, r := range []rune(line) {
			if r == '#' {
				elves[Elf{x, y}] = true
			}
		}
		y++
	}
	moves := []Dir{North, South, West, East}
	moved := true
	round := 0
	for moved {
		moved = false
		proposedLocs := make(map[Elf][]Elf, len(elves))
		for elf := range elves {
			if noneInMap(elves, elf.adj()) {
				continue
			}
			for i := 0; i < len(moves); i++ {
				d := moves[(i+round)%len(moves)]
				newLoc, checks := elf.move(d)
				if noneInMap(elves, checks) {
					proposedLocs[newLoc] = append(proposedLocs[newLoc], elf)
					break
				}
			}
		}
		for newLoc, pelfves := range proposedLocs {
			if len(pelfves) == 1 {
				delete(elves, pelfves[0])
				elves[newLoc] = true
				moved = true
			}
		}
		round++
	}
	fmt.Printf("Part 2: %d\n", round)
}
