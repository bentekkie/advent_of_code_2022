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

type rng struct {
	low, high int
}

func newRng(line string) *rng {
	parts := strings.Split(line, "-")
	low, _ := strconv.Atoi(parts[0])
	high, _ := strconv.Atoi(parts[1])
	return &rng{
		low, high,
	}
}

func (r *rng) contains(o *rng) bool {
	return r.low <= o.low && r.high >= o.high
}

func (r *rng) overlaps(o *rng) bool {
	return r.low <= o.high && o.low <= r.high
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
		parts := strings.Split(line, ",")
		a := newRng(parts[0])
		b := newRng(parts[1])
		if a.contains(b) || b.contains(a) {
			total += 1
		}
	}
	fmt.Printf("Part 1: %d\n", total)
}

func part2() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	total := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Split(line, ",")
		a := newRng(parts[0])
		b := newRng(parts[1])
		if a.overlaps(b) {
			total += 1
		}
	}
	fmt.Printf("Part 2: %d\n", total)
}
