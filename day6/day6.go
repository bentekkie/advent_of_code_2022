package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	part1()
	part2()
}

func part1() {
	b, _ := os.ReadFile("input.txt")
	input := strings.TrimSpace(string(b))
	fmt.Printf("Part 1: %d\n", marker(input, 4))
}
func part2() {
	b, _ := os.ReadFile("input.txt")
	input := strings.TrimSpace(string(b))
	fmt.Printf("Part 2: %d\n", marker(input, 14))
}

type window[T comparable] struct {
	keyCnt int
	m      map[T]int
}

func (w *window[T]) Put(r T) {
	if w.m[r] == 0 {
		w.keyCnt += 1
	}
	w.m[r] += 1
}
func (w *window[T]) Remove(r T) {
	if w.m[r] == 1 {
		w.keyCnt -= 1
	}
	w.m[r] -= 1
}

func marker(input string, n int) int {
	runes := []rune(input)
	w := &window[rune]{
		m: make(map[rune]int, n),
	}
	for i, r := range runes {
		if i >= n {
			w.Remove(runes[i-n])
		}
		w.Put(r)
		if i > n-2 {
			if w.keyCnt == n {
				return i + 1
			}
		}
	}
	return -1
}
