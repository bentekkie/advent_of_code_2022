package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func main() {
	part1()
	part2()
}

type Point struct {
	x, y int
}

type Object int

const (
	Beacon Object = iota
	NoBeacon
	Sensor
)

var lineRe = regexp.MustCompile(`Sensor at x=(?P<sx>-?\d+), y=(?P<sy>-?\d+): closest beacon is at x=(?P<bx>-?\d+), y=(?P<by>-?\d+)`)

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	points := map[Point]Object{}
	targetY := 2000000
	k := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		matches := reSubMatchMap(lineRe, line)
		sx, _ := strconv.Atoi(matches["sx"])
		sy, _ := strconv.Atoi(matches["sy"])
		s := Point{sx, sy}
		bx, _ := strconv.Atoi(matches["bx"])
		by, _ := strconv.Atoi(matches["by"])
		b := Point{bx, by}
		dist := mahattenDist(s, b)
		if sy+dist >= targetY || sy-dist <= targetY {
			xdist := dist - abs(sy-targetY)
			for x := sx - xdist; x <= sx+xdist; x++ {
				points[Point{
					x: x,
					y: targetY,
				}] = NoBeacon
			}
		}
		if by == targetY {
			points[b] = Beacon
		}
		k++
	}
	nobs := 0
	for _, o := range points {
		if o == NoBeacon {
			nobs++
		}
	}
	fmt.Printf("Part 1: %d\n", nobs)
}

type Range struct {
	s, e int
}

type YRange struct {
	y int
	r Range
}

func part2() {
	cmin, cmax := 0, 4000000
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	xranges := map[int][]Range{}
	sensors := map[Point]int{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		matches := reSubMatchMap(lineRe, line)
		sx, _ := strconv.Atoi(matches["sx"])
		sy, _ := strconv.Atoi(matches["sy"])
		s := Point{sx, sy}
		bx, _ := strconv.Atoi(matches["bx"])
		by, _ := strconv.Atoi(matches["by"])
		b := Point{bx, by}
		sensors[s] = mahattenDist(s, b)
	}
	c := make(chan YRange)
	go func() {
		for yr := range c {
			xranges[yr.y] = append(xranges[yr.y], yr.r)
		}
	}()
	var wg sync.WaitGroup
	for s, d := range sensors {
		wg.Add(1)
		go rangesForSensor(s, d, cmin, cmax, c, len(sensors), &wg)
	}
	wg.Wait()
	close(c)

	p := findEmpty(xranges, cmax)
	fmt.Printf("Part 2: %d\n", 4000000*p.x+p.y)

}

func findEmpty(xranges map[int][]Range, maxx int) Point {
	out := make(chan Point, 1)
	for y, rs := range xranges {
		go func(y int, rs []Range) {
			x := 0
			slices.SortFunc(rs, func(a, b Range) bool {
				return a.s < b.s || (a.s == b.s && a.e < b.e)
			})
			if rs[0].s == 0 {
				for _, r := range rs {
					if x > r.e {
						continue
					}
					if x >= r.s {
						x = r.e + 1
					}
					if x < r.s {
						break
					}
				}
			}
			if x <= maxx {
				out <- Point{x, y}
			}
		}(y, rs)
	}
	return <-out
}

func rangesForSensor(s Point, d, cmin, cmax int, c chan YRange, ns int, wg *sync.WaitGroup) {
	ymin, ymax := max(s.y-d, cmin), min(s.y+d, cmax)
	for y := ymin; y <= ymax; y++ {
		xdist := d - abs(s.y-y)
		r := Range{
			s: s.x - xdist,
			e: s.x + xdist,
		}
		if r.s < cmin && r.e < cmin {
			continue
		} else if r.s < cmin {
			r.s = cmin
		}
		if r.s > cmax && r.e > cmax {
			continue
		} else if r.e > cmax {
			r.e = cmax
		}
		c <- YRange{y, r}
	}
	wg.Done()
}

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

func mahattenDist(a, b Point) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func abs[T constraints.Integer | constraints.Float](i T) T {
	if i < 0 {
		return -i
	}
	return i
}

func min[T constraints.Ordered](a T, bs ...T) T {
	m := a
	for _, b := range bs {
		if b < m {
			m = b
		}
	}
	return m
}

func max[T constraints.Ordered](a T, bs ...T) T {
	m := a
	for _, b := range bs {
		if b > m {
			m = b
		}
	}
	return m
}
