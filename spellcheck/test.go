package main

import (
	"fmt"

	"github.com/sajari/fuzzy"
)

func main() {
	model := fuzzy.NewModel()

	// // For testing only, this is not advisable on production
	model.SetThreshold(2)

	// // This expands the distance searched, but costs more resources (memory and time).
	// // For spell checking, "2" is typically enough, for query suggestions this can be higher
	model.SetDepth(2)

	// // Train multiple words simultaneously by passing an array of strings to the "Train" function
	model.Train(fuzzy.SampleEnglish())
	words := []string{"SCSE", "SCSE", "SCSE", "SCSE", "SCSE", "SCSE", "SCSE", "SCSE", "SCSE", "SCSE"}
	model.Train(words)
	// // model := fuzzy.Load("model")

	// // use SaveLight to get smaller model, but not useful for incremental corpus
	model.Save("model")

	// model, _ := fuzzy.Load("model")

	// Check Spelling
	fmt.Println("\nSPELL CHECKS")
	fmt.Println("	Deletion test (yor) : ", model.SpellCheck("SCSE"))
}
