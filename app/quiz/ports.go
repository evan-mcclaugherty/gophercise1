package quiz

import (
	"gophercise1/app/quiz/structs"
	"time"
)

type Port interface {
	Setup(fileLocation string) error
	StartTimer(timeLimit int) (done <-chan time.Time)
	Questions() <-chan structs.Question
	NumOfQuestions() int
	IncrementScore()
	Score() int
}
