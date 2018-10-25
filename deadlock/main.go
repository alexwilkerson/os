package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	processes map[int]*process
	resources map[int]*resource
)

type resource struct {
	id       int
	checked  bool
	owned    *process
	neededBy []*process
}

func (r *resource) Next() {
	if len(r.neededBy) > 0 {
		r.owned = r.neededBy[0]
		r.owned.usedResources[r.id] = r
		r.neededBy = r.neededBy[1:]
		r.owned.neededResource = nil
		c := color.New(color.FgYellow)
		c.Printf("Resource %d is allocated to process %d.\n", r.id, r.owned.id)
		checkDeadlock(r.owned, []int{}, []int{})
	} else {
		r.owned = nil
		c := color.New(color.FgBlue)
		c.Printf("Resource %d is now free.\n", r.id)
	}
}

type process struct {
	id             int
	checked        bool
	neededResource *resource
	usedResources  map[int]*resource
}

func (p *process) Needs(res int) {
	if r, ok := resources[res]; ok {
		if r.owned != nil {
			r.neededBy = append(r.neededBy, p)
			p.neededResource = r
			c := color.New(color.FgMagenta)
			c.Printf("Process %d must wait.\n", p.id)
		} else {
			r.owned = p
			c := color.New(color.FgYellow)
			c.Printf("Resource %d is allocated to process %d.\n", r.id, p.id)
			p.usedResources[res] = r
		}
	} else {
		resources[res] = &resource{id: res, owned: p}
		c := color.New(color.FgYellow)
		c.Printf("Resource %d is allocated to process %d.\n", resources[res].id, p.id)
		p.usedResources[res] = resources[res]
	}
	checkDeadlock(p, []int{}, []int{})
}

func (p *process) Releases(res int) {
	r := p.usedResources[res]
	delete(p.usedResources, res)
	r.Next()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Wrong number of arguments.\nUsage: deadlock [input file]")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Could not open file \"%s\".\n", os.Args[1])
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	processes = make(map[int]*process)
	resources = make(map[int]*resource)

	for scanner.Scan() {
		parseLine(scanner.Text())
	}

	c := color.New(color.FgGreen)

	c.Println("EXECUTION COMPLETED: No deadlock encountered.")
}

func parseLine(line string) {
	fields := strings.Fields(line)
	id, _ := strconv.Atoi(fields[0])
	res, _ := strconv.Atoi(fields[2])
	var proc *process
	if val, ok := processes[id]; ok {
		proc = val
	} else {
		proc = &process{id: id}
		proc.usedResources = make(map[int]*resource)
		processes[id] = proc
	}
	c := color.New(color.FgCyan)
	if fields[1] == "N" {
		c.Printf("Process %d needs resource %d – ", id, res)
		proc.Needs(res)
	} else if fields[1] == "R" {
		c.Printf("Process %d releases resource %d – ", id, res)
		proc.Releases(res)
	} else {
		fmt.Println("Input file formatting error.")
		os.Exit(1)
	}
}

func checkDeadlock(proc *process, pids, rids []int) {
	if proc.checked && contains(pids, proc.id) {
		printDeadlock(pids, rids)
		os.Exit(0)
	} else if proc.checked {
		return
	}
	if res := proc.neededResource; res != nil {
		if nextProc := res.owned; nextProc != nil {
			proc.checked = true
			pids = append(pids, proc.id)
			rids = append(rids, proc.neededResource.id)
			checkDeadlock(nextProc, pids, rids)
		}
	}
	removeChecked()
}

func removeChecked() {
	for _, p := range processes {
		p.checked = false
	}

	for _, r := range resources {
		r.checked = false
	}
}

func contains(l []int, v int) bool {
	for _, n := range l {
		if n == v {
			return true
		}
	}
	return false
}

func printDeadlock(pids, rids []int) {
	sort.Ints(pids)
	sort.Ints(rids)

	c := color.New(color.FgRed)

	c.Printf("DEADLOCK DETECTED: Processes ")
	for _, p := range pids {
		c.Printf("%v, ", p)
	}
	c.Printf("and Resources ")
	for _, r := range rids {
		c.Printf("%v, ", r)
	}
	c.Printf("are found in a cycle.\n")
}
