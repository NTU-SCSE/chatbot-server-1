package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"./config"
	"./course"
	"./handler"
	"./storage"
	"./utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sajari/fuzzy"
)

// "github.com/kamalpy/apiai-go"

type param map[string]interface{}

type context map[string]interface{}

type metadata struct {
	IntentName string `json:intentName`
}

type query_struct struct {
	Query     string
	SessionID string
	Enum      []string `json:",omitempty"`
}

type response struct {
	Response string
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world! Your URL: %s", r.URL.Path[1:])
}

func main() {
	// spell checking module
	fmt.Println("Loading Spell Checking module...")
	model, _ := fuzzy.Load("spellcheck/model")

	// Read configuration file
	fmt.Println("Reading Configurations...")

	file, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		fmt.Println(err.Error())
	}

	conf := utils.GetServerConfig(file)
	googleConf := utils.GetGoogleSearchConfig(file)
	dialogflowConf := utils.GetDialogflowConfig(file)
	externalAgents := utils.GetExternalAgentConfig(file)

	// Loading database
	fmt.Println("Loading database...")
	db, err := storage.NewDB("test.sqlite3")
	course := course.NewCourse()

	// Get the data of courses from json files
	//temp := []string{"course description", "course name", "au", "prereq", "course code", "time", "venue"}
	//result := utils.GetEnum(temp)

	// start server
	fmt.Println("Starting the server...")
	r := mux.NewRouter()

	r.HandleFunc("/", defaultHandler)
	r.HandleFunc("/query", handler.NewQueryHandler(course))
	r.HandleFunc("/webhook", handler.WebhookHandler)
	r.HandleFunc("/webhook-v1", handler.NewWebhookHandlerV1(&googleConf, db, conf.UseSpellchecker))
	r.HandleFunc("/classifier-webhook", handler.NewClassifierWebhookHandler(&dialogflowConf, &externalAgents))
	r.HandleFunc("/internal-query", handler.NewInternalHandler(config.GetAgentConfigByName(&dialogflowConf, "faqs")))
	r.HandleFunc("/spellcheck", handler.NewSpellCheckHandler(config.GetAgentConfigByName(&dialogflowConf, "faqs"), model))

	// Apply the CORS middleware to our top-level router, with the defaults.
	if conf.IsProduction {
		fmt.Println("Starting a production server on port %d...\n", conf.Port)
		err := http.ListenAndServeTLS(":"+strconv.Itoa(conf.Port), conf.CertFile, conf.KeyFile, cors.Default().Handler(r))
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Starting a development local server on port %d...\n", conf.Port)
		http.ListenAndServe(":"+strconv.Itoa(conf.Port), cors.Default().Handler(r))
	}
}

// todo: fix typo in application security json data
// todo: fix computer security entity
// todo: fix ce1004
