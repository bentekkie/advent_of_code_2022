package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	part1()
	part2real()
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
}

func (m *MathMonkey) Number() float64 {
	leftn, rightn := m.left.Number(), m.right.Number()
	switch m.op {
	case '+':
		return leftn + rightn
	case '-':
		return leftn - rightn
	case '*':
		return leftn * rightn
	case '/':
		return leftn / rightn
	}
	panic("Invalid op")
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

func regulaFalsi(goalFunc func(float64) (float64, error), low, high float64, maxItt int) (float64, error) {
	var next float64
	for i := 0; i < maxItt; i++ {
		highFunc, err := goalFunc(high)
		if err != nil {
			return next, err
		}
		lowFunc, err := goalFunc(low)
		if err != nil {
			return next, err
		}
		next := (low*highFunc - high*lowFunc) / (highFunc - lowFunc)
		nextFunc, err := goalFunc(next)
		if err != nil {
			return next, err
		}
		if math.Abs(low-next) <= .00001 || math.Abs(high-next) <= .00001 {
			return next, nil
		}
		if (nextFunc > 0) == (lowFunc > 0) {
			low = next
		} else {
			high = next
		}
	}
	return next, errors.New("regula falsi failed to converge")
}

func part2real() {
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
	root := ms["root"].(*MathMonkey)
	root.op = '-'
	me := ms["humn"].(*NumMonkey)
	low := me.Number()
	for root.Number() > 0 {
		me.n *= 2
	}
	high := me.Number()
	val, _ := regulaFalsi(func(in float64) (float64, error) {
		me.n = in
		return root.Number(), nil
	}, low, high, 20000)

	fmt.Printf("Part 2: %d\n", int(val))
}
