package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
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
	total := ""
	actual := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		total = addSnafu(total, line)
		fmt.Printf("%s = %d\n", line, snaduToNum(line))
		actual += snaduToNum(line)
	}
	fmt.Printf("Part 1: %s  %d %d\n", total, snaduToNum(total), actual)
}

func addSnafu(a, b string) string {
	ars, brs := []rune(a), []rune(b)
	res := ""
	place := 0
	var carry rune
	for place < max(len(a), len(b)) {
		placeVal := 0
		if carry != 0 {
			placeVal += snafuDigitToInt(carry)
		}
		carry = 0
		if place < len(a) {
			placeVal += snafuDigitToInt(ars[len(ars)-1-place])
		}
		if place < len(b) {
			placeVal += snafuDigitToInt(brs[len(brs)-1-place])
		}

		switch placeVal {
		case 0:
			carry = 0
			res = "0" + res
		case 1:
			carry = 0
			res = "1" + res
		case 2:
			carry = 0
			res = "2" + res
		case 3:
			carry = '1'
			res = "=" + res
		case 4:
			carry = '1'
			res = "-" + res
		case 5:
			carry = '1'
			res = "0" + res
		case -1:
			carry = 0
			res = "-" + res
		case -2:
			carry = 0
			res = "=" + res
		case -3:
			carry = '-'
			res = "2" + res
		case -4:
			carry = '-'
			res = "1" + res
		case -5:
			carry = '-'
			res = "0" + res
		default:
			panic("Invalid result " + " " + string(ars[len(ars)-1-place]) + " + " + string(brs[len(brs)-1-place]) + " " + string(carry) + " = " + strconv.Itoa(placeVal))
		}
		place++
	}
	if carry != 0 {
		res = string(carry) + res
	}
	return res
}

func max[T constraints.Ordered](a, b T) T {
	if a < b {
		return b
	}
	return a
}

func snaduToNum(snafu string) int {
	total := 0
	rs := []rune(snafu)
	placeVal := 1
	for i := 0; i < len(rs); i++ {
		total += snafuDigitToInt(rs[len(rs)-1-i]) * placeVal
		placeVal *= 5
	}
	return total
}

func snafuDigitToInt(r rune) int {
	if r == '2' {
		return 2
	}
	if r == '1' {
		return 1
	}
	if r == '0' {
		return 0
	}
	if r == '-' {
		return -1
	}
	if r == '=' {
		return -2
	}
	panic("Invalid rune " + string(r))
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
	}
	fmt.Printf("Part 2: %d\n", total)
}
