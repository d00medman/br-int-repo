package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

type pathingConfig struct {
	SourceNode   int `yaml:"SourceNode,omitempty"`
	DestNode     int `yaml:"DestNode,omitempty"`
	pathLength   int
	shortestPath []int
	Edges        map[string]int `yaml:"Edges"`
	mu           sync.Mutex
}

func main() {
	var targetPath string
	flag.StringVar(&targetPath, "file", "", "File path to read from; defaults to empty and generates a new yaml")
	flag.Parse()

	if targetPath == "" {
		input := takeUserInput("No file provided: y to manually generate a case, g to autogenerate a case, n to quit: ")
		input = strings.ToLower(input)
		for {

			if input == "y" {
				generator := createFileGenerator()
				targetPath = generator.createInputYaml()
			} else if input == "n" {
				log.Println("Ending process")
				return
			} else if input == "g" {
				fmt.Println("todo: implement autogen")
				generator := createFileGenerator()
				targetPath = generator.createInputYaml()
			}
		}
	}
	var wg sync.WaitGroup
	wg.Add(1)
	cfg := yamlGraphGen(targetPath)
	go cfg.findShortestPath(cfg.SourceNode, 0, []int{}, &wg)
	wg.Wait()
	cfg.displayResult()
}

func yamlGraphGen(targetPath string) *pathingConfig {
	pc := pathingConfig{}
	data, err := ioutil.ReadFile(targetPath)
	if err != nil {
		log.Fatal("error, file path is not viable")
	}
	if err := yaml.Unmarshal(data, &pc); err != nil {
		log.Fatalf("problem with yaml unmarshalling: %v", err)
	}
	pc.pathLength = int(math.MaxInt64)
	pc.shortestPath = []int{pc.SourceNode}
	return &pc

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

func getSyncMapKey(a, b int) string {
	builder := strings.Builder{}
	builder.WriteString(strconv.Itoa(a))
	builder.WriteString("-")
	builder.WriteString(strconv.Itoa(b))

	return builder.String()
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

func (cfg *pathingConfig) displayResult() {
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
			builder.WriteString("\nShortest distance from ")
			builder.WriteString(strconv.Itoa(cfg.SourceNode))
			builder.WriteString(" to ")
			builder.WriteString(strconv.Itoa(cfg.DestNode))
			builder.WriteString(": ")
			builder.WriteString(strconv.Itoa(cfg.pathLength))
			builder.WriteString(strconv.Itoa(cfg.pathLength))
		}
	}
	fmt.Println(builder.String())
}
