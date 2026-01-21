package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	filePath := flag.String("file", "problems.csv", "path to csv file")
	timeLimit := flag.Int("limit", 10, "time limit")
	flag.Parse()
	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %s\n", *filePath)
		os.Exit(1)
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		os.Exit(1)
	}
	problems := parseRecords(records)
	inputReader := bufio.NewReader(os.Stdin)
	t := time.NewTimer(time.Second * time.Duration(*timeLimit))
	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d %s=", i, p.q)

		answerChan := make(chan string)
		go func() {
			input, err := inputReader.ReadString('\n')
			if err != nil {
				fmt.Println("erroring")
			}
			answerChan <- strings.TrimSpace(input)
		}()

		select {
		case <-t.C:
			fmt.Printf("\nYou've got %d right out of %d.\n", correct, len(records))
			return
		case answer := <-answerChan:
			if answer == p.a {
				correct++
			}
		}
	}
	fmt.Printf("You've got %d right out of %d.\n", correct, len(records))
}

func parseRecords(r [][]string) []problem {
	p := make([]problem, len(r))
	for i, row := range r {
		p[i] = problem{
			q: strings.TrimSpace(row[0]),
			a: strings.TrimSpace(row[1]),
		}
	}
	return p
}
