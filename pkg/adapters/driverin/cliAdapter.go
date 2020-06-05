package driverin

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
	"time"
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
	scanner := bufio.NewReader(os.Stdin)

	err := C.Quiz.Setup(fileLocation)
	if err != nil {
		log.Fatal(err)
	}

	questions := C.Quiz.Questions()
	timesUp := C.Quiz.StartTimer(timeLimit)
	answerCh := make(chan int)

	wg.Add(1)
	go getAnswers(&wg, answerCh, C)

	go startQuestions(timesUp, answerCh, questions, scanner)

	wg.Wait()
}

func startQuestions(timesUp <-chan time.Time, answerCh chan int, questions <-chan structs.Question, scanner *bufio.Reader) {
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

func getAnswers(wg *sync.WaitGroup, answerCh chan int, C *CLIAdapter) {
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
