package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"../config"
	"../utils"
	"github.com/marcossegovia/apiai-go"
	"github.com/sajari/fuzzy"
)

func NewSpellCheckHandler(conf *config.AgentConfig, model *fuzzy.Model) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		defer utils.TimeFunction(time.Now(), "spell")
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			panic(err)
		}

		resultMap := make(map[string]interface{})

		var t query_struct

		err = json.Unmarshal(body, &t)
		if err != nil {
			panic(err)
		}

		t.Query = utils.SpellCheck(model, t.Query)

		client, err := apiai.NewClient(
			&apiai.ClientConfig{
				Token:      conf.Token,
				QueryLang:  conf.QueryLang,
				SpeechLang: conf.SpeechLang,
			},
		)

		if err != nil {
			fmt.Printf("%v", err)
		}

		// log into file
		f, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("error %v", err)
		}
		defer f.Close()

		f.WriteString("Query from: " + "11111111-1111-1111-1111-111111111111" + "\r\n")
		f.WriteString(t.Query + "\r\n")
		f.WriteString("----------------------\r\n")

		//Set the query string and your current user identifier.
		var qr *apiai.QueryResponse
		qr, err = client.Query(apiai.Query{Query: []string{t.Query}, SessionId: "11111111-1111-1111-1111-111111111111"})

		resultMap["Result"] = qr.Result.Fulfillment.Speech

		resultJson, _ := json.Marshal(resultMap)

		rw.Header().Set("Content-Type", "application/json")

		rw.Write(resultJson)
	}
}
