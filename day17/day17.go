package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

type Point struct {
	x, y int64
}

type Rock struct {
	loc     Point
	relPoss []Point
}

type FillType int

const (
	WallType FillType = iota
	RockType
	FallingType
)

var rockHeights []int64 = []int64{
	1,
	3,
	3,
	4,
	2,
}

var rocksb [][]int = [][]int{
	{
		0b000111100,
	},
	{
		0b00010000,
		0b00111000,
		0b00010000,
	},
	{
		0b000111000,
		0b000001000,
		0b000001000,
	},
	{
		0b000100000,
		0b000100000,
		0b000100000,
		0b000100000,
	},
	{
		0b000110000,
		0b000110000,
	},
}

var rockTypes = [][]Point{
	{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
	},
	{
		{1, 0},
		{0, 1},
		{1, 1},
		{2, 1},
		{1, 2},
	},
	{
		{2, 1},
		{2, 2},
		{2, 0},
		{1, 0},
		{0, 0},
	},
	{
		{0, 0},
		{0, 1},
		{0, 2},
		{0, 3},
	},
	{
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},
	},
}

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
	var arrows []rune
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			arrows = []rune(line)
		}
	}
	filled := map[Point]FillType{
		{0, 0}: WallType,
		{1, 0}: WallType,
		{2, 0}: WallType,
		{3, 0}: WallType,
		{4, 0}: WallType,
		{5, 0}: WallType,
		{6, 0}: WallType,
	}
	maxHeight := int64(0)
	arrowi := 0
	for rnum := 0; rnum < 2022; rnum++ {
		rheight := rockHeights[rnum%len(rockHeights)]
		for wheight := maxHeight; wheight < maxHeight+rheight+4; wheight++ {
			filled[Point{-1, wheight}] = WallType
			filled[Point{7, wheight}] = WallType
		}
		rock := Rock{
			loc:     Point{2, maxHeight + 4},
			relPoss: append([]Point(nil), rockTypes[rnum%len(rockTypes)]...),
		}
		//printFilledAndRock(rock, filled)
		for {
			arrow := arrows[arrowi%len(arrows)]
			arrowi++
			if arrow == '>' {
				right := Point{rock.loc.x + 1, rock.loc.y}
				if isValid(right, rock, filled) {
					rock.loc = right
				}
			} else if arrow == '<' {
				left := Point{rock.loc.x - 1, rock.loc.y}
				if isValid(left, rock, filled) {
					rock.loc = left
				}
			}
			down := Point{rock.loc.x, rock.loc.y - 1}
			if isValid(down, rock, filled) {
				rock.loc = down
			} else {
				break
			}
		}
		for _, rpos := range rock.relPoss {
			y := rock.loc.y + rpos.y
			if y > maxHeight {
				maxHeight = y
			}
			filled[Point{rock.loc.x + rpos.x, y}] = RockType
		}
	}
	fmt.Printf("Part 1: %d\n", maxHeight)
}

func isValid(p Point, rock Rock, filled map[Point]FillType) bool {
	for _, rpos := range rock.relPoss {
		if _, ok := filled[Point{p.x + rpos.x, p.y + rpos.y}]; ok {
			return false
		}
	}
	return true
}

func printFilledAndRock(rock Rock, filled map[Point]FillType) {
	poss := map[Point]FillType{}
	maxy := int64(0)
	for p, t := range filled {
		if p.y > maxy {
			maxy = p.y
		}
		poss[p] = t
	}
	for _, p := range rock.relPoss {
		poss[Point{p.x + rock.loc.x, p.y + rock.loc.y}] = FallingType
	}
	for y := maxy; y >= 0; y-- {
		for x := int64(-1); x < 8; x++ {
			t, ok := poss[Point{x, y}]
			if ok {
				switch t {
				case RockType:
					fmt.Print("#")
				case WallType:
					fmt.Print("|")
				case FallingType:
					fmt.Print("@")
				default:
					fmt.Print(".")
				}
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

type CacheVal struct {
	rnum,
	heights []int64
}

func part2() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var arrows []rune
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			arrows = []rune(line)
		}
	}
	arrowi := 0
	rows := []int{}
	rows = append(rows, 0b111111111)
	trimmed := 0
	mem := map[string]*CacheVal{}
	target := int64(1000000000000) - 1
	for rnum := 0; rnum < 100000000; rnum++ {
		rocki := rnum % len(rocksb)
		k := genkey(rows, rocki, arrowi%len(arrows))
		if n, ok := mem[k]; ok && len(n.rnum) > 2 && (target-int64(n.rnum[0]))%int64(n.rnum[1]-n.rnum[0]) == 0 {
			var ndelta, hdelta []int64
			for i := 0; i < len(n.rnum)-2; i++ {
				ndelta = append(ndelta, int64(n.rnum[i+1]-n.rnum[i]))
				hdelta = append(hdelta, int64(n.heights[i+1]-n.heights[i]))
			}
			curr := n.rnum[0]
			height := n.heights[0]
			for curr != target {
				curr += ndelta[0]
				height += hdelta[0]
			}
			fmt.Printf("Part 2: %d\n", height)
			return
		}
		rock := append([]int(nil), rocksb[rocki]...)
		for i := 0; i < len(rock)+4; i++ {
			rows = append(rows, 0b100000001)
		}
		offset := len(rows) - 1 - len(rock)
		for {
			arrow := arrows[arrowi%len(arrows)]
			arrowi++
			if arrow == '>' {
				rightRock := make([]int, len(rock))
				for i, v := range rock {
					rightRock[i] = v >> 1
				}
				if isValidb(rightRock, rows, offset) {
					rock = rightRock
				}
			} else if arrow == '<' {
				rockLeft := make([]int, len(rock))
				for i, v := range rock {
					rockLeft[i] = v << 1
				}
				if isValidb(rockLeft, rows, offset) {
					rock = rockLeft
				}
			}
			if isValidb(rock, rows, offset-1) {
				offset--
			} else {
				break
			}
		}
		for i, rockLine := range rock {
			rows[offset+i] |= rockLine
		}
		for i := len(rows); i > 0; i-- {
			if rows[i-1] != 0b100000001 {
				rows = rows[:i]
				break
			}
		}
		for i := len(rows) - 1; i >= 0; i-- {
			if rows[i] == 0b111111111 {
				trimmed += i
				rows = rows[i:]
				break
			}
		}
		if m, ok := mem[k]; ok {
			m.rnum = append(m.rnum, int64(rnum))
			m.heights = append(m.heights, int64(len(rows)-1+trimmed))
		} else {
			mem[k] = &CacheVal{
				rnum:    []int64{int64(rnum)},
				heights: []int64{int64(len(rows) - 1 + trimmed)},
			}
		}
	}
	fmt.Printf("Part 2: %d\n", len(rows)-1+trimmed)
}
func hasBit(n int, pos uint) bool {
	val := n & (int(1) << pos)
	return (val > 0)
}

func isValidb(rock []int, rows []int, off int) bool {
	for i, rockLine := range rock {
		if rows[off+i]&rockLine != 0 {
			return false
		}
	}
	return true
}

func genkey(rows []int, rocki int, arrowi int) string {
	var sb bytes.Buffer
	for _, row := range rows {
		sb.WriteRune(rune(row))
		sb.WriteRune(',')
	}
	sb.WriteRune(':')
	sb.WriteRune(rune(rocki))
	sb.WriteRune(':')
	sb.WriteRune(rune(arrowi))
	return sb.String()
}

func printFilledb(rows []int, rock []int, offset int) {
	for i := len(rows) - 1; i >= 0; i-- {
		for pos := 8; pos >= 0; pos-- {
			if hasBit(rows[i], uint(pos)) {
				fmt.Print("#")
			} else if i >= offset && i < offset+len(rock) && hasBit(rock[i-offset], uint(pos)) {
				fmt.Print("@")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}
