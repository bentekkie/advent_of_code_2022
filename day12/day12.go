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
	g := parseGrid("input.txt")
	part1(g)
	part2(g)
}

type posNode struct {
	row, col int
}

func (p *posNode) up() posNode {
	return posNode{
		row: p.row - 1,
		col: p.col,
	}
}
func (p *posNode) down() posNode {
	return posNode{
		row: p.row + 1,
		col: p.col,
	}
}
func (p *posNode) left() posNode {
	return posNode{
		row: p.row,
		col: p.col - 1,
	}
}
func (p *posNode) right() posNode {
	return posNode{
		row: p.row,
		col: p.col + 1,
	}
}

type grid struct {
	start, end posNode
	poss       map[posNode]rune
	runeToPoss map[rune][]posNode
}

func (g *grid) neighbours(p posNode) []posNode {
	ns := []posNode{}
	currR := g.poss[p]
	if n, ok := g.poss[p.up()]; ok && n <= currR+1 {
		ns = append(ns, p.up())
	}
	if n, ok := g.poss[p.down()]; ok && n <= currR+1 {
		ns = append(ns, p.down())
	}
	if n, ok := g.poss[p.left()]; ok && n <= currR+1 {
		ns = append(ns, p.left())
	}
	if n, ok := g.poss[p.right()]; ok && n <= currR+1 {
		ns = append(ns, p.right())
	}
	return ns
}

func parseGrid(path string) *grid {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	g := &grid{
		poss:       map[posNode]rune{},
		runeToPoss: map[rune][]posNode{},
	}
	row := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		for col, r := range []rune(line) {
			p := posNode{row: row, col: col}
			g.poss[p] = r
			if r == 'S' {
				g.start = p
				g.poss[p] = 'a'
			} else if r == 'E' {
				g.end = p
				g.poss[p] = 'z'
			}
			g.runeToPoss[g.poss[p]] = append(g.runeToPoss[g.poss[p]], p)
		}
		row += 1
	}
	return g
}

func (g *grid) findDistFrom(s posNode, max int) int {
	seen := map[posNode]bool{}
	toVisit := []posNode{s}
	dist := 0
	for !seen[g.end] {
		if max > 0 && dist > max {
			return -1
		}
		dist += 1
		next := []posNode{}
		for _, p := range toVisit {
			if !seen[p] {
				next = append(next, g.neighbours(p)...)
				seen[p] = true
			}
		}
		toVisit = next
	}
	return dist - 1
}

func part1(g *grid) {
	fmt.Printf("Part 1: %d\n", g.findDistFrom(g.start, -1))
}

func part2(g *grid) {
	dists := []int{}
	max := g.findDistFrom(g.start, -1)
	for _, p := range g.runeToPoss['a'] {
		if d := g.findDistFrom(p, max); d != -1 {
			dists = append(dists, d)
		}
	}
	sort.Ints(dists)
	fmt.Printf("Part 2: %d\n", dists[0])
}
