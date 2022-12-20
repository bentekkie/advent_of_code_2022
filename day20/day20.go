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

type Num[T any] struct {
	prev, next *Num[T]
	n, actual  T
}

func (n *Num[T]) forward(k int64) {
	for i := int64(0); i < k; i++ {
		cnext := n.next
		cprev := n.prev
		n.next = cnext.next
		n.next.prev = n
		n.prev = cnext
		n.prev.next = n
		cprev.next = cnext
		cnext.prev = cprev
	}
}

func (n *Num[T]) backward(k int64) {
	for i := int64(0); i < k; i++ {
		cnext := n.next
		cprev := n.prev
		n.prev = cprev.prev
		n.prev.next = n
		n.next = cprev
		n.next.prev = n
		cprev.next = cnext
		cnext.prev = cprev
	}
}

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var zero *Num[int64]
	nums := []*Num[int64]{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			n, _ := strconv.Atoi(line)
			num := &Num[int64]{n: int64(n), actual: int64(n)}
			nums = append(nums, num)
			if n == 0 {
				zero = num
			}
		}
	}
	nums[0].prev = nums[len(nums)-1]
	nums[len(nums)-1].next = nums[0]
	for i, num := range nums {
		if i > 0 {
			num.prev = nums[i-1]
		}
		if i < len(nums)-1 {
			num.next = nums[i+1]
		}
	}
	for _, num := range nums {
		if num.n > 0 {
			num.forward(num.n)
		}
		if num.n < 0 {
			num.backward(-num.n)
		}
	}
	ns := toArr(zero)
	total := ns[1000%len(ns)] + ns[2000%len(ns)] + ns[3000%len(ns)]
	fmt.Printf("Part 1: %d\n", total)
}

func printList[T any](start *Num[T]) {
	fmt.Printf("%v, ", start.actual)
	curr := start.next
	for curr != start {
		fmt.Printf("%v, ", curr.actual)
		curr = curr.next
	}
	fmt.Print("\n")
}

func toArr[T any](zero *Num[T]) []T {
	out := []T{zero.actual}
	curr := zero.next
	for curr != zero {
		out = append(out, curr.actual)
		curr = curr.next
	}
	return out
}

func part2() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var zero *Num[int64]
	nums := []*Num[int64]{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			n, _ := strconv.Atoi(line)
			num := &Num[int64]{actual: int64(n) * 811589153}
			nums = append(nums, num)
			if n == 0 {
				zero = num
			}
		}
	}
	nums[0].prev = nums[len(nums)-1]
	nums[len(nums)-1].next = nums[0]
	for i, num := range nums {
		num.n = num.actual % int64(len(nums)-1)
		if i > 0 {
			num.prev = nums[i-1]
		}
		if i < len(nums)-1 {
			num.next = nums[i+1]
		}
	}
	for round := 0; round < 10; round++ {
		for _, num := range nums {
			if num.n > 0 {
				num.forward(num.n)
			}
			if num.n < 0 {
				num.backward(-num.n)
			}
		}
	}
	ns := toArr(zero)
	total := ns[1000%len(ns)] + ns[2000%len(ns)] + ns[3000%len(ns)]
	fmt.Printf("Part 2: %d\n", total)
}
