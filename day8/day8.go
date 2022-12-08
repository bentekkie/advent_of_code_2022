package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	part1()
	part2()
}

type tree struct {
	height  int
	visible bool
}

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	trees := [][]*tree{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		row := []*tree{}
		for _, c := range strings.Split(line, "") {
			height, _ := strconv.Atoi(c)
			row = append(row, &tree{
				height: height,
			})
		}
		trees = append(trees, row)
	}
	treesR := transpose(trees)
	for i, row := range treesR {
		treesR[i] = reverse(row)
	}
	for rowi := 0; rowi < len(trees); rowi++ {
		checkRow(trees[rowi])
	}
	for rowi := 0; rowi < len(treesR); rowi++ {
		checkRow(treesR[rowi])
	}
	vtotal := 0
	for _, row := range trees {
		for _, t := range row {
			if t.visible {
				vtotal += 1
			}
		}
	}
	fmt.Printf("Part 1: %d\n", vtotal)
}

func checkRow(row []*tree) {
	row[0].visible = true
	max := row[0].height
	for colj := 1; colj < len(row); colj++ {
		if row[colj].height > max {
			row[colj].visible = true
			max = row[colj].height
		}
	}
	row[len(row)-1].visible = true
	max = row[len(row)-1].height
	for colj := len(row) - 2; colj >= 0; colj-- {
		if row[colj].height > max {
			row[colj].visible = true
			max = row[colj].height
		}
	}
}

func reverse[T interface{}](orig []T) []T {
	for i, j := 0, len(orig)-1; i < j; i, j = i+1, j-1 {
		orig[i], orig[j] = orig[j], orig[i]
	}
	return orig
}

func transpose[T interface{}](slice [][]T) [][]T {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]T, xl)
	for i := range result {
		result[i] = make([]T, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}

type treeP2 struct {
	height                int
	up, down, left, right *treeP2
}

func (t *treeP2) GoUp(h int) int {
	if t == nil {
		return 0
	}
	if t.height < h {
		return 1 + t.up.GoUp(h)
	}
	return 1
}

func (t *treeP2) GoDown(h int) int {
	if t == nil {
		return 0
	}
	if t.height < h {
		return 1 + t.down.GoDown(h)
	}
	return 1
}

func (t *treeP2) GoLeft(h int) int {
	if t == nil {
		return 0
	}
	if t.height < h {
		return 1 + t.left.GoLeft(h)
	}
	return 1
}

func (t *treeP2) GoRight(h int) int {
	if t == nil {
		return 0
	}
	if t.height < h {
		return 1 + t.right.GoRight(h)
	}
	return 1
}

func part2() {

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	trees := [][]*treeP2{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		row := []*treeP2{}
		for _, c := range strings.Split(line, "") {
			height, _ := strconv.Atoi(c)
			row = append(row, &treeP2{
				height: height,
			})
		}
		trees = append(trees, row)
	}
	for i := 0; i < len(trees); i++ {
		for j := 0; j < len(trees[i]); j++ {
			if i > 0 {
				trees[i][j].up = trees[i-1][j]
			}
			if i < len(trees)-1 {
				trees[i][j].down = trees[i+1][j]
			}
			if j > 0 {
				trees[i][j].left = trees[i][j-1]
			}
			if j < len(trees[i])-1 {
				trees[i][j].right = trees[i][j+1]
			}
		}
	}
	max := 0
	for i := 0; i < len(trees); i++ {
		for j := 0; j < len(trees[i]); j++ {
			t := trees[i][j]
			h := trees[i][j].height
			score := t.up.GoUp(h) * t.down.GoDown(h) * t.left.GoLeft(h) * t.right.GoRight(h)
			if score > max {
				max = score
			}
		}
	}
	fmt.Printf("Part 2: %d\n", max)
}
