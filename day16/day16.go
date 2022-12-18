package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
	"gonum.org/v1/gonum/graph/path"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

func main() {
	part1()
	part2()
}

type Node struct {
	name string
	flow int
}

var lineRe = regexp.MustCompile(`Valve (?P<node>.*) has flow rate=(?P<flow>\d*); tunnel[s]? lead[s]? to valve[s]? (?P<tunnels>.*)`)

func reSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}
	return subMatchMap
}

type GInfo struct {
	flows      map[string]int
	adjList    map[string][]string
	nameToNode map[string]graph.Node
	nodeToName map[graph.Node]string
	paths      path.AllShortest
}

func part1() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	ginfo := &GInfo{
		flows:      map[string]int{},
		adjList:    map[string][]string{},
		nameToNode: map[string]graph.Node{},
		nodeToName: map[graph.Node]string{},
	}
	g := simple.NewUndirectedGraph()
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		matches := reSubMatchMap(lineRe, line)
		flow, _ := strconv.Atoi(matches["flow"])
		name := matches["node"]
		tunnels := strings.Split(matches["tunnels"], ", ")
		ginfo.flows[name] = flow
		ginfo.adjList[name] = tunnels
		node := g.NewNode()
		ginfo.nameToNode[name] = node
		ginfo.nodeToName[node] = name
		g.AddNode(node)
	}
	isOpen := map[string]bool{}
	for name, adjs := range ginfo.adjList {
		isOpen[name] = ginfo.flows[name] == 0
		for _, aname := range adjs {
			g.SetEdge(g.NewEdge(ginfo.nameToNode[name], ginfo.nameToNode[aname]))
		}
	}
	ginfo.paths, _ = path.FloydWarshall(g)
	fmt.Printf("Part 1: %d\n", traverse("AA", 0, isOpen, ginfo, 30))
}

type Weight struct {
	to   string
	w, f int
}

func traverse(curr string, time int, isOpen map[string]bool, ginfo *GInfo, totalTime int) int {
	oldOpen := isOpen[curr]
	isOpen[curr] = true
	defer func() {
		isOpen[curr] = oldOpen
	}()
	currNode := ginfo.nameToNode[curr]
	neighbours := make([]Weight, 0, len(isOpen))
	for name, open := range isOpen {
		if name != curr && !open {
			w := int(ginfo.paths.Weight(currNode.ID(), ginfo.nameToNode[name].ID()))
			if w+time+1 < totalTime {
				neighbours = append(neighbours, Weight{
					to: name,
					w:  w,
					f:  ginfo.flows[name],
				})
			}
		}
	}
	slices.SortFunc(neighbours, func(a, b Weight) bool {
		return a.w < b.w || (a.w == b.w && a.f > b.f)
	})
	max := 0
	for _, n := range neighbours {
		timeValveStarts := time + n.w + 1
		value := (totalTime-timeValveStarts)*ginfo.flows[n.to] + traverse(n.to, timeValveStarts, isOpen, ginfo, totalTime)
		if value > max {
			max = value
		}
	}
	return max
}

func part2() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	ginfo := &GInfo{
		flows:      map[string]int{},
		adjList:    map[string][]string{},
		nameToNode: map[string]graph.Node{},
		nodeToName: map[graph.Node]string{},
	}
	g := simple.NewUndirectedGraph()
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		matches := reSubMatchMap(lineRe, line)
		flow, _ := strconv.Atoi(matches["flow"])
		name := matches["node"]
		tunnels := strings.Split(matches["tunnels"], ", ")
		ginfo.flows[name] = flow
		ginfo.adjList[name] = tunnels
		node := g.NewNode()
		ginfo.nameToNode[name] = node
		ginfo.nodeToName[node] = name
		g.AddNode(node)
	}
	isOpen := map[string]bool{}
	keys := []string{}
	var open int64
	for name, adjs := range ginfo.adjList {
		isOpen[name] = ginfo.flows[name] == 0
		keys = append(keys, name)
		for _, aname := range adjs {
			g.SetEdge(g.NewEdge(ginfo.nameToNode[name], ginfo.nameToNode[aname]))
		}
	}
	sort.Strings(keys)
	indices := map[string]uint{}
	for i, name := range keys {
		indices[name] = uint(i)
		if isOpen[name] {
			open = setBit(open, indices[name])
		}
	}
	ginfo.paths, _ = path.FloydWarshall(g)
	fmt.Printf("Part 2: %d\n", traverseWithElephant("AA", "AA", 0, 0, 0, open, ginfo, 26, keys, indices, map[cachekey]int{}))
}

// Sets the bit at pos in the integer n.
func setBit(n int64, pos uint) int64 {
	n |= (1 << pos)
	return n
}

// Clears the bit at pos in n.
func clearBit(n int64, pos int64) int64 {
	mask := ^(int64(1) << pos)
	n &= mask
	return n
}
func hasBit(n int64, pos uint) bool {
	val := n & (int64(1) << pos)
	return (val > 0)
}

type Bitmask uint64

func (f Bitmask) Has(pos uint) bool { return f&(Bitmask(1)<<pos) != 0 }
func (f *Bitmask) Set(pos uint)     { *f |= (1 << pos) }

type cachekey struct {
	you, elephant                 string
	open                          int64
	time, youtimefree, eltimefree int
}

func tkey(you, elephant string, time, youtimefree, eltimefree int, open int64, keys []string) cachekey {
	k := cachekey{
		you:         you,
		elephant:    elephant,
		time:        time,
		youtimefree: youtimefree,
		eltimefree:  eltimefree,
		open:        open,
	}
	if k.eltimefree < k.youtimefree {
		k.eltimefree, k.youtimefree = k.youtimefree, k.eltimefree
		k.elephant, k.you = k.you, k.elephant
	}
	return k
}

func traverseWithElephant(you, elephant string, time, youtimefree, eltimefree int, open int64, ginfo *GInfo, totalTime int, keys []string, indices map[string]uint, mem map[cachekey]int) int {
	key := tkey(you, elephant, time, youtimefree, eltimefree, open, keys)
	if v, ok := mem[key]; ok {
		return v
	}
	open = setBit(setBit(open, indices[you]), indices[elephant])
	currNodeYou := ginfo.nameToNode[you]
	currNodeElephant := ginfo.nameToNode[elephant]
	neighboursYou := make([]Weight, 0, len(keys))
	neighboursElephant := make([]Weight, 0, len(keys))
	for _, name := range keys {
		if !hasBit(open, indices[name]) {
			if time == youtimefree {
				wyou := int(ginfo.paths.Weight(currNodeYou.ID(), ginfo.nameToNode[name].ID()))
				if wyou+time+1 < totalTime {
					neighboursYou = append(neighboursYou, Weight{
						to: name,
						w:  wyou,
						f:  ginfo.flows[name],
					})
				}
			}
			if time == eltimefree {
				welephant := int(ginfo.paths.Weight(currNodeElephant.ID(), ginfo.nameToNode[name].ID()))
				if welephant+time+1 < totalTime {
					neighboursElephant = append(neighboursElephant, Weight{
						to: name,
						w:  welephant,
						f:  ginfo.flows[name],
					})
				}
			}
		}
	}
	if len(neighboursYou) == 0 && len(neighboursElephant) == 0 {
		mem[key] = 0
		return 0
	}
	if len(neighboursYou) == 0 {
		max := 0
		for _, n := range neighboursElephant {
			eltimefree = time + n.w + 1
			value := (totalTime-eltimefree)*n.f + traverseWithElephant(you, n.to, min(youtimefree, eltimefree), youtimefree, eltimefree, open, ginfo, totalTime, keys, indices, mem)
			if value > max {
				max = value
			}
		}
		mem[key] = max
		return max
	}
	if len(neighboursElephant) == 0 {
		max := 0
		for _, n := range neighboursYou {
			youtimefree = time + n.w + 1
			value := (totalTime-youtimefree)*n.f + traverseWithElephant(n.to, elephant, min(youtimefree, eltimefree), youtimefree, eltimefree, open, ginfo, totalTime, keys, indices, mem)
			if value > max {
				max = value
			}
		}
		mem[key] = max
		return max
	}
	max := 0
	for _, nyou := range neighboursYou {
		for _, nel := range neighboursElephant {
			if nyou.to != nel.to {
				youtimefree = time + nyou.w + 1
				eltimefree = time + nel.w + 1
				value := (totalTime-eltimefree)*nel.f + (totalTime-youtimefree)*nyou.f + traverseWithElephant(nyou.to, nel.to, min(youtimefree, eltimefree), youtimefree, eltimefree, open, ginfo, totalTime, keys, indices, mem)
				if value > max {
					max = value
				}
			}
		}
	}
	mem[key] = max
	return max
}

func min[T constraints.Ordered](a T, b T) T {
	if b < a {
		return b
	}
	return a
}
