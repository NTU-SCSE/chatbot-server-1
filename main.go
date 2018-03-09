package main

import (
    "fmt"
	"net/http"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "./handler"    
    "./course"
    "./config"
    "io/ioutil"
    "encoding/json"
    "strconv"
)

// "github.com/kamalpy/apiai-go"

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

func defaultHandler(w http.ResponseWriter, r *http.Request) {
    // TODO: Handle this properly.
	fmt.Fprintf(w, "Hello, world! Your URL: %s", r.URL.Path[1:])
}


func main() {
    // Read configuration file
    var conf config.ServerConfig
    var googleConf config.GoogleSearchConfig
    file, err := ioutil.ReadFile("./config/config.json")
    if err != nil {
        fmt.Println(err.Error())
    }
    json.Unmarshal(file, &conf)
    json.Unmarshal(file, &googleConf)

    // Get the data of courses from json files
    //temp := []string{"course description", "course name", "au", "prereq", "course code", "time", "venue"}
    //result := utils.GetEnum(temp)

    // start server
    course := course.NewCourse()

    r := mux.NewRouter()

    r.HandleFunc("/", defaultHandler)
    r.HandleFunc("/query", handler.NewQueryHandler(course))
    r.HandleFunc("/webhook", handler.WebhookHandler)
    r.HandleFunc("/webhook-v1", handler.NewWebhookHandlerV1(&googleConf))
    r.HandleFunc("/dummy-webhook", handler.DummyWebhookHandler)
    r.HandleFunc("/internal-query", handler.InternalHandler)

    // Apply the CORS middleware to our top-level router, with the defaults.
    if(conf.IsProduction) {
        fmt.Println("Starting a production server on port %d...\n", conf.Port)
        err := http.ListenAndServeTLS(":" + strconv.Itoa(conf.Port), conf.CertFile, conf.KeyFile, cors.Default().Handler(r))
        fmt.Println(err.Error())
    } else {
        fmt.Printf("Starting a development local server on port %d...\n", conf.Port)
        http.ListenAndServe(":" + strconv.Itoa(conf.Port), cors.Default().Handler(r))
    }
    
    
    
}
// todo: fix typo in application security json data
// todo: fix computer security entity
// todo: fix ce1004
