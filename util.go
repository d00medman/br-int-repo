package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func takeUserInput(message string) string {
	fmt.Printf(message)
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("An error occured while reading input. Please try again", err)
	}
	input = strings.TrimSuffix(input, "\n")
	return input
}

func userInputToInt(message string) int {
	input := takeUserInput(message)

	output, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error in strconv in yamlGen")
	}
	return output
}

func (cfg *pathingConfig) findPointDistance(a, b int) int {
	k := getSyncMapKey(a, b)
	cfg.mu.Lock()
	dist := cfg.Edges[k]
	cfg.mu.Unlock()
	return dist
}
