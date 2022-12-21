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

type Monkey interface {
	Number() float64
}

type NumMonkey struct {
	n float64
}

func (m *NumMonkey) Number() float64 {
	return m.n
}

type MathMonkey struct {
	lefts, rights string
	left, right   Monkey
	op            rune
	calced        bool
	n             float64
}

func (m *MathMonkey) Number() float64 {
	if m.calced {
		return m.n
	}
	leftn, rightn := m.left.Number(), m.right.Number()
	switch m.op {
	case '+':
		m.n = leftn + rightn
	case '-':
		m.n = leftn - rightn
	case '*':
		m.n = leftn * rightn
	case '/':
		m.n = leftn / rightn
	}
	m.calced = true
	return m.n
}

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	mathms := []*MathMonkey{}
	ms := map[string]Monkey{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, ": ")
		name := parts[0]
		rparts := strings.Split(parts[1], " ")
		if len(rparts) == 1 {
			n, _ := strconv.ParseFloat(rparts[0], 64)
			ms[name] = &NumMonkey{n}
		} else {
			m := &MathMonkey{
				lefts:  rparts[0],
				rights: rparts[2],
				op:     []rune(rparts[1])[0],
			}
			mathms = append(mathms, m)
			ms[name] = m
		}
	}
	for _, m := range mathms {
		m.left = ms[m.lefts]
		m.right = ms[m.rights]
	}
	fmt.Printf("Part 1: %d\n", int(ms["root"].Number()))
}

func part2() {

	fmt.Printf("Part 2: Use excel solver :P\n")
}
