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

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var maxCal, currCal int

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			cal, _ := strconv.Atoi(line)
			currCal += cal
		} else {
			if currCal > maxCal {
				maxCal = currCal
			}
			currCal = 0
		}
	}
	if currCal > maxCal {
		maxCal = currCal
	}
	fmt.Printf("Part1: Max calories=%d\n", maxCal)
}

func part2() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var first, second, third, currCal int

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			cal, _ := strconv.Atoi(line)
			currCal += cal
		} else {
			if currCal > first {
				third = second
				second = first
				first = currCal
			} else if currCal > second {
				third = second
				second = currCal
			} else if currCal > third {
				third = currCal
			}
			currCal = 0
		}
	}
	if currCal > first {
		third = second
		second = first
		first = currCal
	} else if currCal > second {
		third = second
		second = currCal
	} else if currCal > third {
		third = currCal
	}
	fmt.Printf("Part2: Max calories=%d\n", first+second+third)
}
