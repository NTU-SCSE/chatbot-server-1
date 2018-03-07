package handler
import (
    "net/http"
    "io/ioutil"
    "encoding/json"
	"github.com/marcossegovia/apiai-go"    
	"github.com/tidwall/gjson"
	"../utils"
	"time"
)

func DummyWebhookHandler(rw http.ResponseWriter, req *http.Request) {
    defer utils.TimeFunction(time.Now(), "d")
    body, _ := ioutil.ReadAll(req.Body)
    fullJSON := string(body[:])
    originalRequest := gjson.Get(fullJSON, "result.resolvedQuery")
    client, _ := apiai.NewClient(
        &apiai.ClientConfig{
			// Token:      "c49f70c867c54ff49d054455b3153e61 ",
			Token:      "031636d290f341729417585f09f1ebc4",
            QueryLang:  "en",    //Default en
            SpeechLang: "en-US", //Default en-US
        },
	)
	sessionID := "456"
    // _, _ = client.Query(apiai.Query{Query: []string{"do you know me?"}, SessionId: sessionID})
    _, _ = client.Query(apiai.Query{Query: []string{originalRequest.String()}, SessionId: sessionID})
    
    resultMap := make(map[string]interface{})
    
       
    resultMap["displayText"] = "Test Response"
    resultMap["speech"] = "Response not found"
    resultMap["data"] = ""
    resultMap["contextOut"] = []string{}
    resultMap["source"] = "Hello"

    resultJson, _ := json.Marshal(resultMap)
    
    rw.Header().Set("Content-Type", "application/json")
        
    rw.Write(resultJson)
}