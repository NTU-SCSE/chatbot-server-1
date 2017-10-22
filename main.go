package main

import (
    "fmt"
    "sort"
    "reflect"
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
    
    if err != nil {
        fmt.Printf("%v", err)
    }

	//Set the query string and your current user identifier.

    qr, err := client.Query(apiai.Query{Query: []string{t.Query}, SessionId: t.SessionID})
    qwordValue := "What"
    entityValue := ""
    groupValue := make([]string, 0)
    
    if(qr.Result.Params["Entity"] != nil) {
        entityValue = qr.Result.Params["Entity"].(string)
    } else if(qr.Result.Metadata.IntentName == "") {
        entityValue = qr.Result.Metadata.IntentName
    }
    
    if(qr.Result.Params["Group"] != nil && len(qr.Result.Params["Group"].([]string)) > 0) {
        // TODO: handle multiple group values
        groupValue = append(groupValue, qr.Result.Params["Group"].([]string)...)
        sort.Strings(groupValue)
    } else if(entityValue != "") {
        groupValue = append(groupValue, "general")
    }
    
    var resultMap map[string]string
    resultMap = make(map[string]string)
    
    resultMap["Result"] = qr.Result.Fulfillment.Speech
    resultMap["Context"] = ""

    // TODO: Handle multiple intents.
    if(qr.Result.Contexts != nil && len(qr.Result.Contexts) > 0) {
        resultMap["Context"] = qr.Result.Contexts[0].Name
    }
    
    fmt.Printf("-----")
    fmt.Printf("%v %v %v",qwordValue, entityValue, groupValue)
    fmt.Printf("-----")
    
    

    // Handling context
    // TODO: Handle multiple contexts
    // if(len(qr.Result.Contexts) > 0) {
    //     resultMap["Context"] = qr.Result.Contexts[0].Name
    // }
    // fmt.Printf("Context: %v", resultMap["Context"])

    
    for _, elem := range all {
        dbValue := strings.Split(elem.Query, ",")
        sort.Strings(dbValue)
        if(strings.Compare(entityValue, elem.Entity) == 0 && reflect.DeepEqual(groupValue, dbValue)) { //&& strings.Compare(qwordValue, elem.QWord) == 0) {
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
    r := mux.NewRouter()
    r.HandleFunc("/", handler)
    r.HandleFunc("/query", queryHandler)

    // Apply the CORS middleware to our top-level router, with the defaults.
    http.ListenAndServe(":8080", cors.Default().Handler(r))
}