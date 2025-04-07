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
	"time"
)

type Problem struct {
	question string
	answer   string
}

func main() {
	// Parse command-line flags
	filePath := flag.String("file", "problems.csv", "path to a CSV file in format 'question,answer'")
	timeLimit := flag.Int("limit", 30, "time limit for quiz in seconds")
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

	// Create timer
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	defer timer.Stop()

	// Run the quiz
	correct, attempted := runQuiz(problems, os.Stdin, os.Stdout, timer)
	fmt.Printf("\nScore: %d correct out of %d attempted (total questions: %d)\n",
		correct, attempted, len(problems))
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
func runQuiz(problems []Problem, input io.Reader, output io.Writer, timer *time.Timer) (correct, attempted int) {
	scanner := bufio.NewScanner(input)

	answerCh := make(chan string)

	for i, p := range problems {
		fmt.Fprintf(output, "%d> %s = ", i+1, p.question)

		go func() {
			if scanner.Scan() {
				answerCh <- strings.TrimSpace(scanner.Text())
			} else {
				answerCh <- ""
			}
		}()

		select {
		case <-timer.C:
			fmt.Fprintln(output, "\nTime's up!")
			return correct, attempted
		case answer := <-answerCh:
			attempted++
			if strings.EqualFold(answer, p.answer) {
				correct++
			}
		}
	}

	return correct, attempted
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
