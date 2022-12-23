package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Facing int

const (
	Right Facing = iota
	Down
	Left
	Up
)

func (f Facing) String() string {
	switch f {
	case Right:
		return "R"
	case Down:
		return "D"
	case Left:
		return "L"
	case Up:
		return "U"
	}
	panic("Invalid facing")
}

func (f Facing) Val() int {
	switch f {
	case Right:
		return 0
	case Down:
		return 1
	case Left:
		return 2
	case Up:
		return 3
	}
	panic("Invalid facing")
}

type Tile struct {
	x, y, face, edge          int
	up, down, left, right     *Tile
	upf, downf, leftf, rightf Facing
	isWall                    bool
	last                      rune
}

var instructionRe = regexp.MustCompile(`[0-9]+`)

func (f Facing) rotate(dir rune) Facing {
	switch f {
	case Right:
		if dir == 'R' {
			return Down
		}
		return Up
	case Left:
		if dir == 'R' {
			return Up
		}
		return Down
	case Up:
		if dir == 'R' {
			return Right
		}
		return Left
	case Down:
		if dir == 'R' {
			return Left
		}
		return Right
	}
	panic("Invalid facing")
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
	tileLines := [][]rune{}
	var instructions string
	maxWidth := 0
	for scanner.Scan() {
		origLine := []rune(scanner.Text())
		if len(origLine) == 0 {
			continue
		}
		if origLine[0] == ' ' || origLine[0] == '.' || origLine[0] == '#' {
			tileLines = append(tileLines, origLine)
			if len(origLine) > maxWidth {
				maxWidth = len(origLine)
			}
		}
		instructions = string(origLine)
	}
	tiles := [][]*Tile{}
	for _, line := range tileLines {
		row := make([]*Tile, maxWidth)
		for i, r := range line {
			if r == '.' {
				row[i] = &Tile{
					upf:    Up,
					leftf:  Left,
					rightf: Right,
					downf:  Down,
				}
			}
			if r == '#' {
				row[i] = &Tile{
					isWall: true,
				}
			}
		}
		tiles = append(tiles, row)
	}
	for y, row := range tiles {
		for x, tile := range row {
			if tile == nil {
				continue
			}
			tile.x, tile.y = x+1, y+1
			right := x + 1
			for tiles[y][mod(right, maxWidth)] == nil {
				right++
			}
			tile.right = tiles[y][mod(right, maxWidth)]

			left := x - 1
			for tiles[y][mod(left, maxWidth)] == nil {
				left--
			}
			tile.left = tiles[y][mod(left, maxWidth)]

			down := y + 1
			for tiles[mod(down, len(tiles))][x] == nil {
				down++
			}
			tile.down = tiles[mod(down, len(tiles))][x]

			up := y - 1
			for tiles[mod(up, len(tiles))][x] == nil {
				up--
			}
			tile.up = tiles[mod(up, len(tiles))][x]
		}
	}
	match := instructionRe.FindAllStringIndex(instructions, -1)
	moves := []Move{}
	for _, idx := range match {
		dist, _ := strconv.Atoi(instructions[idx[0]:idx[1]])
		moves = append(moves, Walk{dist})
		if idx[1] < len(instructions) {
			moves = append(moves, Turn{[]rune(instructions[idx[1] : idx[1]+1])[0]})
		}
	}
	f := Right
	x := 0
	for tiles[0][x] == nil || tiles[0][x].isWall {
		x++
	}
	loc := tiles[0][x]
	for _, move := range moves {
		//fmt.Printf("x=%d y=%d %s %s\n", loc.x, loc.y, f, move)
		loc, f = move.Apply(loc, f)
	}
	fmt.Printf("Part 1: %v\n", 1000*loc.y+4*loc.x+f.Val())
}

type Move interface {
	fmt.Stringer
	Apply(loc *Tile, f Facing) (*Tile, Facing)
}

type Turn struct {
	dir rune
}

func (t Turn) String() string {
	return "Turn{" + string(t.dir) + "}"
}

func (t Turn) Apply(loc *Tile, f Facing) (*Tile, Facing) {
	return loc, f.rotate(t.dir)
}

type Walk struct {
	dist int
}

func (w Walk) String() string {
	return fmt.Sprintf("Walk{%d}", w.dist)
}

func (w Walk) Apply(loc *Tile, f Facing) (*Tile, Facing) {
	curr := loc
	for i := 0; i < w.dist; i++ {
		switch f {
		case Right:
			curr.last = '>'
			if curr.right.isWall {
				return curr, f
			}
			f = curr.rightf
			curr = curr.right
		case Left:
			curr.last = '<'
			if curr.left.isWall {
				return curr, f
			}
			f = curr.leftf
			curr = curr.left
		case Up:
			curr.last = '^'
			if curr.up.isWall {
				return curr, f
			}
			f = curr.upf
			curr = curr.up
		case Down:
			curr.last = 'v'
			if curr.down.isWall {
				return curr, f
			}
			f = curr.downf
			curr = curr.down
		}
	}
	return curr, f
}

func mod(a, b int) int {
	return (a%b + b) % b
}

func part2() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	tileLines := [][]rune{}
	var instructions string
	maxWidth := 0
	for scanner.Scan() {
		origLine := []rune(scanner.Text())
		if len(origLine) == 0 {
			continue
		}
		if origLine[0] == ' ' || origLine[0] == '.' || origLine[0] == '#' {
			tileLines = append(tileLines, origLine)
			if len(origLine) > maxWidth {
				maxWidth = len(origLine)
			}
		}
		instructions = string(origLine)
	}
	tiles := [][]*Tile{}
	cubeFaces := [6][][]*Tile{}
	for y, line := range tileLines {
		row := make([]*Tile, maxWidth)
		for x, r := range line {
			if r == '.' {
				row[x] = &Tile{
					upf:    Up,
					leftf:  Left,
					rightf: Right,
					downf:  Down,
					x:      x + 1,
					y:      y + 1,
					face:   -1,
					edge:   -1,
					last:   '.',
				}
			}
			if r == '#' {
				row[x] = &Tile{
					upf:    Up,
					leftf:  Left,
					rightf: Right,
					downf:  Down,
					isWall: true,
					edge:   -1,
				}
			}
		}
		tiles = append(tiles, row)
	}
	for y := 0; y < 50; y++ {
		cubeFaces[0] = append(cubeFaces[0], tiles[y][50:100])
		for _, t := range tiles[y][50:100] {
			t.face = 0
		}
	}
	for y := 0; y < 50; y++ {
		cubeFaces[1] = append(cubeFaces[1], tiles[y][100:150])
		for _, t := range tiles[y][100:150] {
			t.face = 1
		}
	}
	for y := 50; y < 100; y++ {
		cubeFaces[2] = append(cubeFaces[2], tiles[y][50:100])
		for _, t := range tiles[y][50:100] {
			t.face = 2
		}
	}
	for y := 100; y < 150; y++ {
		cubeFaces[3] = append(cubeFaces[3], tiles[y][0:50])
		for _, t := range tiles[y][0:50] {
			t.face = 3
		}
	}
	for y := 100; y < 150; y++ {
		cubeFaces[4] = append(cubeFaces[4], tiles[y][50:100])
		for _, t := range tiles[y][50:100] {
			t.face = 4
		}
	}
	for y := 150; y < 200; y++ {
		cubeFaces[5] = append(cubeFaces[5], tiles[y][0:50])
		for _, t := range tiles[y][0:50] {
			t.face = 5
		}
	}
	for y, row := range tiles {
		for x, tile := range row {
			if tile == nil {
				continue
			}
			if x < len(tiles[y])-1 {
				tile.right = tiles[y][x+1]
			}

			if x > 0 {
				tile.left = tiles[y][x-1]
			}

			if y < len(tiles)-1 {
				tile.down = tiles[y+1][x]
			}

			if y > 0 {
				tile.up = tiles[y-1][x]
			}
		}
	}
	for i, t := range cubeFaces[0][0] {
		o := cubeFaces[5][i][0]
		t.edge = o.face
		o.edge = t.face
		t.up = o
		t.upf = Right
		o.left = t
		o.leftf = Down
	}
	for i, t := range cubeFaces[1][0] {
		o := cubeFaces[5][len(cubeFaces[5])-1][i]
		t.edge = o.face
		o.edge = t.face
		t.up = o
		t.upf = Up
		o.down = t
		o.downf = Down
	}
	for i, t := range cubeFaces[1][len(cubeFaces[1])-1] {
		o := cubeFaces[2][i][len(cubeFaces[2])-1]
		t.edge = o.face
		o.edge = t.face
		t.down = o
		t.downf = Left
		o.right = t
		o.rightf = Up
	}
	for i, row := range cubeFaces[0] {
		t := row[0]
		o := cubeFaces[3][len(cubeFaces[3])-1-i][0]
		t.edge = o.face
		o.edge = t.face
		t.left = o
		t.leftf = Right
		o.left = t
		o.leftf = Right
	}
	for i, row := range cubeFaces[1] {
		t := row[len(cubeFaces[1])-1]
		o := cubeFaces[4][len(cubeFaces[1])-1-i][len(cubeFaces[1][0])-1]
		t.edge = o.face
		o.edge = t.face
		t.right = o
		t.rightf = Left
		o.right = t
		o.rightf = Left
	}
	for i, row := range cubeFaces[2] {
		t := row[0]
		o := cubeFaces[3][0][i]
		t.edge = o.face
		o.edge = t.face
		t.left = o
		t.leftf = Down
		o.up = t
		o.upf = Right
	}
	for i, t := range cubeFaces[4][len(cubeFaces[4])-1] {
		o := cubeFaces[5][i][len(cubeFaces[5])-1]
		t.edge = o.face
		o.edge = t.face
		t.down = o
		t.downf = Left
		o.right = t
		o.rightf = Up
	}

	match := instructionRe.FindAllStringIndex(instructions, -1)
	moves := []Move{}
	for _, idx := range match {
		dist, _ := strconv.Atoi(instructions[idx[0]:idx[1]])
		moves = append(moves, Walk{dist})
		if idx[1] < len(instructions) {
			moves = append(moves, Turn{[]rune(instructions[idx[1] : idx[1]+1])[0]})
		}
	}
	f := Right
	x := 0
	for tiles[0][x] == nil || tiles[0][x].isWall {
		x++
	}
	loc := tiles[0][x]
	for _, move := range moves {
		loc, f = move.Apply(loc, f)
	}

	fmt.Printf("Part 2: %v\n", 1000*loc.y+4*loc.x+f.Val())

}
