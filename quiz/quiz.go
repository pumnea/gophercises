// Package quiz provides utilities to parse a csv file.
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

type Problem struct {
	question string
	answer   string
}

func main() {
	// Parse command-line flags
	filePath := flag.String("file", "problems.csv", "path to a CSV file in format 'question,answer'")
	flag.Parse()

	// Open and read the CSV file
	file, err := os.Open(*filePath)
	if err != nil {
		exit(fmt.Errorf("opening file: %w", err))
	}
	defer file.Close()

	lines, err := readFile(file)
	if err != nil {
		exit(fmt.Errorf("reading CSV: %w", err))
	}

	// Parse problems
	problems, err := parseProblems(lines)
	if err != nil {
		exit(err)
	}

	// Shuffle the questions
	shuffleProblems(problems)

	// Run the quiz
	correct := runQuiz(problems, os.Stdin, os.Stdout)
	fmt.Printf("\nScore: %d correct out of %d total\n", correct, len(problems))
}

// readCSV reads a CSV file from an io.Reader and returns its contents.
func readFile(r io.Reader) ([][]string, error) {
	reader := csv.NewReader(r)

	// Read all lines
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse csv file: %w", err)
	}

	return lines, nil
}

// parseProblems converts CSV lines into a slice of Problem structs.
func parseProblems(lines [][]string) ([]Problem, error) {
	problems := make([]Problem, 0, len(lines))

	for _, line := range lines {
		problems = append(problems, Problem{line[0], strings.TrimSpace(line[1])})
	}
	return problems, nil
}

// runQuiz administers the quiz and returns the number of correct answers.
func runQuiz(problems []Problem, input io.Reader, output io.Writer) int {
	scanner := bufio.NewScanner(input)
	correct := 0

	for i, p := range problems {
		fmt.Fprintf(output, "%d> %s = ", i+1, p.question)
		scanner.Scan()
		answer := strings.TrimSpace(scanner.Text())
		if answer == p.answer {
			correct++
		}
	}
	return correct
}

// shuffleProblems randomizes the order of quiz problems.
func shuffleProblems(problems []Problem) {
	rand.Shuffle(len(problems), func(i int, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})
}

// exit prints an error message to stderr and terminates the program.
func exit(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
