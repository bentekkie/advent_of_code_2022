package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

func main() {
	part1()
	part2()
}

type Point struct {
	x, y int
}

type Material int

const (
	Air Material = iota
	Rock
	Sand
)

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	//total := 0
	grid := map[Point]Material{}
	for scanner.Scan() {
		line := strings.Split(strings.TrimSpace(scanner.Text()), " -> ")
		for i := 0; i < len(line)-1; i++ {
			s := newPoint(line[i])
			e := newPoint(line[i+1])
			for _, p := range pointsInLine(s, e) {
				grid[p] = Rock
			}
		}
	}
	maxY := 0
	for p, _ := range grid {
		if p.y > maxY {
			maxY = p.y
		}
	}
	sandCount := 0
	for {
		sandPos := dropSand(Point{500, 0}, grid, maxY)
		if sandPos.y > maxY {
			break
		}
		grid[sandPos] = Sand
		sandCount++
	}
	/*
		for y := 0; y < maxY+5; y++ {
			for x := 490; x < 510; x++ {
				switch grid[Point{x, y}] {
				case Rock:
					fmt.Print("#")
				case Sand:
					fmt.Print("o")
				default:
					fmt.Print(".")
				}
			}
			fmt.Print("\n")
		}
	*/
	fmt.Printf("Part 1: %d\n", sandCount)
}

func min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
func max[T constraints.Ordered](a, b T) T {
	if a < b {
		return b
	}
	return a
}

func dropSand(sandPos Point, grid map[Point]Material, maxY int) Point {
	if sandPos.y > maxY {
		return sandPos
	}
	down := Point{
		sandPos.x, sandPos.y + 1,
	}
	if m, ok := grid[down]; !ok || m == Air {
		return dropSand(down, grid, maxY)
	}
	downLeft := Point{
		sandPos.x - 1, sandPos.y + 1,
	}
	if m, ok := grid[downLeft]; !ok || m == Air {
		return dropSand(downLeft, grid, maxY)
	}
	downRight := Point{
		sandPos.x + 1, sandPos.y + 1,
	}
	if m, ok := grid[downRight]; !ok || m == Air {
		return dropSand(downRight, grid, maxY)
	}
	return sandPos
}

func newPoint(str string) Point {
	parts := strings.Split(str, ",")
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	return Point{
		x, y,
	}
}

func pointsInLine(s, e Point) []Point {
	var ps []Point
	if s.x == e.x {
		for y := min(s.y, e.y); y <= max(s.y, e.y); y++ {
			ps = append(ps, Point{s.x, y})
		}
	} else {
		for x := min(s.x, e.x); x <= max(s.x, e.x); x++ {
			ps = append(ps, Point{x, s.y})
		}
	}
	return ps
}

func part2() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	//total := 0
	grid := map[Point]Material{}
	for scanner.Scan() {
		line := strings.Split(strings.TrimSpace(scanner.Text()), " -> ")
		for i := 0; i < len(line)-1; i++ {
			s := newPoint(line[i])
			e := newPoint(line[i+1])
			for _, p := range pointsInLine(s, e) {
				grid[p] = Rock
			}
		}
	}
	maxY := 0
	for p, _ := range grid {
		if p.y > maxY {
			maxY = p.y
		}
	}
	minX := -1
	maxX := 0
	for p, _ := range grid {
		if p.x > maxX {
			maxX = p.x
		}
		if minX < 0 || minX > p.x {
			minX = p.x
		}
	}

	floor := maxY + 1
	sandCount := 0
	for {
		sandPos := dropSandWithFloor(Point{500, 0}, grid, floor)
		sandCount++
		grid[sandPos] = Sand
		if sandPos.x == 500 && sandPos.y == 0 {
			break
		}
		/*
			if sandCount%10 == 0 {
				for y := 0; y < floor+2; y++ {
					for x := minX - 5; x < maxX+10; x++ {
						switch grid[Point{x, y}] {
						case Rock:
							fmt.Print("#")
						case Sand:
							fmt.Print("o")
						default:
							fmt.Print(".")
						}
					}
					fmt.Print("\n")
				}
			}
		*/

	}

	fmt.Printf("Part 2: %d\n", sandCount)
}

func dropSandWithFloor(sandPos Point, grid map[Point]Material, floor int) Point {
	moved := true
	for sandPos.y < floor && moved {
		moved = false
		down := Point{
			sandPos.x, sandPos.y + 1,
		}
		downLeft := Point{
			sandPos.x - 1, sandPos.y + 1,
		}
		downRight := Point{
			sandPos.x + 1, sandPos.y + 1,
		}
		if m, ok := grid[down]; !ok || m == Air {
			moved = true
			sandPos = down
		} else if m, ok := grid[downLeft]; !ok || m == Air {
			moved = true
			sandPos = downLeft
		} else if m, ok := grid[downRight]; !ok || m == Air {
			moved = true
			sandPos = downRight
		}
	}
	return sandPos
}
