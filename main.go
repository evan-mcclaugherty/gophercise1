package main

import (
	"gophercise1/app/quiz"
	"gophercise1/pkg/adapters/drivenout/csv"
	"gophercise1/pkg/adapters/drivenout/storage"
	"gophercise1/pkg/adapters/driverin"
)

func main() {
	var csvPort csv.FilePort
	csvPort = csv.NewFileAdapter()

	var storagePort storage.Port
	storagePort = storage.NewMemoryAdapter()

	var quizPort quiz.Port
	quizPort = quiz.NewQuiz(csvPort, storagePort)

	var cliAdapter driverin.Port
	cliAdapter = driverin.NewCLIAdapter(quizPort)

	cliAdapter.Run()
}
