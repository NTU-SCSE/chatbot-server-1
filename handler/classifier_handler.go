package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"../config"
	"../utils"
	"github.com/marcossegovia/apiai-go"
	"github.com/tidwall/gjson"
)

func getClass(query string, extAgentsConf *config.ExternalAgentsConfig) string {
	defer utils.TimeFunction(time.Now(), "classify")
	var jsonStr = []byte(`{"question":"` + query + `"}`)
	request, err := http.NewRequest("POST", extAgentsConf.ClassifierUrl, bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	classifierResp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer classifierResp.Body.Close()

	body, _ := ioutil.ReadAll(classifierResp.Body)
	class := gjson.Get(string(body[:]), "result").String()
	fmt.Println("--------------------")
	fmt.Println(class)
	return class
}

func NewClassifierWebhookHandler(conf *config.DialogflowConfig, extAgentsConf *config.ExternalAgentsConfig) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		defer utils.TimeFunction(time.Now(), "master")
		body, _ := ioutil.ReadAll(req.Body)
		fullJSON := string(body[:])
		originalRequest := gjson.Get(fullJSON, "result.resolvedQuery")
		sessionID := gjson.Get(fullJSON, "sessionId")

		// Send request to classifier
		class := getClass(originalRequest.String(), extAgentsConf)

		// Send query to the respective agent
		// TODO: Fix this when we have proper classifier integrations and agent names
		classifierResult := class // temporary var for debugging purpose
		class = "faqs"
		agentConfig := config.GetAgentConfigByName(conf, class)
		apiaiClient, _ := apiai.NewClient(
			&apiai.ClientConfig{
				Token:      agentConfig.Token,
				QueryLang:  agentConfig.QueryLang,
				SpeechLang: agentConfig.SpeechLang,
			},
		)

		qr, _ := apiaiClient.Query(apiai.Query{Query: []string{originalRequest.String()}, SessionId: sessionID.String()})

		resultMap := make(map[string]interface{})

		resultMap["speech"] = qr.Result.Fulfillment.Speech
		resultMap["source"] = class + "_agent:" + classifierResult

		resultJson, _ := json.Marshal(resultMap)

		rw.Header().Set("Content-Type", "application/json")

		rw.Write(resultJson)
	}
}
