package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	b, _ := os.ReadFile("input.txt")
	input := strings.TrimSpace(string(b))
	part1(input)
	part2(input)
}

func part1(input string) {
	fmt.Printf("Part 1: %d\n", marker(input, 4))
}
func part2(input string) {
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
