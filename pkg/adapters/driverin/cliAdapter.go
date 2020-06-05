package driverin

import "C"
import (
	"bufio"
	"flag"
	"fmt"
	"gophercise1/app/quiz"
	"gophercise1/app/quiz/structs"
	"log"
	"os"
	"strings"
	"sync"
)

type CLIAdapter struct {
	Quiz quiz.Port
}

func NewCLIAdapter(quiz quiz.Port) *CLIAdapter {
	return &CLIAdapter{
		Quiz: quiz,
	}
}

func (C *CLIAdapter) Run() {
	var wg sync.WaitGroup
	var fileLocation string
	var timeLimit int
	parseFlags(&fileLocation, &timeLimit)

	err := C.Quiz.Setup(fileLocation)
	if err != nil {
		log.Fatal(err)
	}

	questionCh := C.Quiz.Questions()
	answerCh := make(chan int)

	wg.Add(1)
	go C.getAnswers(&wg, answerCh)

	go startQuestions(timeLimit, answerCh, questionCh)

	wg.Wait()
}

func startQuestions(timeLimit int, answerCh chan int, questions <-chan structs.Question) {
	timesUp := C.Quiz.StartTimer(timeLimit)
	scanner := bufio.NewReader(os.Stdin)
loop:
	for {
		select {
		case <-timesUp:
			close(answerCh)
			fmt.Println("out of time")
			break loop
		case question, more := <-questions:
			if !more {
				close(answerCh)
				break loop
			} else {
				fmt.Printf("%v) %v : ", question.QuestionNumber, question.Question)
				str, _ := scanner.ReadString('\n')
				if strings.TrimSpace(str) == question.Answer {
					answerCh <- 1
				}
			}
		}
	}
}

func (c *CLIAdapter) getAnswers(wg *sync.WaitGroup, answerCh chan int) {
	defer wg.Done()
loop:
	for {
		select {
		case _, ok := <-answerCh:
			if ok {
				C.Quiz.IncrementScore()
			} else {
				fmt.Printf("You scored: %v/%v", C.Quiz.Score(), C.Quiz.NumOfQuestions())
				break loop
			}
		default:
		}
	}
}

func parseFlags(fileLocation *string, timeLimit *int) {
	flag.StringVar(fileLocation, "fileLocation", "problems.csv", "location of the csv file")
	flag.IntVar(timeLimit, "timeLimit", 30, "time limit of quiz")
	flag.Parse()
}
