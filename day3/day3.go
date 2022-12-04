package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zyedidia/generic/set"
)

func main() {
	part1()
	part2()
}

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	total := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		sackSize := len(line) / 2
		sackA := set.NewMapset[rune]()
		sackB := set.NewMapset[rune]()
		for i, r := range line {
			if i < sackSize {
				sackA.Put(r)
			} else {
				sackB.Put(r)
			}
		}
		common := sackA.Intersection(sackB).Keys()[0]
		total += letterVal(common)
	}
	fmt.Printf("Part 1: Total Score=%d\n", total)
}

func part2() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	total := 0
	groups := [][]set.Set[rune]{}
	g := []set.Set[rune]{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		m := set.NewMapset([]rune(line)...)
		g = append(g, m)
		if len(g) == 3 {
			groups = append(groups, g)
			g = []set.Set[rune]{}
		}
	}
	for _, g := range groups {
		badge := g[0].Intersection(g[1], g[2]).Keys()[0]
		total += letterVal(badge)
	}
	fmt.Printf("Part 2: Total Score=%d\n", total)
}

const (
	a = int('a')
	A = int('A')
)

func letterVal(r rune) int {
	if r >= 'a' && r <= 'z' {
		return int(r) - a + 1
	}
	return int(r) - A + 27
}
