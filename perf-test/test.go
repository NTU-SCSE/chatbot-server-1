package main

import (
	"github.com/marcossegovia/apiai-go"    
	"time"
	"log"
    "strconv"
    "fmt"
)

type param map[string]interface{}

type context map[string]interface{}

type metadata struct {
    IntentName string `json:intentName`
}

type query_struct struct {
    Query string
    SessionID string
    Enum []string   `json:",omitempty"`
}

type response struct {
    Response string
}

type course struct {
    Code string `json:"code"`
    Name string `json:"name"`
    AU int `json:"AU"`
    PreReq string `json:"preReq"`
    Description string `json:"description"`
}

type class struct {
    Code string `json:"code"`
    Index string `json:"index"`
    Type string `json:"type"`
    Group string `json:"group"`
    Day string `json:"day"`
    Time string `json:"time"`
    Venue string `json:"venue"`
    Remark string `json:"remark"`
}
func timeFunction(start time.Time, name string) {
    elapsed := time.Since(start)
    log.Printf("%s took %s seconds.", name, fmt.Sprintf("%.9f", elapsed.Seconds()))
}

func testQueryOnly(iteration int) {
	defer timeFunction(time.Now(), "testQueryOnly with " + strconv.Itoa(iteration))
	client, _ := apiai.NewClient(
        &apiai.ClientConfig{
			Token:      "c49f70c867c54ff49d054455b3153e61 ", // masterbot
            // Token:      "031636d290f341729417585f09f1ebc4", // SCSE BOT
            // Token: "58be6f8f4fb9447693edd36fb975bece", // Chatbot.v1
            QueryLang:  "en",    //Default en
            SpeechLang: "en-US", //Default en-US
        },
	)
    // sessionID := "123"
    sessionID := "105"
	for i := 0 ; i < iteration; i++ {
		// _, _ = client.Query(apiai.Query{Query: []string{"do you know me?"}, SessionId: sessionID})
		_, _ = client.Query(apiai.Query{Query: []string{"SCSE rank"}, SessionId: sessionID})
	}
}

func main() {
    testQueryOnly(70)
}
