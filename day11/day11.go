package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	flag.Parse()
	ms := parseMonkeys()
	part1(ms)
	ms = parseMonkeys()
	part2(ms)
}

type opPart interface {
	val(old int) int
}

type useOld struct{}

func (*useOld) val(old int) int {
	return old
}

type useVal struct {
	v int
}

func (u *useVal) val(old int) int {
	return u.v
}

type thing struct {
	origVal int
	worries map[int]int
}

type monkey struct {
	divBy           int
	ifTrue, ifFalse int
	opLeft, opRight opPart
	op              rune
	items           []*thing
	timesInspected  int
}

func parseMonkeys() map[int]*monkey {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	ms := map[int]*monkey{}
	var currM int
	for _, line := range lines {
		if strings.HasPrefix(line, "Monkey ") {
			i, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(line, "Monkey "), ":"))
			currM = i
			ms[currM] = &monkey{}
		} else if strings.HasPrefix(line, "Starting items: ") {
			items := strings.Split(strings.TrimPrefix(line, "Starting items: "), ", ")
			for _, it := range items {
				val, _ := strconv.Atoi(it)
				ms[currM].items = append(ms[currM].items, &thing{
					origVal: val,
					worries: map[int]int{},
				})
			}
		} else if strings.HasPrefix(line, "Operation: new = old ") {
			op := strings.TrimPrefix(line, "Operation: new = old ")
			ms[currM].opLeft = &useOld{}
			rs := []rune(op)
			ms[currM].op = rs[0]
			rop := string(rs[2:])
			if rop == "old" {
				ms[currM].opRight = &useOld{}
			} else {
				val, _ := strconv.Atoi(rop)
				ms[currM].opRight = &useVal{val}
			}
		} else if strings.HasPrefix(line, "Test: divisible by ") {
			i, _ := strconv.Atoi(strings.TrimPrefix(line, "Test: divisible by "))
			ms[currM].divBy = i
		} else if strings.HasPrefix(line, "If true: throw to monkey ") {
			ms[currM].ifTrue, _ = strconv.Atoi(strings.TrimPrefix(line, "If true: throw to monkey "))
		} else if strings.HasPrefix(line, "If false: throw to monkey ") {
			ms[currM].ifFalse, _ = strconv.Atoi(strings.TrimPrefix(line, "If false: throw to monkey "))
		}
	}
	for mi := 0; mi < len(ms); mi++ {
		for mj := 0; mj < len(ms); mj++ {
			for _, item := range ms[mj].items {
				item.worries[mi] = item.origVal % ms[mi].divBy
			}
		}
	}
	return ms
}
func part1(ms map[int]*monkey) {
	for round := 0; round < 20; round++ {
		for mi := 0; mi < len(ms); mi++ {
			m := ms[mi]
			for _, it := range m.items {
				it.origVal = applyOp(it.origVal, m.opLeft, m.opRight, m.op)
				it.origVal /= 3
				if it.origVal%m.divBy == 0 {
					ms[m.ifTrue].items = append(ms[m.ifTrue].items, it)
				} else {
					ms[m.ifFalse].items = append(ms[m.ifFalse].items, it)
				}
				m.timesInspected++
			}
			m.items = []*thing{}
		}
	}
	tis := []int{}
	for mi := 0; mi < len(ms); mi++ {
		tis = append(tis, ms[mi].timesInspected)
	}
	sort.Ints(tis)
	fmt.Printf("Part 1: %v\n", tis[len(tis)-1]*tis[len(tis)-2])
}

func applyOp(old int, left, right opPart, op rune) int {
	lval := left.val(old)
	rval := right.val(old)
	switch op {
	case '+':
		return lval + rval
	case '*':
		return lval * rval
	}
	panic("Invalid op")
}
func applyOpMod(old int, left, right opPart, op rune, mod int) int {
	lval := left.val(old)
	rval := right.val(old)
	switch op {
	case '+':
		return (lval + rval) % mod
	case '*':
		return (lval * rval) % mod
	}
	panic("Invalid op")
}

func part2(ms map[int]*monkey) {
	for round := 0; round < 10000; round++ {
		for mi := 0; mi < len(ms); mi++ {
			m := ms[mi]
			for _, it := range m.items {
				for mj := 0; mj < len(ms); mj++ {
					it.worries[mj] = applyOpMod(it.worries[mj], m.opLeft, m.opRight, m.op, ms[mj].divBy)
				}
				if it.worries[mi] == 0 {
					ms[m.ifTrue].items = append(ms[m.ifTrue].items, it)
				} else {
					ms[m.ifFalse].items = append(ms[m.ifFalse].items, it)
				}
				m.timesInspected++
			}
			m.items = []*thing{}
		}
	}
	tis := []int{}
	for mi := 0; mi < len(ms); mi++ {
		tis = append(tis, ms[mi].timesInspected)
	}
	sort.Ints(tis)
	fmt.Printf("Part 2: %v\n", tis[len(tis)-1]*tis[len(tis)-2])
}
