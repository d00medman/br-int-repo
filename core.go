package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type pathingConfig struct {
	SourceNode   int `yaml:"SourceNode,omitempty"`
	DestNode     int `yaml:"DestNode,omitempty"`
	pathLength   int
	shortestPath []int
	Edges        map[string]int `yaml:"Edges"`
	mu           sync.Mutex
	start        time.Time
}

func (cfg *pathingConfig) findShortestPath(source, routeScore int, routeNodes []int, wg *sync.WaitGroup) {
	defer wg.Done()
	cfg.mu.Lock()
	if routeScore > cfg.pathLength {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	cfg.mu.Lock()
	if routeScore < cfg.pathLength && source == cfg.DestNode {
		cfg.shortestPath = append(routeNodes, source)
		cfg.pathLength = routeScore
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()
	nextSteps := cfg.findNextSteps(source)
	for _, k := range nextSteps {
		nextNodes := append(routeNodes, source)
		nextSource, err := strconv.Atoi(strings.Split(k, "-")[1])
		if err != nil {
			log.Fatal("Error in strconv in findShortestPath")
		}
		cfg.mu.Lock()
		nextScore := routeScore + cfg.Edges[k]
		cfg.mu.Unlock()
		wg.Add(1)
		go cfg.findShortestPath(nextSource, nextScore, nextNodes, wg)
	}
}

func (cfg *pathingConfig) findNextSteps(source int) []string {
	nextSteps := []string{}
	cfg.mu.Lock()
	for k, _ := range cfg.Edges {
		splitPrefix, err := strconv.Atoi(strings.Split(k, "-")[0])
		if err != nil {
			log.Fatal("Error in strconv in findNextSteps")
		}
		if splitPrefix == source {
			nextSteps = append(nextSteps, k)
		}
	}
	cfg.mu.Unlock()
	return nextSteps
}

func (cfg *pathingConfig) findPointDistance(a, b int) int {
	k := getSyncMapKey(a, b)
	cfg.mu.Lock()
	dist := cfg.Edges[k]
	cfg.mu.Unlock()
	return dist
}

func (cfg *pathingConfig) displayResult() {
	elapsed := time.Since(cfg.start)
	builder := strings.Builder{}
	var prev int
	for i, n := range cfg.shortestPath {
		prev = n
		if i == 0 {
			builder.WriteString("Source")
		} else if i == len(cfg.shortestPath)-1 {
			builder.WriteString("Dest")
		}
		builder.WriteString("Node")
		builder.WriteString(strconv.Itoa(n))
		if i != len(cfg.shortestPath)-1 {
			builder.WriteString("--(")
			b := cfg.shortestPath[i+1]
			builder.WriteString(strconv.Itoa(cfg.findPointDistance(prev, b)))
			builder.WriteString(")-->")
		} else {
			builder.WriteString("\nShortest distance from nodes ")
			builder.WriteString(strconv.Itoa(cfg.SourceNode))
			builder.WriteString(" to ")
			builder.WriteString(strconv.Itoa(cfg.DestNode))
			builder.WriteString("\nNodes Visited: ")
			builder.WriteString(strconv.Itoa(len(cfg.shortestPath)))
			builder.WriteString("\nDistance traveled: ")
			builder.WriteString(strconv.Itoa(cfg.pathLength))
			builder.WriteString("\nExecution time: ")
			builder.WriteString(elapsed.String())
		}
	}
	fmt.Println(builder.String())
}
func getSyncMapKey(a, b int) string {
	builder := strings.Builder{}
	builder.WriteString(strconv.Itoa(a))
	builder.WriteString("-")
	builder.WriteString(strconv.Itoa(b))
	return builder.String()
}
