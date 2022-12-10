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

type cycle struct {
	during, end int
}

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	xs := []cycle{{during: 1, end: 1}}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "noop" {
			last := xs[len(xs)-1]
			xs = append(
				xs,
				cycle{
					during: last.end,
					end:    last.end,
				},
			)
		} else {
			parts := strings.Split(line, " ")
			delta, _ := strconv.Atoi(parts[1])
			last := xs[len(xs)-1]
			xs = append(
				xs,
				cycle{
					during: last.end,
					end:    last.end,
				},
				cycle{
					during: last.end,
					end:    last.end + delta,
				},
			)
		}
	}
	total := 0
	for i := 20; i < len(xs); i += 40 {
		total += i * xs[i].during
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
	xs := []cycle{{during: 1, end: 1}}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "noop" {
			last := xs[len(xs)-1]
			xs = append(
				xs,
				cycle{
					during: last.end,
					end:    last.end,
				},
			)
		} else {
			parts := strings.Split(line, " ")
			delta, _ := strconv.Atoi(parts[1])
			last := xs[len(xs)-1]
			xs = append(
				xs,
				cycle{
					during: last.end,
					end:    last.end,
				},
				cycle{
					during: last.end,
					end:    last.end + delta,
				},
			)
		}
	}
	var sb strings.Builder
	for row := 0; row < 6; row++ {
		for col := 0; col < 40; col++ {
			curr := xs[1+row*40+col].during
			if curr-1 <= col && col <= curr+1 {
				sb.WriteRune('#')
			} else {
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune('\n')
	}
	fmt.Printf("Part 2: \n%s", sb.String())
}
