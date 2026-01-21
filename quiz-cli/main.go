package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	filepath := flag.String("file", "problems.csv", "path to csv file")
	flag.Parse()
	file, err := os.Open(*filepath)
	if err != nil {
		fmt.Printf("Failed to open file: %s\n", *filepath)
		os.Exit(1)
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		os.Exit(1)
	}

	inputReader := bufio.NewReader(os.Stdin)
	correct := 0
	for i, row := range records {
		fmt.Printf("Problem #%d %s=", i, row[0])
		input, err := inputReader.ReadString('\n')
		if err != nil {
			os.Exit(1)
		}
		if strings.TrimSpace(input) == row[1] {
			correct++
		}
	}
	fmt.Printf("You've got %d right and %d wrong.\n", correct, len(records)-correct)
}
