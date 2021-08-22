package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type inputFileGenerator struct {
	nodeLabels        []int  // `yaml:"SourceNode,omitempty"`
	NodeDistanceRange int    `yaml:"NodeDistanceRange,omitempty"`
	destNode          int    // `yaml:"SourceNode,omitempty"`
	sourceNode        int    // `yaml:"SourceNode,omitempty"`
	TotalNodes        int    `yaml:"TotalNodes,omitempty"`
	NodeLabelMax      int    `yaml:"NodeLabelMax,omitempty"`
	FileLocation      string `yaml:"FileLocation ,omitempty"`
	FileName          string `yaml:"FileName,omitempty"`
}

func createFileGeneratorFromYaml(targetPath string) inputFileGenerator {
	gen := inputFileGenerator{}
	data, err := ioutil.ReadFile(targetPath)
	if err != nil {
		log.Fatal("error, file path is not viable")
	}
	if err := yaml.Unmarshal(data, &gen); err != nil {
		log.Fatalf("problem with generator yaml unmarshalling: %v", err)
	}
	labels := []int{}
	for i := 0; i < gen.TotalNodes; i++ {
		labels = append(labels, i)
	}
	fmt.Printf("Input yaml will be generated with the following labels: %v\n", labels)
	sourceIndex := rand.Intn(gen.TotalNodes)
	destIndex := rand.Intn(gen.TotalNodes)
	for sourceIndex == destIndex {

		destIndex = rand.Intn(gen.TotalNodes)
	}
	gen.sourceNode = labels[sourceIndex]
	gen.destNode = labels[destIndex]
	gen.nodeLabels = labels

	return gen
}

func createFileGenerator() inputFileGenerator {
	inLabels := func(labels []int, target int) bool {
		for _, i := range labels {
			if i == target {
				return true
			}
		}
		return false
	}
	labels := []int{}
	totalNodes := userInputToInt("How many nodes do you want to create: ")
	//todo: remove me when we find out
	fmt.Println("what happens when you hit enter with this function?", totalNodes)
	fmt.Printf("Please provide numbers from %d to %d to act as the distance range between nodes\n", int(math.MinInt64), int(math.MaxInt64))
	for i := 0; i < totalNodes; i++ {

		nli := userInputToInt(fmt.Sprintf("Node: %d: ", i))
		for inLabels(labels, nli) {
			fmt.Printf("there is alreaady a node with label %d in the set of nodes %v, please provide a different number.\n", nli, labels)
			nli = userInputToInt(fmt.Sprintf("Node: %d: ", i))
		}
		labels = append(labels, nli)
		fmt.Printf("Input yaml will be generated with the following labels: %v\n", labels)
	}

	source := userInputToInt("Please choose a source node: ")
	for !inLabels(labels, source) {
		source = userInputToInt(fmt.Sprintf("Please choose a number from %v: ", labels))
	}
	dest := userInputToInt("Please choose a destination node: ")
	for dest == source || !inLabels(labels, source) {

		dest = userInputToInt(fmt.Sprintf("Please choose a node from %v other than %d: ", labels, source))
	}

	return inputFileGenerator{
		nodeLabels:        labels,
		NodeDistanceRange: userInputToInt("Please provide the maximum distance between each node: "),
		sourceNode:        source,
		destNode:          dest,
	}
}

func (g *inputFileGenerator) createInputYaml(fullPath string) {
	file, err := os.Create(fullPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	yamlBuilder := strings.Builder{}
	yamlBuilder.WriteString("SourceNode: ")
	yamlBuilder.WriteString(strconv.Itoa(g.sourceNode))
	yamlBuilder.WriteString("\n")
	yamlBuilder.WriteString("DestNode: ")
	yamlBuilder.WriteString(strconv.Itoa(g.destNode))
	yamlBuilder.WriteString("\n")
	yamlBuilder.WriteString("Edges:\n")
	for _, i := range g.nodeLabels {
		for _, j := range g.nodeLabels {
			if i == j {
				continue
			}
			yamlBuilder.WriteString("    ")
			yamlBuilder.WriteString(strconv.Itoa(i))
			yamlBuilder.WriteString("-")
			yamlBuilder.WriteString(strconv.Itoa(j))
			yamlBuilder.WriteString(": ")

			yamlBuilder.WriteString(strconv.Itoa(rand.Intn(g.NodeDistanceRange)))
			yamlBuilder.WriteString("\n")
		}
	}
	_, err = file.WriteString(yamlBuilder.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("---------- Case file successfully generated ----------")
}
func userGenerateSavePath() string {
	directory := takeUserInput("Provide a directory to place the new yaml (blank input creates file locally) : ")
	fileName := takeUserInput("Provide a name for the new yaml file (blank input defaults to 'autogen'): ")
	return findFullSavePath(fileName, directory)
}
func findFullSavePath(fileName, directory string) string {
	if fileName == "" {
		fileName = "autogen"
	}
	if directory == "" {
		directory = "./"
	}
	fullPathBuilder := strings.Builder{}

	fullPathBuilder.WriteString(directory)
	fullPathBuilder.WriteString(fileName)
	basePath := fullPathBuilder.String()
	fullPathBuilder.WriteString(".yaml")
	tmpSlug := 0
	fullPath := fullPathBuilder.String()
	_, err := os.Stat(fullPath)
	for !os.IsNotExist(err) {

		fullPathBuilder.Reset()
		fullPathBuilder.WriteString(basePath)
		fullPathBuilder.WriteString(strconv.Itoa(tmpSlug))
		fullPathBuilder.WriteString(".yaml")
		fullPath = fullPathBuilder.String()
		_, err = os.Stat(fullPath)
		tmpSlug += 1
	}
	return fullPath
}
