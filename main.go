package main

import (
	"fmt"
	"net/http"
	"github.com/marcossegovia/apiai-go"
    "io/ioutil"
    "log"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

type query_struct struct {
    Query string
}

type response struct {
    Response string
}
func handler(w http.ResponseWriter, r *http.Request) {
    // TODO: Handle this properly.
	fmt.Fprintf(w, "Hello, world! Your URL: %s", r.URL.Path[1:])
}


func queryHandler(rw http.ResponseWriter, req *http.Request) {
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        // TODO: Don't use panic, handle properly.
        panic(err)
    }

    var t query_struct
    err = json.Unmarshal(body, &t)
    if err != nil {
        panic(err)
    }

    client, err := apiai.NewClient(
        &apiai.ClientConfig{
            Token:      "031636d290f341729417585f09f1ebc4",
            QueryLang:  "en",    //Default en
            SpeechLang: "en-US", //Default en-US
        },
    )
    
    if err != nil {
        fmt.Printf("%v", err)
    }

	//Set the query string and your current user identifier.
    // TODO: Set the proper sessionID per user.
    qr, err := client.Query(apiai.Query{Query: []string{t.Query}, SessionId: "123123"})
    profile := response{qr.Result.Fulfillment.Speech}

    js, err := json.Marshal(profile)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }

    rw.Header().Set("Content-Type", "application/json")
    
    rw.Write(js)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", handler)
    r.HandleFunc("/query", queryHandler)

    // Apply the CORS middleware to our top-level router, with the defaults.
    http.ListenAndServe(":8080", cors.Default().Handler(r))
}