package main

import (
    "fmt"
    "sort"
    "reflect"
	"net/http"
    "./utils"
    "io/ioutil"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "./storage"
    "strconv"
    "strings"    
    "os"
    "github.com/marcossegovia/apiai-go"    
    "github.com/tidwall/gjson"
    "regexp"
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

type course struct {
    Code string `json:"code"`
    Name string `json:"name"`
    AU int `json:"AU"`
    PreReq string `json:"preReq"`
    Description string `json:"description"`
}

type class struct {
    Code string `json:"code"`
    Index string `json:"index"`
    Type string `json:"type"`
    Group string `json:"group"`
    Day string `json:"day"`
    Time string `json:"time"`
    Venue string `json:"venue"`
    Remark string `json:"remark"`
}

var courses []course
var classes []class

func parseCourseCode(param string) string {
    var validCode = regexp.MustCompile(`^c[z|e][0-9]{4}$`)
    return validCode.FindString(param)
}

func getCourseCode(param string) (string, string) {
    param = strings.TrimSpace(param)
    
    // TODO: if query is "CZ#### code", it won't return CE/CZ#### format, fix this
    // TODO: api.ai still not trained yet to recognize CE/CZ#### format
    if _, err := strconv.Atoi(param[len(param)-4:]); err == nil {
        return param, ""
    }
    result := ""
    auxResult := ""

    for _, course := range courses {
        if strings.ToLower(strings.TrimSpace(course.Name)) == strings.ToLower(param) {
            if(result == "") {
                result = course.Code
            } else if course.Code[2] != '/' {
                auxResult = "CE/CZ" + course.Code[2:]
            }
        }
    }
    return result, auxResult
}

func getSchedulePrint(param class) string {
    var result string
    if(param.Type == "LEC/STUDIO") {
        result = "Lecture"
    } else if(param.Type == "TUT") {
        result = "Tutorial"
    } else {
        result = param.Type
    }
    result = result + " on " + param.Day + ", " + param.Time + " at " + param.Venue
    return result
}

func getIndex(code string) map[string]bool {
    result := map[string]bool{}
    for _, class := range classes {
        if strings.ToLower(strings.TrimSpace(class.Code)) == strings.ToLower(code) {
            result[class.Index] = true
        }
    }
    return result
}

func getIndexString(code string) string {
    indexList := getIndex(code)
    indexStrings := []string{}
    for key, _ := range indexList {
        indexStrings = append(indexStrings, key)
    }
    sort.Strings(indexStrings)
    result := ""
    for _, key := range indexStrings {
        result = result + key + "\n"
    }
    return result
}

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

func webhookHandler(rw http.ResponseWriter, req *http.Request) {
    db, err := storage.NewDB("test.sqlite3")
    all, _ := db.ListAll()
    body, err := ioutil.ReadAll(req.Body)

    if(err != nil) {
        panic(err)
    }
    //fmt.Println(string(body[:]))
    fullJSON := string(body[:])

    // Get query parameters, sort them
    paramsJSON := gjson.Get(fullJSON, "result.parameters")
    params := make([]string, 0)
    paramsJSON.ForEach(func(key, value gjson.Result) bool {
        for _, elem := range value.Array() {
            if(elem.String() != "") {
                params = append(params, strings.ToLower(elem.String()))
            }
        }
        return true
    })
    sort.Strings(params)

    //originalRequest := gjson.Get(string(body[:]), "originalRequest.data.message.text")

    // get original request text, intent and contexts
    originalRequest := gjson.Get(fullJSON, "result.resolvedQuery")
    intent := gjson.Get(fullJSON, "result.metadata.intentName")
    contextsJSON := gjson.Get(fullJSON, "result.contexts")
    contexts := make([]string, 0)
    for _, elem := range contextsJSON.Array() {
        contexts = append(contexts, gjson.Get(elem.String(), "name").String())
    }

    // preparing response
    var resultMap map[string]interface{}
    resultMap = make(map[string]interface{})

    // for debugging
    fmt.Println(originalRequest)
    fmt.Println(params)
    fmt.Println(intent)

    // TODO: fill with proper values here
    resultMap["displayText"] = "Test Response"
    resultMap["speech"] = "Response not found"
    resultMap["data"] = ""
    resultMap["contextOut"] = []string{}
    resultMap["source"] = "Hello"

    if(strings.Compare(intent.String(), "Course") == 0) {
        // course related
        for _, elem := range params {
            courseCode := parseCourseCode(elem)
            if(courseCode != "") {
                course, _ := db.GetCourseByCode(courseCode)
                resultMap["speech"] = course.Description
            }
        }
    } else {
        // other queries
        for _, elem := range all {
            dbValue := strings.Split(elem.Query, ",")    
            sort.Strings(dbValue)
            for index, _ := range dbValue {
                dbValue[index] = strings.ToLower(dbValue[index])
            }
            if(strings.Compare(intent.String(), elem.Intent) == 0 && reflect.DeepEqual(params, dbValue)) { //&& strings.Compare(qwordValue, elem.QWord) == 0) {
                resultMap["speech"] = elem.Value
            }
        }
    }

    // default fallback: direct to google search, get the first result
    if strings.Compare(resultMap["speech"].(string), "Response not found") == 0 {
        resp, err := http.Get("https://www.googleapis.com/customsearch/v1?q=" + 
            strings.Replace(originalRequest.String(), " ", "+", -1) + "&cx=000348109821987500770%3Ar1ufthpxqxg&key=AIzaSyDW0l64m7xweAo28Z_q3yAskU_d5fbevGw")
        if err != nil {
            // handle error
        }
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        results := gjson.Get(string(body[:]), "items")
        for _, elem := range results.Array() {
            link := gjson.Get(elem.String(), "link").String()
            resultMap["speech"] = "You can find out more about it at " + link + "\r\n"
            break
        }
    }

    resultJson, _ := json.Marshal(resultMap)
    
    rw.Header().Set("Content-Type", "application/json")
        
    rw.Write(resultJson)
}

func queryHandler(rw http.ResponseWriter, req *http.Request) {
    // db, err := storage.NewDB("test.sqlite3")
    // all, _ := db.ListAll()
    body, err := ioutil.ReadAll(req.Body)

    var resultMap map[string]interface{}
    resultMap = make(map[string]interface{})
    
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

    // log into file
    f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY, 0644)
    if(err != nil) {
        fmt.Printf("error %v", err)
    }
    defer f.Close()
    f.WriteString("Query from: " + t.SessionID + "\r\n")
    f.WriteString(t.Query + "\r\n")

    //Set the query string and your current user identifier.
    var qr *apiai.QueryResponse
    if(t.Query == "reset") {
        qr, err = client.Query(apiai.Query{Query: []string{t.Query}, SessionId: t.SessionID, ResetContexts: true})
        resultMap["Result"] = "reset"
        resultJson, _ := json.Marshal(resultMap)
    
        rw.Header().Set("Content-Type", "application/json")
        rw.Write(resultJson)
        f.WriteString("----------\r\n")
        return
    } else {
        if ind, err := strconv.Atoi(t.Query); err == nil && len(t.Enum) > 0 && ind > 0 && ind <= len(t.Enum) {
            qr, err = client.Query(apiai.Query{Query: []string{t.Enum[ind - 1]}, SessionId: t.SessionID})
        } else {
            qr, err = client.Query(apiai.Query{Query: []string{t.Query}, SessionId: t.SessionID})
        }        
    }
    //qwordValue := "What"
    entityValue := ""
    intentValue := ""
    number := "" // TODO: should we use integer instead?
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
    if(qr.Result.Params["number"] != nil) {
        number = qr.Result.Params["number"].(string)
    }
    

    possibleParams := []string{"Course", "Group", "Course1"}
    // for key, _:= range qr.Result.Params {
    //     fmt.Printf("%v\n", key)
    //     paramValue = key
    // }
    for _, paramValue := range possibleParams {
        if(qr.Result.Params[paramValue] != nil && len(qr.Result.Params[paramValue].([]interface{})) > 0) {
            // TODO: handle multiple group values
            fmt.Printf("Param %v %v\n", paramValue, qr.Result.Params[paramValue])
            for _, v := range qr.Result.Params[paramValue].([]interface{}) {
                groupValue = append(groupValue, v.(string))
            }
            
            sort.Strings(groupValue)
        }
    }
    if(len(groupValue) == 0 && entityValue != "") {
        groupValue = append(groupValue, "general")
    }
    
    
    
    resultMap["Result"] = ""
    resultMap["Context"] = ""
    resultMap["Enum"] = ""
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
    // for _, elem := range all {
    //     dbValue := strings.Split(elem.Query, ",")    
    //     sort.Strings(dbValue)
    //     if(strings.Compare(entityValue, elem.Entity) == 0 && reflect.DeepEqual(groupValue, dbValue)) { //&& strings.Compare(qwordValue, elem.QWord) == 0) {
    //         resultMap["Result"] = elem.Value
    //     }
    // }

    courseIntent := []string{"course description", "course name", "au", "prereq", "course code", "time", "venue"}
    courseCode := ""
    auxCourseCode := ""
    courseAttr := ""
    // matching with json
    if(resultMap["Result"] == "" && len(intentValue) >= 6 && intentValue[:6] == "Course") {    
        for _, param := range groupValue {
            if !contains(courseIntent, param) {
                courseCode, auxCourseCode = getCourseCode(param)
            } else {
                courseAttr = param
            }
        }
        
        if(courseAttr == "venue" || courseAttr == "time") {
            for _, class := range classes {
                if strings.ToLower(class.Code) == strings.ToLower(courseCode) &&
                class.Index == number {
                    resultMap["Result"] = resultMap["Result"].(string) + getSchedulePrint(class) + "\n"
                }
            }
        } else {
            // TODO: Fixed white space trimming properly
            for _, course := range courses {
                if strings.ToLower(course.Code) == strings.ToLower(courseCode) {
                    if(courseAttr == "course description") {
                        resultMap["Result"] = course.Description
                    } else if(courseAttr == "course name") {
                        resultMap["Result"] = course.Name
                    } else if(courseAttr == "au") {
                        resultMap["Result"] = strconv.Itoa(course.AU)
                    } else if(courseAttr == "prereq") {
                        resultMap["Result"] = course.PreReq
                    } else if(courseAttr == "course code") {
                        if(auxCourseCode == "") {
                            resultMap["Result"] = courseCode
                        } else {
                            resultMap["Result"] = auxCourseCode
                        }
                    }
                }
            }
        }
    }

    // nothing matched
    // TODO: Handle this properly
    fmt.Printf("%v %v %v %v\n", courseCode, intentValue, courseAttr, number)
    if(resultMap["Result"] == "") {
        if intentValue[:4] == "SCSE" {
            resultMap["Result"] = "Find out more about SCSE courses by specifying the course code or course name.\n"
            qr, err = client.Query(apiai.Query{Query: []string{t.Query}, SessionId: t.SessionID, ResetContexts: true})
        } else if intentValue[:6] == "Course" && courseCode == "" {
            resultMap["Result"] = "Please specify the course code or course name."
            qr, err = client.Query(apiai.Query{Query: []string{t.Query}, SessionId: t.SessionID, ResetContexts: true})
        } else if intentValue[:6] == "Course" && courseAttr != "venue" && courseAttr != "time" {
            temp := []string{"Course Name", "Academic Units", "Course description", "Class schedule", "Venues"}
            str := utils.GetEnum(temp)
            resultMap["Result"] = "What do you want to know about " + courseCode + "?"+ str
            resultMap["Enum"] = temp
        } else if intentValue[:6] == "Course" {
            resultMap["Result"] = "Please specify your index number.\n" + getIndexString(courseCode)
        } else if intentValue[:6] == "Hostel" {
            temp := []string{"Application", "Criteria", "Fee"}
            str := utils.GetEnum(temp)
            resultMap["Result"] = "What do you want to know about NTU Hostel Accomodation?" + str
            resultMap["Enum"] = temp
        } else {
            resultMap["Result"] = "One more time?"
        }
    }
    // or we can use:
    // qr.Result.Fulfillment.Speech
    

    // log to file
    f.WriteString("Response:\r\n")
    f.WriteString(strings.Replace(resultMap["Result"].(string), "\n", "\r\n", -1) + "\r\n")
    f.WriteString("----------\r\n")

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
    //temp := []string{"course description", "course name", "au", "prereq", "course code", "time", "venue"}
    //result := utils.GetEnum(temp)
    
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

    // Get the data of course schedules and venues
    file, err = ioutil.ReadFile("./schedules.json")
    if err != nil {
        fmt.Println(err.Error())
    }
    json.Unmarshal(file, &classes)

    // start server
    fmt.Println("Server started...")
    r := mux.NewRouter()
    r.HandleFunc("/", handler)
    r.HandleFunc("/query", queryHandler)
    r.HandleFunc("/webhook", webhookHandler)

    // Apply the CORS middleware to our top-level router, with the defaults.
    http.ListenAndServe(":8080", cors.Default().Handler(r))
}
// todo: fix typo in application security json data
// todo: fix computer security entity
// todo: fix ce1004