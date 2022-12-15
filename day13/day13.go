package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

func main() {
	part1()
	part2()
}

type EntryType int

const (
	numType  EntryType = iota
	listType EntryType = iota
)

type CompareResult int

const (
	correctOrder   CompareResult = iota
	sameOrder      CompareResult = iota
	incorrectOrder CompareResult = iota
)

func (res CompareResult) String() string {
	switch res {
	case correctOrder:
		return "correct"
	case incorrectOrder:
		return "incorrect"
	case sameOrder:
		return "same"
	}
	return "invalid"
}

type Entry struct {
	isDivider bool
	t         EntryType
	i         int64
	list      []*Entry
}

func (e *Entry) String() string {
	if e.t == numType {
		return fmt.Sprintf("%d", e.i)
	}
	return fmt.Sprintf("%s", e.list)
}

func compare(left, right *Entry) CompareResult {
	if left.t == numType && right.t == numType {
		if left.i < right.i {
			return correctOrder
		}
		if left.i > right.i {
			return incorrectOrder
		}
		return sameOrder
	}
	if left.t == numType {
		return compare(&Entry{
			t:    listType,
			list: []*Entry{left},
		}, right)
	}
	if right.t == numType {
		return compare(left, &Entry{
			t:    listType,
			list: []*Entry{right},
		})
	}
	for i, le := range left.list {
		if i == len(right.list) {
			return incorrectOrder
		}
		res := compare(le, right.list[i])
		if res != sameOrder {
			return res
		}
	}
	if len(left.list) == len(right.list) {
		return sameOrder
	}
	return correctOrder
}

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	//total := 0
	linePairs := [][]string{}
	currPair := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			linePairs = append(linePairs, currPair)
			currPair = []string{}
		} else {
			currPair = append(currPair, line)
		}
	}
	linePairs = append(linePairs, currPair)
	total := 0
	for i, linePair := range linePairs {
		res := compare(
			newEntry(decodeJson(linePair[0])),
			newEntry(decodeJson(linePair[1])),
		)
		if res == correctOrder {
			total += i + 1
		}
	}
	fmt.Printf("Part 1: %d\n", total)
}

func decodeJson(str string) []interface{} {
	d := json.NewDecoder(strings.NewReader(str))
	d.UseNumber()
	raw := []interface{}{}
	d.Decode(&raw)
	return raw
}

func newEntry(raw []interface{}) *Entry {
	curr := &Entry{
		t: listType,
	}
	for _, val := range raw {
		switch val := val.(type) {
		case json.Number:
			i, _ := val.Int64()
			curr.list = append(curr.list, &Entry{
				t: numType,
				i: i,
			})
		case []interface{}:
			curr.list = append(curr.list, newEntry(val))
		}
	}
	return curr
}

func part2() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	entries := []*Entry{
		newEntry(decodeJson("[[2]]")),
		newEntry(decodeJson("[[6]]")),
	}
	entries[0].isDivider = true
	entries[1].isDivider = true
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			entries = append(entries, newEntry(decodeJson(line)))
		}
	}
	slices.SortFunc(entries, func(a, b *Entry) bool { return compare(a, b) != incorrectOrder })
	total := 1
	for i, e := range entries {
		if e.isDivider {
			total *= i + 1
		}
	}
	fmt.Printf("Part 2: %d\n", total)
}
