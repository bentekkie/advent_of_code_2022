package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
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
	var buf, towerLines []string
	inMoves := false
	for scanner.Scan() {
		line := scanner.Text()
		if !inMoves && line == "" {
			towerLines = append(towerLines, buf...)
			buf = []string{}
			inMoves = true
		}
		buf = append(buf, line)
	}
	t := newTower(towerLines)
	for _, m := range parseMoves(buf[1:]) {
		t.ApplyMovePart1(m)
	}
	fmt.Printf("Part 1: %s\n", t.Message())
}

func part2() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var buf, towerLines []string
	inMoves := false
	for scanner.Scan() {
		line := scanner.Text()
		if !inMoves && line == "" {
			towerLines = append(towerLines, buf...)
			buf = []string{}
			inMoves = true
		}
		buf = append(buf, line)
	}
	t := newTower(towerLines)
	for _, m := range parseMoves(buf[1:]) {
		t.ApplyMovePart2(m)
	}
	fmt.Printf("Part 2: %s\n", t.Message())
}

func getParams(regEx *regexp.Regexp, line string) (paramsMap map[string]string) {
	match := regEx.FindStringSubmatch(line)

	paramsMap = make(map[string]string)
	for i, name := range regEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}

type tower struct {
	m map[rune][]rune
}

func newTower(lines []string) *tower {
	var tr [][]rune
	for _, line := range lines {
		tr = append(tr, []rune(line))
	}
	tr = transpose(tr)
	t := &tower{
		m: map[rune][]rune{},
	}
	for _, row := range tr {
		line := strings.TrimSpace(reverse(row))
		if len(line) > 0 && !strings.ContainsAny(line, "[]") {
			lr := []rune(line)
			t.m[lr[0]] = lr[1:]
		}
	}
	return t
}

var mvr = regexp.MustCompile(`.*move (?P<num>\d*) from (?P<from>\d) to (?P<to>\d)`)

func parseMoves(lines []string) []*move {
	var moves []*move
	for _, m := range lines {
		ps := getParams(mvr, m)
		n, _ := strconv.Atoi(ps["num"])
		moves = append(moves, &move{
			n:    n,
			from: []rune(ps["from"])[0],
			to:   []rune(ps["to"])[0],
		})
	}
	return moves
}

type move struct {
	n        int
	from, to rune
}

func (t *tower) String() string {
	s := ""
	for _, k := range t.Keys() {
		s += fmt.Sprintf("%s: %s\n", string(k), string(t.m[k]))
	}
	return s
}

func (t *tower) Message() string {
	msg := []rune{}
	for _, k := range t.Keys() {
		msg = append(msg, t.m[k][len(t.m[k])-1])
	}
	return string(msg)
}

func (t *tower) Keys() []rune {
	keys := []rune{}
	for k := range t.m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func (t *tower) ApplyMovePart1(mv *move) {
	for i := 0; i < mv.n; i++ {
		r := t.m[mv.from][len(t.m[mv.from])-1]
		t.m[mv.to] = append(t.m[mv.to], r)
		t.m[mv.from] = t.m[mv.from][:len(t.m[mv.from])-1]
	}
}

func (t *tower) ApplyMovePart2(mv *move) {
	rs := t.m[mv.from][len(t.m[mv.from])-mv.n : len(t.m[mv.from])]
	t.m[mv.to] = append(t.m[mv.to], rs...)
	t.m[mv.from] = t.m[mv.from][:len(t.m[mv.from])-mv.n]
}

func reverse(runes []rune) string {
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func transpose[T comparable](slice [][]T) [][]T {
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
