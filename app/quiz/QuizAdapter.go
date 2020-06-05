package quiz

import (
	"fmt"
	"gophercise1/app/quiz/structs"
	"gophercise1/pkg/adapters/drivenout/csv"
	"gophercise1/pkg/adapters/drivenout/storage"
	"time"
)

type Quiz struct {
	csv                 csv.Port
	Storage             storage.Port
	numOfCorrectAnswers int
	numOfQuestions      int
}

func NewQuiz(csv csv.Port, storage storage.Port) *Quiz {
	return &Quiz{
		csv:     csv,
		Storage: storage,
	}
}

func (q *Quiz) Setup(fileLocation string) error {
	if filePort, ok := q.csv.(csv.FilePort); ok {
		filePort.SetFileLocation(fileLocation)
	}
	records, err := q.csv.GetRecords()
	if err != nil {
		return fmt.Errorf("starting quiz: %w", err)
	}
	q.numOfQuestions = len(records)
	q.Storage.SaveRecords(records)
	return nil
}

func (q *Quiz) StartTimer(timeLimit int) <-chan time.Time {
	return time.After(time.Second * time.Duration(timeLimit))
}

func (q *Quiz) NumOfQuestions() int {
	return q.numOfQuestions
}

func (q *Quiz) Questions() <-chan structs.Question {
	questions := make(chan structs.Question)

	go func() {
		for index, question := range q.Storage.Records() {
			questions <- structs.Question{
				QuestionNumber: index + 1,
				Qustion:        question[0],
				Answer:         question[1],
			}
		}
		close(questions)
	}()
	return questions
}

func (q *Quiz) IncrementScore() {
	q.numOfCorrectAnswers++
}

func (q *Quiz) Score() int {
	return q.numOfCorrectAnswers
}
