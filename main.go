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
    "strconv"
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

type course struct {
    Code string `json:"code"`
    Name string `json:"name"`
    AU int `json:"AU"`
    PreReq string `json:"preReq"`
    Description string `json:"description"`
}

var courses []course


func contains(slice []string, item string) bool {
    set := make(map[string]struct{}, len(slice))
    for _, s := range slice {
        set[s] = struct{}{}
    }

    _, ok := set[item] 
    return ok
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
    //qwordValue := "What"
    entityValue := ""
    intentValue := ""
    groupValue := make([]string, 0)
    
    if(qr.Result.Metadata.IntentName != "") {
        intentValue = qr.Result.Metadata.IntentName
    }

    if(qr.Result.Params["Entity"] != nil) {
        entityValue = qr.Result.Params["Entity"].(string)
    }
    if(entityValue == "") {
        entityValue = intentValue
    }
    
    paramValue := ""
    for key, _:= range qr.Result.Params {
        fmt.Printf("%v\n", key)
        paramValue = key
    }
    if(qr.Result.Params[paramValue] != nil && len(qr.Result.Params[paramValue].([]interface{})) > 0) {
        // TODO: handle multiple group values
        for _, v := range qr.Result.Params[paramValue].([]interface{}) {
            groupValue = append(groupValue, v.(string))
        }
        
        sort.Strings(groupValue)
    } else if(entityValue != "") {
        groupValue = append(groupValue, "general")
    }
    
    var resultMap map[string]string
    resultMap = make(map[string]string)
    
    resultMap["Result"] = ""
    resultMap["Context"] = ""

    // TODO: Handle multiple intents.
    if(qr.Result.Contexts != nil && len(qr.Result.Contexts) > 0) {
        resultMap["Context"] = qr.Result.Contexts[0].Name
    }

    // Handling context
    // TODO: Handle multiple contexts
    // if(len(qr.Result.Contexts) > 0) {
    //     resultMap["Context"] = qr.Result.Contexts[0].Name
    // }
    // fmt.Printf("Context: %v", resultMap["Context"])

    // matching with Database
    for _, elem := range all {
        dbValue := strings.Split(elem.Query, ",")    
        sort.Strings(dbValue)
        if(strings.Compare(entityValue, elem.Entity) == 0 && reflect.DeepEqual(groupValue, dbValue)) { //&& strings.Compare(qwordValue, elem.QWord) == 0) {
            resultMap["Result"] = elem.Value
        }
    }

    // matching with json
    if(resultMap["Result"] == "") {
        courseIntent := []string{"course description", "course name", "au", "prereq", "course code"}
        courseCode := ""
        courseName := ""
        courseAttr := ""
        for _, param := range groupValue {
            if !contains(courseIntent, param) {
                courseCode = param
                courseName = param
            } else {
                courseAttr = param
            }
        }
        
        // TODO: Fixed white space trimming properly
        for _, course := range courses {
            if strings.ToLower(course.Code) == strings.ToLower(courseCode) || 
            strings.ToLower(strings.TrimSpace(course.Name)) == strings.ToLower(courseName) {
                if(courseAttr == "course description") {
                    resultMap["Result"] = course.Description
                } else if(courseAttr == "course name") {
                    resultMap["Result"] = course.Name
                } else if(courseAttr == "au") {
                    resultMap["Result"] = strconv.Itoa(course.AU)
                } else if(courseAttr == "prereq") {
                    resultMap["Result"] = course.PreReq
                } else if(courseAttr == "course code") {
                    if(resultMap["Result"] == "") {
                        resultMap["Result"] = course.Code
                    } else if course.Code[2] != '/' {
                        resultMap["Result"] = "CE/CZ" + course.Code[2:]
                    }
                    
                }
                
            }
        }
    }

    // nothing matched
    // TODO: Handle this properly
    if(resultMap["Result"] == "") {
        resultMap["Result"] = "One more time?"
    }
    // or we can use:
    // qr.Result.Fulfillment.Speech

    resultJson, _ := json.Marshal(resultMap)


    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
        return
    }

    rw.Header().Set("Content-Type", "application/json")
    
    rw.Write(resultJson)
}

func main() {
    // Get the data of courses from json files
    file, err := ioutil.ReadFile("./cs.json")
    if err != nil {
        fmt.Println(err.Error())
    }
    json.Unmarshal(file, &courses)
    var CECourses []course
    file, err = ioutil.ReadFile("./ce.json")
    if err != nil {
        fmt.Println(err.Error())
    }
    json.Unmarshal(file, &CECourses)
    courses = append(courses, CECourses...)
    // start server

    r := mux.NewRouter()
    r.HandleFunc("/", handler)
    r.HandleFunc("/query", queryHandler)

    // Apply the CORS middleware to our top-level router, with the defaults.
    http.ListenAndServe(":8080", cors.Default().Handler(r))
}
// todo: fix typo in application security json data
// todo: fix computer security entity
// todo: fix ce1004