package test
import (
	"testing"
	"encoding/json"
	"io/ioutil"
	"../config"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/marcossegovia/apiai-go"
	"strings"
	"../utils"
)

var dialogflowConf config.DialogflowConfig
var testConfig config.TestConfig

type query_response struct {
	Query string `json="query"`
	Response string `json="response"`
}

var queries []query_response


func TestDialogFlow(t *testing.T) {
	agentConfig := config.GetAgentConfigByName(&dialogflowConf, "faqs")
	client, err := apiai.NewClient(
		&apiai.ClientConfig{
			Token:      agentConfig.Token,
			QueryLang:  agentConfig.QueryLang,
			SpeechLang: agentConfig.SpeechLang,
		},
	)

	if err != nil {
		fmt.Printf("%v", err)
	}


	//Set the query string and your current user identifier.
	var qr *apiai.QueryResponse
	testSessionID := "10101010-1010-1010-1010-101010101010"
	for _, query := range queries {
		qr, err = client.Query(apiai.Query{Query: []string{query.Query}, SessionId: testSessionID})
		expected := query.Response
		actual := qr.Result.Fulfillment.Speech
		if !testConfig.IsExactMatching {
			expected = strings.Replace(expected, "\n", "", -1)
			expected = strings.Replace(expected, "\r", "", -1)
			expected = strings.Replace(expected, " ", "", -1)
			actual = strings.Replace(actual, "\n", "", -1)
			actual = strings.Replace(actual, "\r", "", -1)
			actual = strings.Replace(actual, " ", "", -1)
		}
		assert.Equal(t, expected, actual)
	}
	
	
}

func TestMain(m *testing.M) {
	// Read configuration file
	fmt.Println("Reading Configurations...")
	

	file, err := ioutil.ReadFile("./test_config.json")
	if err != nil {
		fmt.Println(err.Error())
	}

	dialogflowConf = utils.GetDialogflowConfig(file)
	json.Unmarshal(file, &testConfig)

	queryResponse, err := ioutil.ReadFile("./query_response.json")
	json.Unmarshal(queryResponse, &queries)
	m.Run()
}