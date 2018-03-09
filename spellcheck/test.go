package main 

import(
	"github.com/sajari/fuzzy"
	"fmt"
)

func main() {
	model := fuzzy.NewModel()
	
	// For testing only, this is not advisable on production
	model.SetThreshold(2)

	// This expands the distance searched, but costs more resources (memory and time). 
	// For spell checking, "2" is typically enough, for query suggestions this can be higher
	model.SetDepth(2)

	// Train multiple words simultaneously by passing an array of strings to the "Train" function
	// words := []string{"bob", "your", "uncle", "dynamite", "delicate", "biggest", "big", "bigger", "aunty", "you're"}
	model.Train(fuzzy.SampleEnglish())
	model.SaveLight("model_small")
	// Check Spelling
	fmt.Println("\nSPELL CHECKS")
	fmt.Println("	Deletion test (yor) : ", model.SpellCheck("yor"))
	fmt.Println("	Swap test (uncel) : ", model.SpellCheck("uncel"))
	fmt.Println("	Replace test (dynemite) : ", model.SpellCheck("dynemite"))
	fmt.Println("	Insert test (dellicate) : ", model.SpellCheck("dellicate"))
	fmt.Println("	Two char test (dellicade) : ", model.SpellCheck("dellicade"))
}