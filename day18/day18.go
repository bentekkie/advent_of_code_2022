package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

func main() {
	part1()
	part2()
}

type Point struct {
	x, y, z float64
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
	Forward
	Backward
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Left:
		return "Left"
	case Right:
		return "Right"
	case Forward:
		return "Forward"
	case Backward:
		return "Backward"
	}
	return ""
}

func (d Direction) flip() Direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	case Forward:
		return Backward
	case Backward:
		return Forward
	}
	return 0
}

type Face struct {
	dir    Direction
	center Point
}

func (f Face) move(d Direction) Face {
	switch d {
	case Up:
		return f.up()
	case Down:
		return f.down()
	case Left:
		return f.left()
	case Right:
		return f.right()
	case Forward:
		return f.forward()
	case Backward:
		return f.backward()
	}
	return f
}

func (f Face) left() Face {
	return Face{
		dir: f.dir,
		center: Point{
			f.center.x - 1, f.center.y, f.center.z,
		},
	}
}
func (f Face) right() Face {
	return Face{
		dir: f.dir,
		center: Point{
			f.center.x + 1, f.center.y, f.center.z,
		},
	}
}
func (f Face) up() Face {
	return Face{
		dir: f.dir,
		center: Point{
			f.center.x, f.center.y, f.center.z + 1,
		},
	}
}
func (f Face) down() Face {
	return Face{
		dir: f.dir,
		center: Point{
			f.center.x, f.center.y, f.center.z - 1,
		},
	}
}
func (f Face) forward() Face {
	return Face{
		dir: f.dir,
		center: Point{
			f.center.x, f.center.y + 1, f.center.z,
		},
	}
}
func (f Face) backward() Face {
	return Face{
		dir: f.dir,
		center: Point{
			f.center.x, f.center.y - 1, f.center.z,
		},
	}
}

func (f Face) flip() Face {
	return Face{
		center: f.center,
		dir:    f.dir.flip(),
	}
}
func (f Face) neighbours() [][]Face {
	var pReg, pInv Point
	ns := [][]Face{}
	switch f.dir {
	case Up:
		pInv = Point{f.center.x - .5, f.center.y - .5, f.center.z}
		pReg = Point{f.center.x - .5, f.center.y - .5, f.center.z - 1}
	case Down:
		pReg = Point{f.center.x - .5, f.center.y - .5, f.center.z}
		pInv = Point{f.center.x - .5, f.center.y - .5, f.center.z - 1}
	case Right:
		pInv = Point{f.center.x, f.center.y - .5, f.center.z - .5}
		pReg = Point{f.center.x - 1, f.center.y - .5, f.center.z - .5}
	case Left:
		pReg = Point{f.center.x, f.center.y - .5, f.center.z - .5}
		pInv = Point{f.center.x - 1, f.center.y - .5, f.center.z - .5}
	case Forward:
		pInv = Point{f.center.x - .5, f.center.y, f.center.z - .5}
		pReg = Point{f.center.x - .5, f.center.y - 1, f.center.z - .5}
	case Backward:
		pReg = Point{f.center.x - .5, f.center.y, f.center.z - .5}
		pInv = Point{f.center.x - .5, f.center.y - 1, f.center.z - .5}
	}
	regs := map[Direction]Face{}
	invs := map[Direction]Face{}

	for _, fi := range cube(pInv) {
		flipped := fi.flip()
		if fi.dir != f.dir && flipped.dir != f.dir {
			invs[fi.dir] = flipped
		}
	}
	for _, fi := range cube(pReg) {
		flipped := fi.flip()
		if fi.dir != f.dir && flipped.dir != f.dir {
			regs[fi.dir] = fi
		}
	}
	for d, reg := range regs {
		ns = append(ns, []Face{invs[d], f.move(d), reg})
	}
	return ns
}

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	total := 0
	faces := map[Face]bool{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			parts := strings.Split(line, ",")
			x, _ := strconv.ParseFloat(parts[0], 64)
			y, _ := strconv.ParseFloat(parts[1], 64)
			z, _ := strconv.ParseFloat(parts[2], 64)
			for _, f := range cube(Point{
				x, y, z,
			}) {
				faces[f] = true
			}
		}
	}
	for f := range faces {
		if !faces[f.flip()] {
			total++

		}
	}
	fmt.Printf("Part 1: %d\n", total)
}

func cube(p Point) [6]Face {
	return [6]Face{
		{
			dir: Down,
			center: Point{
				x: p.x + .5,
				y: p.y + .5,
				z: p.z,
			},
		},
		{
			dir: Up,
			center: Point{
				x: p.x + .5,
				y: p.y + .5,
				z: p.z + 1,
			},
		},
		{
			dir: Left,
			center: Point{
				x: p.x,
				y: p.y + .5,
				z: p.z + .5,
			},
		},
		{
			dir: Right,
			center: Point{
				x: p.x + 1,
				y: p.y + .5,
				z: p.z + .5,
			},
		},
		{
			dir: Backward,
			center: Point{
				x: p.x + .5,
				y: p.y,
				z: p.z + .5,
			},
		},
		{
			dir: Forward,
			center: Point{
				x: p.x + .5,
				y: p.y + 1,
				z: p.z + .5,
			},
		},
	}
}

func firstIn(m map[Face]bool, l []Face) (Face, bool) {
	for _, f := range l {
		if m[f] {
			return f, true
		}
	}
	return Face{}, false
}

//	0   1   2
//
// 0
//
// 1
//
// 2
func part2() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	//total := 0
	faces := map[Face]bool{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			parts := strings.Split(line, ",")
			x, _ := strconv.ParseFloat(parts[0], 64)
			y, _ := strconv.ParseFloat(parts[1], 64)
			z, _ := strconv.ParseFloat(parts[2], 64)
			for _, f := range cube(Point{
				x, y, z,
			}) {
				faces[f] = true
			}
		}
	}
	surface := map[Face]bool{}
	for f := range faces {
		if !faces[f.flip()] {
			surface[f] = true
		}
	}
	notSeen := map[Face]bool{}
	maps.Copy(notSeen, surface)
	groups := []map[Face]bool{}
	for {
		currGroup := map[Face]bool{}
		var firstFace Face
		var found bool
		for k, v := range notSeen {
			if v {
				firstFace = k
				found = true
			}
		}
		if !found {
			break
		}
		toVisit := []Face{firstFace}
		for len(toVisit) > 0 {
			newToVisit := []Face{}
			for _, f := range toVisit {
				currGroup[f] = true
				notSeen[f] = false
				for _, ns := range f.neighbours() {
					if n, ok := firstIn(surface, ns); ok {
						if notSeen[n] {
							newToVisit = append(newToVisit, n)
						}
					}
				}
			}
			toVisit = newToVisit
		}
		groups = append(groups, currGroup)
	}
	max := 0
	for _, g := range groups {
		if len(g) > max {
			max = len(g)
		}
	}
	fmt.Printf("Part 2:%v\n", max)
}
