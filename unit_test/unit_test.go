package test
import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
	"../utils"
	"github.com/sajari/fuzzy"
	"io/ioutil"
	"../config"
	"reflect"
)

var file []byte

func TestSpellCheck(t *testing.T) {
	model, _ := fuzzy.Load("../spellcheck/model")
	assert.Equal(t, "scholarship application", utils.SpellCheck(model, "schlarsrhip aplictaion"))
	assert.Equal(t, "scholarship application", utils.SpellCheck(model, " schlarsrhip  aplictaion "))
	assert.Equal(t, "hostel applications", utils.SpellCheck(model, "hostel applications"))
}

func TestContains(t *testing.T) {
	list := []string{"string1", "string2", "string 3"}
	assert.Equal(t, true, utils.Contains(list, "string1"))
	assert.Equal(t, true, utils.Contains(list, "string2"))
	assert.Equal(t, true, utils.Contains(list, "string 3"))
	assert.Equal(t, false, utils.Contains(list, "string 1"))
	assert.Equal(t, false, utils.Contains(list, "string3"))
	assert.Equal(t, false, utils.Contains(list, ""))
	assert.Equal(t, false, utils.Contains(list, " string2"))
	assert.Equal(t, false, utils.Contains(list, "string2 "))

	var emptyList []string
	assert.Equal(t, false, utils.Contains(emptyList, ""))
	assert.Equal(t, false, utils.Contains(emptyList, "string1"))
}

func TestGetServerConfig(t *testing.T) {
	conf := utils.GetServerConfig(file)
	assert.Equal(t, false, conf.IsProduction)
	assert.Equal(t, 8080, conf.Port)
	assert.Equal(t, "/etc/letsencrypt/live/www.test.org/fullchain.pem", conf.CertFile)
	assert.Equal(t, "/etc/letsencrypt/live/www.test.org/privkey.pem", conf.KeyFile)
	assert.Equal(t, true, conf.UseSpellchecker)

	expected := [][]string{
		{"IsProduction", "bool"},
		{"Port", "int"},
		{"CertFile", "string"},
		{"KeyFile", "string"},
		{"UseSpellchecker", "bool"}}

	serverStruct := reflect.ValueOf(&conf).Elem()
	fields := serverStruct.Type()
	assert.Equal(t, len(expected), fields.NumField())
	for i := 0; i < serverStruct.NumField(); i++ {
		f := serverStruct.Field(i)
		assert.Equal(t, expected[i][0], fields.Field(i).Name)
		assert.Equal(t, expected[i][1], f.Type().String())
	}
}

func TestGetDialogflowConfig(t *testing.T) {
	conf := utils.GetDialogflowConfig(file)
	assert.Equal(t, 2, len(conf.Agents))

	faqs := config.GetAgentConfigByName(&conf, "faqs")
	assert.Equal(t, "faqs", faqs.Name)
	assert.Equal(t, "abc123", faqs.Token)
	assert.Equal(t, "en", faqs.QueryLang)
	assert.Equal(t, "en-US", faqs.SpeechLang)

	masterbot := config.GetAgentConfigByName(&conf, "masterbot")
	assert.Equal(t, "masterbot", masterbot.Name)
	assert.Equal(t, "123aaa", masterbot.Token)
	assert.Equal(t, "en", masterbot.QueryLang)
	assert.Equal(t, "en-US", masterbot.SpeechLang)

	expected := [][]string{
		{"Name", "string"},
		{"Token", "string"},
		{"QueryLang", "string"},
		{"SpeechLang", "string"}}

	agentStruct := reflect.ValueOf(faqs).Elem()
	fields := agentStruct.Type()
	assert.Equal(t, len(expected), fields.NumField())
	for i := 0; i < agentStruct.NumField(); i++ {
		f := agentStruct.Field(i)
		assert.Equal(t, expected[i][0], fields.Field(i).Name)
		assert.Equal(t, expected[i][1], f.Type().String())
	}
}

func TestGetExternalAgentConfig(t *testing.T) {
	conf := utils.GetExternalAgentConfig(file)
	assert.Equal(t, "http://10.255.47.200:5000/predict", conf.ClassifierUrl)

	expected := [][]string{
		{"ClassifierUrl", "string"}}

	externalAgentStruct := reflect.ValueOf(&conf).Elem()
	fields := externalAgentStruct.Type()
	assert.Equal(t, len(expected), fields.NumField())
	for i := 0; i < externalAgentStruct.NumField(); i++ {
		f := externalAgentStruct.Field(i)
		assert.Equal(t, expected[i][0], fields.Field(i).Name)
		assert.Equal(t, expected[i][1], f.Type().String())
	}
}

func TestGetGoogleSearchConfig(t *testing.T) {
	conf := utils.GetGoogleSearchConfig(file)
	assert.Equal(t, "123123123", conf.SearchEngineID)
	assert.Equal(t, "aaa123aaa", conf.ApiKey)

	expected := [][]string{
		{"SearchEngineID", "string"},
		{"ApiKey", "string"}}

	googleSearchStruct := reflect.ValueOf(&conf).Elem()
	fields := googleSearchStruct.Type()
	assert.Equal(t, len(expected), fields.NumField())
	for i := 0; i < googleSearchStruct.NumField(); i++ {
		f := googleSearchStruct.Field(i)
		assert.Equal(t, expected[i][0], fields.Field(i).Name)
		assert.Equal(t, expected[i][1], f.Type().String())
	}
}
func TestMain(m *testing.M) {
	fmt.Println("Starting unit test...")
	file, _ = ioutil.ReadFile("./config.json")
	m.Run()
}