package main

import (
	"fmt"
	"net/http"
	
    "io/ioutil"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "./storage"
"strings"    
"github.com/marcossegovia/apiai-go"    
)

// "github.com/kamalpy/apiai-go"

type query_struct struct {
    Query string
    SessionID string
}

type response struct {
    Response string
}
func handler(w http.ResponseWriter, r *http.Request) {
    // TODO: Handle this properly.
	fmt.Fprintf(w, "Hello, world! Your URL: %s", r.URL.Path[1:])
}


func queryHandler(rw http.ResponseWriter, req *http.Request) {
    db, err := storage.NewDB("test.sqlite3")
    all, _ := db.ListAll()
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
    // ai := apiaigo.APIAI{
	// 	AuthToken: "031636d290f341729417585f09f1ebc4",
	// 	Language:  "en-US",
	// 	SessionID: "32314214",
	// 	Version:   "20150910",
	// }

    if err != nil {
        fmt.Printf("%v", err)
    }

	//Set the query string and your current user identifier.
    // TODO: Set the proper sessionID per user.
    // qr, err := client.Query(apiai.Query{Query: []string{t.Query}, SessionId: "123123", Contexts: []apiai.Context{apiai.Context{Name: t.Context}} })
    qr, err := client.Query(apiai.Query{Query: []string{t.Query}, SessionId: t.SessionID})
    
    // qr, err := client.Query(apiai.Query{Query: []string{"What are the available scholarship for ASEAN students?"}, SessionId: "123123"})
    
    // params, _ := json.Marshal(qr.Result.Params)
    // qr, err := ai.SendText("What are the available scholarship for ASEAN students?")
    fmt.Printf("%v", err)
    // fmt.Printf("%v", qr.Result.Action)
    // qwordValue := qr.Result.Params["QWord.original"].(string)
    qwordValue := "What"
    entityValue := ""
    groupValue := ""
    if(qr.Result.Params["Entity"] != nil) {
        entityValue = qr.Result.Params["Entity"].(string)
    }
    if(qr.Result.Params["Group"] != nil && qr.Result.Params["Group"] != "") {
        groupValue = qr.Result.Params["Group"].(string)
    } else if(entityValue != "") {
        groupValue = "general"
    }
    // entityValue := "Scholarship"
    // groupValue := "ASEAN"
    fmt.Printf("-----")
    fmt.Printf("%v %v %v",qwordValue, entityValue, groupValue)
    fmt.Printf("-----")
    
    var resultMap map[string]string
    resultMap = make(map[string]string)
    
    resultMap["Result"] = qr.Result.Fulfillment.Speech
    resultMap["Context"] = ""

    // Handling context
    // TODO: Handle multiple contexts
    // if(len(qr.Result.Contexts) > 0) {
    //     resultMap["Context"] = qr.Result.Contexts[0].Name
    // }
    // fmt.Printf("Context: %v", resultMap["Context"])

    
    for _, elem := range all {
        if(strings.Compare(entityValue, elem.Entity) == 0 && strings.Compare(groupValue, elem.Query) == 0) { //&& strings.Compare(qwordValue, elem.QWord) == 0) {
            fmt.Printf("Found: %v", elem.Value)
            resultMap["Result"] = elem.Value
        }
    }

    // TODO: Handle this properly
    if(resultMap["Result"] == "") {
        resultMap["Result"] = "One more time?"
    }
    
    resultJson, _ := json.Marshal(resultMap)


    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }

    rw.Header().Set("Content-Type", "application/json")
    
    rw.Write(resultJson)
}

func main() {
    // _, _ = storage.NewDB("test.sqlite3")
    // if(err != nil) {
    //     fmt.Printf("%v\n", err)
    // }
    // all, err := db.ListAll()
    // if(err != nil) {
    //     fmt.Printf("%v\n", err)
    // }
    // for _, elem := range all {
    //     fmt.Printf(elem.Query)
    // }
    // fmt.Printf("%v", all)
    r := mux.NewRouter()
    r.HandleFunc("/", handler)
    r.HandleFunc("/query", queryHandler)

    // Apply the CORS middleware to our top-level router, with the defaults.
    http.ListenAndServe(":8080", cors.Default().Handler(r))
}