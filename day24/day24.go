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

type Dir int

const (
	Up Dir = iota
	Down
	Left
	Right
)

type Stage int

const (
	ToEnd Stage = iota
	Back
	ToEndAgain
)

type Point struct {
	x, y int
}

var start = Point{-1, 0}

func (p Point) neighbours(maxX, maxY int, end Point) []Point {
	if p == start {
		return []Point{
			p,
			{0, 0},
		}
	}
	if p == end {
		return []Point{
			p,
			{maxX, maxY},
		}
	}
	out := []Point{
		p,
	}
	if p.x > 0 {
		out = append(out, Point{p.x - 1, p.y})
	}
	if p.y > 0 {
		out = append(out, Point{p.x, p.y - 1})
	}
	if p.x < maxX {
		out = append(out, Point{p.x + 1, p.y})
	}
	if p.y < maxY {
		out = append(out, Point{p.x, p.y + 1})
	}
	if p.x == maxX && p.y == maxY {
		out = append(out, end)
	}
	if p.x == 0 && p.y == 0 {
		out = append(out, start)
	}
	return out
}

type Bliz struct {
	Point
	d Dir
	r rune
}

type Square struct {
	up, down, left, right int
}

func (s Square) Total() int {
	return s.down + s.up + s.left + s.right
}

func (s Square) String() string {
	total := s.Total()
	if total == 0 {
		return "."
	}
	if total > 1 {
		return strconv.Itoa(total)
	}
	if s.up == 1 {
		return "^"
	}
	if s.down == 1 {
		return "v"
	}
	if s.left == 1 {
		return "<"
	}
	return ">"
}

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	y := -1
	bizs := []*Bliz{}
	maxX := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line)-3 > maxX {
			maxX = len(line) - 3
		}
		for x, r := range []rune(line)[1:] {
			switch r {
			case '>':
				bizs = append(bizs, &Bliz{
					Point: Point{
						x: x,
						y: y,
					},
					r: r,
					d: Right,
				})
			case '<':
				bizs = append(bizs, &Bliz{
					Point: Point{
						x: x,
						y: y,
					},
					r: r,
					d: Left,
				})
			case '^':
				bizs = append(bizs, &Bliz{
					Point: Point{
						x: x,
						y: y,
					},
					r: r,
					d: Up,
				})
			case 'v':
				bizs = append(bizs, &Bliz{
					Point: Point{
						x: x,
						y: y,
					},
					r: r,
					d: Down,
				})
			}
		}
		y++
	}
	m := map[Point][]*Bliz{}
	for _, b := range bizs {
		m[b.Point] = append(m[b.Point], b)
	}
	maxY := y - 1
	field := make([][]Square, maxY)
	for y := 0; y < maxY; y++ {
		field[y] = make([]Square, maxX+1)
		for x := 0; x <= maxX; x++ {
			bs := m[Point{x, y}]
			//fmt.Printf("%v %v\n", Point{x, y}, bs)
			for _, b := range bs {
				switch b.r {
				case '>':
					field[y][x].right++
				case '<':
					field[y][x].left++
				case '^':
					field[y][x].up++
				case 'v':
					field[y][x].down++
				}
			}
		}
	}
	mes := map[Point]bool{start: true}
	end := Point{x: len(field[0]) - 1, y: len(field)}
	round := 0
	for {
		field = nextField(field)
		newMes := map[Point]bool{}
		for me := range mes {
			if me == end {
				fmt.Printf("Part 1: %d\n", round)
				return
			}
			for _, n := range me.neighbours(len(field[0])-1, len(field)-1, end) {
				if n == start || n == end || field[n.y][n.x].Total() == 0 {
					newMes[n] = true
				}
			}
		}
		round++
		mes = newMes
	}
}

func nextField(field [][]Square) [][]Square {
	next := make([][]Square, len(field))
	for y, row := range field {
		next[y] = make([]Square, len(row))
	}
	for y, row := range field {
		for x, s := range row {
			if x > 0 {
				next[y][x-1].left += s.left
			} else {
				next[y][len(next[y])-1].left += s.left
			}
			if x < len(next[y])-1 {
				next[y][x+1].right += s.right
			} else {
				next[y][0].right += s.right
			}
			if y > 0 {
				next[y-1][x].up += s.up
			} else {
				next[len(next)-1][x].up += s.up
			}
			if y < len(next)-1 {
				next[y+1][x].down += s.down
			} else {
				next[0][x].down += s.down
			}
		}
	}
	return next
}

func part2() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	y := -1
	bizs := []*Bliz{}
	maxX := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line)-3 > maxX {
			maxX = len(line) - 3
		}
		for x, r := range []rune(line)[1:] {
			switch r {
			case '>':
				bizs = append(bizs, &Bliz{
					Point: Point{
						x: x,
						y: y,
					},
					r: r,
					d: Right,
				})
			case '<':
				bizs = append(bizs, &Bliz{
					Point: Point{
						x: x,
						y: y,
					},
					r: r,
					d: Left,
				})
			case '^':
				bizs = append(bizs, &Bliz{
					Point: Point{
						x: x,
						y: y,
					},
					r: r,
					d: Up,
				})
			case 'v':
				bizs = append(bizs, &Bliz{
					Point: Point{
						x: x,
						y: y,
					},
					r: r,
					d: Down,
				})
			}
		}
		y++
	}
	m := map[Point][]*Bliz{}
	for _, b := range bizs {
		m[b.Point] = append(m[b.Point], b)
	}
	maxY := y - 1
	field := make([][]Square, maxY)
	for y := 0; y < maxY; y++ {
		field[y] = make([]Square, maxX+1)
		for x := 0; x <= maxX; x++ {
			bs := m[Point{x, y}]
			//fmt.Printf("%v %v\n", Point{x, y}, bs)
			for _, b := range bs {
				switch b.r {
				case '>':
					field[y][x].right++
				case '<':
					field[y][x].left++
				case '^':
					field[y][x].up++
				case 'v':
					field[y][x].down++
				}
			}
		}
	}

	mes := map[PointStage]bool{{p: start, s: ToEnd}: true}
	end := Point{x: len(field[0]) - 1, y: len(field)}
	round := 0
	toEndRound, backRound := -1, -1
	for {
		field = nextField(field)
		newMes := map[PointStage]bool{}
		for me := range mes {
			s := me.s
			if me.p == end && s == ToEnd {
				if toEndRound == -1 {
					toEndRound = round
				}
				s = Back
			} else if me.p == start && s == Back {
				if backRound == -1 {
					backRound = round
				}
				s = ToEndAgain
			} else if me.p == end && s == ToEndAgain {
				fmt.Printf("Part 2: %d\n", round)
				return
			}
			for _, n := range me.p.neighbours(len(field[0])-1, len(field)-1, end) {
				if n == start || n == end || field[n.y][n.x].Total() == 0 {
					newMes[PointStage{p: n, s: s}] = true
				}
			}
		}
		round++
		mes = newMes
	}
}

type PointStage struct {
	p Point
	s Stage
}
