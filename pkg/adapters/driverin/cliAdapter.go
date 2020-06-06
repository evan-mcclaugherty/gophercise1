package driverin

import (
	"bufio"
	"flag"
	"fmt"
	"gophercise1/app/quiz"
	"log"
	"os"
	"strings"
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
	var fileLocation string
	var timeLimit int
	parseFlags(&fileLocation, &timeLimit)

	err := C.Quiz.Setup(fileLocation)
	if err != nil {
		log.Fatal(err)
	}

	questionCh := C.Quiz.Questions()
	timesUp := C.Quiz.StartTimer(timeLimit)
	scanner := bufio.NewReader(os.Stdin)

	for question := range questionCh {
		answerCh := make(chan bool)
		go func() {
			fmt.Printf("%v) %v : ", question.QuestionNumber, question.Question)
			str, _ := scanner.ReadString('\n')
			answerCh <- strings.TrimSpace(str) == question.Answer
		}()

		select {
		case <-timesUp:
			fmt.Printf("\nYou got %v our of %v right.\n", C.Quiz.Score(), C.Quiz.NumOfQuestions())
			return
		case isCorrect := <-answerCh:
			if isCorrect {
				C.Quiz.IncrementScore()
			}
		}
	}
}

func parseFlags(fileLocation *string, timeLimit *int) {
	flag.StringVar(fileLocation, "fileLocation", "problems.csv", "location of the csv file")
	flag.IntVar(timeLimit, "timeLimit", 30, "time limit of quiz")
	flag.Parse()
}
