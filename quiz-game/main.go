package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// present quiz problems
// accepts user input
// check for correctness
func main() {
	// flag package is a way to define CLI flags, here is a -csv flag, default value is problems.csv and default text is the 3th arg
	csvFilename := flag.String("csv", "problems.csv", "a csv file i the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	// flag return a pointer, so we pass the value of the pointer here.
	file, err := os.Open(*csvFilename)
	if err != nil {
		// handling error, Sprintf will format a string but not print it. So it will be passed to the exit func
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	// once the file is open, let's create a CSV reader
	// it accepts an io.Reader interface, which our file already is
	r := csv.NewReader(file)
	// let's parse the CSV
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	problems := parseLines(lines)

	// define the timer, uses the timer package and when it expires it sends a message to a channel
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	// start couting how many answers are correct
	correct := 0
	// print the problems to the user
	for i, p := range problems {
		fmt.Printf("Problem %d: %s = \n", i+1, p.question)
		// below we defined a go routine to track if a user inputed an answer
		// we defined a goroutine because fmt.Scanf would block the code, so even if the timer expires, the user could
		// put an answer. By using the goroutine, fmt.Scanf is non blocking
		answerCh := make(chan string)
		go func() {
			// read an answer
			// define a variable to store the user answer
			var answer string
			// scan text from user input. Since the answer is a single number, this is fine. But if the answer is a multiple
			// text, it will give error because Scanf delete all spaces
			// we use the pointer of the answer because when the user types his answer, we want to modify the value of the variable
			fmt.Scanf("%s\n", &answer)
			// and we send the user answer into the outside variable, so it's acessible outside the goroutine
			answerCh <- answer
		}()
		// the select inside the for loop says: if we receive a message in the timer channel, stop the program and print the result. Or if we receive a message in the answerCh, it means the user answered the question before the timer expires.
		select {
		// a linha abaixo cria um channel e escuta esse channel, esperando uma mensagem chegar.
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
			return
		// if there is an user answer, check for correctness
		case answer := <-answerCh:
			if answer == p.answer {
				correct++
				timer = time.NewTimer(time.Duration(*timeLimit) * time.Second)
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

// create the problems struct below, by taken the CSV (which is a slice of slices) and transforming into the problem struct
func parseLines(lines [][]string) []problem {
	// declare the variable we will return, a slice of problems with the length of total number of lines in the csv file
	// we assume every row in the CSV file is a problem
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			// in case the CSV comes with spaces in the answer, we want to delete the spaces because we are using fmt.Scanf
			// to accept user answer, and Scanf trim all spaces
			answer: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

// let's create a struct for using this struct in the rest of our code, so our the code always expect a problem struct
type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	// if user inputs something like -csv=ab.csv, the user will se the name of that file in the error
	fmt.Println(msg)
	os.Exit(1)
}
