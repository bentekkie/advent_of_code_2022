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

type loc struct {
	x, y int
}

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	head := loc{0, 0}
	tail := loc{0, 0}
	locs := map[loc]bool{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Split(line, " ")
		dir := parts[0]
		n, _ := strconv.Atoi(parts[1])
		var f func(loc) loc
		switch dir {
		case "R":
			f = right
		case "L":
			f = left
		case "U":
			f = up
		case "D":
			f = down
		}
		for i := 0; i < n; i++ {
			head = f(head)
			tail = fixTail(head, tail)
			locs[tail] = true
		}
	}
	fmt.Printf("Part 1: %d\n", len(locs))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func up(head loc) loc {
	return loc{head.x, head.y + 1}
}

func down(head loc) loc {
	return loc{head.x, head.y - 1}
}

func left(head loc) loc {
	return loc{head.x - 1, head.y}
}

func right(head loc) loc {
	return loc{head.x + 1, head.y}
}

func fixTail(head, tail loc) loc {
	if max(abs(head.x-tail.x), abs(head.y-tail.y)) <= 1 {
		return tail
	}
	if head.x == tail.x {
		if tail.y < head.y {
			return loc{tail.x, tail.y + 1}
		} else {
			return loc{tail.x, tail.y - 1}
		}
	}
	if head.y == tail.y {
		if tail.x < head.x {
			return loc{tail.x + 1, tail.y}
		} else {
			return loc{tail.x - 1, tail.y}
		}
	}
	if head.y > tail.y && head.x > tail.x {
		return loc{tail.x + 1, tail.y + 1}
	}
	if head.y < tail.y && head.x < tail.x {
		return loc{tail.x - 1, tail.y - 1}
	}
	if head.y > tail.y && head.x < tail.x {
		return loc{tail.x - 1, tail.y + 1}
	}
	if head.y < tail.y && head.x > tail.x {
		return loc{tail.x + 1, tail.y - 1}
	}
	return tail
}

func part2() {

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	rope := make([]loc, 10)
	locs := map[loc]bool{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Split(line, " ")
		dir := parts[0]
		n, _ := strconv.Atoi(parts[1])
		var f func(loc) loc
		switch dir {
		case "R":
			f = right
		case "L":
			f = left
		case "U":
			f = up
		case "D":
			f = down
		}
		for i := 0; i < n; i++ {
			rope[0] = f(rope[0])
			for j := 1; j < 10; j++ {
				rope[j] = fixTail(rope[j-1], rope[j])
			}
			locs[rope[9]] = true
		}
	}
	fmt.Printf("Part 2: %d\n", len(locs))
}
