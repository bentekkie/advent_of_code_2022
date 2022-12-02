package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
	total := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		total += scorePart1(line)
	}
	fmt.Printf("Part 1: Total Score=%d\n", total)
}

// A = Rock
// B = Paper
// C = Scissors

// X = Rock
// Y = Paper
// Z = Scissors

func scorePart1(line string) int {
	switch line {
	case "A X":
		return 1 + 3
	case "A Y":
		return 2 + 6
	case "A Z":
		return 3 + 0
	case "B X":
		return 1 + 0
	case "B Y":
		return 2 + 3
	case "B Z":
		return 3 + 6
	case "C X":
		return 1 + 6
	case "C Y":
		return 2 + 0
	case "C Z":
		return 3 + 3
	}
	return 0
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
		total += scorePart2(line)
	}
	fmt.Printf("Part 2: Total Score=%d\n", total)
}

// A = Rock
// B = Paper
// C = Scissors

// X = Lose
// Y = Draw
// Z = Win

// 1 = Rock
// 2 = Paper
// 3 = Scissors

func scorePart2(line string) int {
	switch line {
	case "A X":
		return 3 + 0
	case "A Y":
		return 1 + 3
	case "A Z":
		return 2 + 6
	case "B X":
		return 1 + 0
	case "B Y":
		return 2 + 3
	case "B Z":
		return 3 + 6
	case "C X":
		return 2 + 0
	case "C Y":
		return 3 + 3
	case "C Z":
		return 1 + 6
	}
	return 0
}
