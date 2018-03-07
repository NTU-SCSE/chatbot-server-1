package main

import (
    "fmt"
	"net/http"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "./handler"    
    "./course"
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
    // Get the data of courses from json files
    //temp := []string{"course description", "course name", "au", "prereq", "course code", "time", "venue"}
    //result := utils.GetEnum(temp)

    // start server
    course := course.NewCourse()

    fmt.Println("Server started...")
    r := mux.NewRouter()

    // set up the registered routes
    // for _, rt := range(routes) {
    //     fmt.Println("Registering...")
    //     fmt.Println(rt.path)
    //     r.HandleFunc(rt.path, rt.handler).Methods(rt.method)
    // }

    r.HandleFunc("/", defaultHandler)
    r.HandleFunc("/query", handler.NewQueryHandler(course))
    r.HandleFunc("/webhook", handler.WebhookHandler)
    r.HandleFunc("/webhook-v1", handler.WebhookHandlerV1)
    r.HandleFunc("/dummy-webhook", handler.DummyWebhookHandler)
    r.HandleFunc("/internal-query", handler.InternalHandler)

    // Apply the CORS middleware to our top-level router, with the defaults.
    // TODO: move this to const
    IS_PRODUCTION := false
    if(IS_PRODUCTION) {
        err := http.ListenAndServeTLS(":8080", "/etc/letsencrypt/live/www.pieceofcode.org/fullchain.pem", "/etc/letsencrypt/live/www.pieceofcode.org/privkey.pem", cors.Default().Handler(r))
        fmt.Println(err.Error())
    } else {
        http.ListenAndServe(":8080", cors.Default().Handler(r))
    }
    
    
    
}
// todo: fix typo in application security json data
// todo: fix computer security entity
// todo: fix ce1004
