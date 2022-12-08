package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	root, sizes := buildTree()
	part1(root, sizes)
	part2(root, sizes)
}

func buildTree() (*dir, map[*dir]int) {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	cmds := []*command{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "$") {
			cmds = append(cmds, &command{
				input: strings.TrimPrefix(line, "$ "),
			})
		} else {
			cmds[len(cmds)-1].output = append(cmds[len(cmds)-1].output, line)
		}
	}
	root := &dir{
		files: map[string]int{},
		dirs: map[string]*dir{
			"/": {
				files: map[string]int{},
				dirs:  map[string]*dir{},
			},
		},
	}
	currentDir := root
	for _, cmd := range cmds {
		if cmd.input == "ls" {
			for _, line := range cmd.output {
				parts := strings.Split(line, " ")
				if parts[0] == "dir" {
					currentDir.dirs[parts[1]] = &dir{
						parent: currentDir,
						files:  map[string]int{},
						dirs:   map[string]*dir{},
					}
				} else {
					size, _ := strconv.Atoi(parts[0])
					currentDir.files[parts[1]] = size

				}
			}
		} else if strings.HasPrefix(cmd.input, "cd") {
			newDirName := strings.TrimPrefix(cmd.input, "cd ")
			if newDirName == ".." {
				currentDir = currentDir.parent
			} else {
				currentDir = currentDir.dirs[newDirName]
			}
		}
	}
	sizes := make(map[*dir]int)
	root.dirs["/"].calculateSizes(sizes)
	return root, sizes
}

func part1(root *dir, sizes map[*dir]int) {
	total := 0
	for _, size := range sizes {
		if size <= 100000 {
			total += size
		}
	}
	fmt.Printf("Part 1: %d\n", total)
}

type command struct {
	input  string
	output []string
}

type dir struct {
	parent *dir
	files  map[string]int
	dirs   map[string]*dir
}

func (d *dir) calculateSizes(sizes map[*dir]int) {
	total := 0
	for _, size := range d.files {
		total += size
	}
	for _, dir := range d.dirs {
		dir.calculateSizes(sizes)
		total += sizes[dir]
	}
	sizes[d] = total
}

func (d *dir) Tree() string {
	return strings.Join(d.stringlines(), "\n")
}

func (d *dir) stringlines() []string {
	var lines []string
	for name, dir := range d.dirs {
		lines = append(lines, fmt.Sprintf("- %s (dir)", name))
		for _, line := range dir.stringlines() {
			lines = append(lines, "  "+line)
		}
	}
	for name, size := range d.files {
		lines = append(lines, fmt.Sprintf("- %s (file, size=%d)", name, size))
	}
	return lines
}

func part2(root *dir, sizes map[*dir]int) {
	rootSize := sizes[root.dirs["/"]]
	freeSpace := 70000000 - rootSize
	allSizes := []int{}
	for _, size := range sizes {
		allSizes = append(allSizes, size)
	}
	sort.Slice(allSizes, func(i, j int) bool { return allSizes[i] < allSizes[j] })
	k := 0
	for freeSpace+allSizes[k] < 30000000 {
		k += 1
	}
	fmt.Printf("Part 2: %d\n", allSizes[k])
}
