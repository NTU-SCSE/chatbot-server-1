package handler
import (
    "net/http"
    "io/ioutil"
    "fmt"
    "encoding/json"
    "../utils"
    "../course"
    "time"
    "../storage"
    "github.com/tidwall/gjson"
    "strconv"
    "sort"
    "strings"
    "os"
    "../config"
    "bytes"
)

func NewWebhookHandlerV1(conf *config.GoogleSearchConfig) func(http.ResponseWriter,*http.Request) {
    return func(rw http.ResponseWriter, req *http.Request) {
        defer utils.TimeFunction(time.Now(), "w")
        db, err := storage.NewDB("test.sqlite3")
        body, err := ioutil.ReadAll(req.Body)

        if(err != nil) {
            panic(err)
        }
        //fmt.Println(string(body[:]))
        fullJSON := string(body[:])
        
        // Get session ID
        sessionID := gjson.Get(fullJSON, "sessionId")

        // Check the source
        spellCheckApplied := true
        if (strings.Compare(sessionID.String(), "11111111-1111-1111-1111-111111111111") == 0) {
            // second query with spell checking applied
            spellCheckApplied = false
        }
        // fmt.Println(spellCheckApplied)
        // querySource := gjson.Get(fullJSON, "originalRequest.source")
        // if(querySource.String() == "") {
        //     // second query with spell checking applied
        //     fromFacebook = false
        // }

        // Get query parameters, sort them
        paramsJSON := gjson.Get(fullJSON, "result.parameters")
        params := make([]string, 0)
        var number string
        hasNumber := false
        paramsJSON.ForEach(func(key, value gjson.Result) bool {
            for _, elem := range value.Array() {
                if(elem.String() != "") {
                    // TODO: check DialogFlow system number entity instead
                    if _, err := strconv.Atoi(elem.String()); err == nil {
                        number = elem.String()
                        hasNumber = true
                    } else {
                        params = append(params, strings.ToLower(elem.String()))
                    }
                }
            }
            return true
        })
        sort.Strings(params)
        if hasNumber {
            params[0] = params[0] + number
        }

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
        // fmt.Println(originalRequest)
        // fmt.Println(params)
        // fmt.Println(intent)

        // TODO: fill with proper values here
        resultMap["displayText"] = "Test Response"
        resultMap["speech"] = "Response not found"
        resultMap["data"] = ""
        resultMap["contextOut"] = []string{}
        resultMap["source"] = "Hello"

        // file logging
        f, err := os.OpenFile("log-alpha2.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
        if(err != nil) {
            fmt.Printf("error %v", err)
        }
        defer f.Close()
        f.WriteString("Query from: " + sessionID.String() + "\r\n")
        f.WriteString(originalRequest.String() + "\r\n")
        f.WriteString("----------\r\n")
        f.WriteString("Intent: \r\n" + intent.String() + "\r\n")
        f.WriteString("----------\r\n")

        if(strings.Compare(intent.String(), "Course") == 0) {
            // course related
            for _, elem := range params {
                courseCode := course.ParseCourseCode(elem)
                if(courseCode != "") {
                    course, _ := db.GetCourseByCode(courseCode)
                    for _, field := range params {
                        if(field == "course code") {
                            resultMap["speech"] = course.Code
                        } else if(field == "course name") {
                            resultMap["speech"] = course.Name
                        } else if(field == "au") {
                            resultMap["speech"] = course.AU
                        } else if(field == "course description") {
                            resultMap["speech"] = course.Description
                        } else if(field == "prereq") {
                            resultMap["speech"] = course.PreReq
                        }
                    }
                }
            }
        } else if strings.Compare(intent.String(), "location") == 0 {
            // location queries
            resultMap["speech"] = "Please refer to http://maps.ntu.edu.sg/maps#q:" +
            strings.Replace(params[0], " ", "%20", -1) + "\r\n"
        } else {
            // other queries
            all, _ := db.ListRecordsByIntent(intent.String())
            maxMatchParams := 0
            for _, elem := range all {
                dbValue := strings.Split(elem.Params, ",")
                for index, _ := range dbValue {
                    dbValue[index] = strings.TrimSpace(dbValue[index])
                    dbValue[index] = strings.ToLower(dbValue[index])
                }
                sort.Strings(dbValue)
                currentMatchParams := 0
                for _, param := range dbValue {
                    if utils.Contains(params, param) {
                        currentMatchParams += 1
                    }
                    // default response, if any
                    if maxMatchParams == 0 && param == "default" {
                        resultMap["speech"] = elem.Response
                    }
                }
                if(currentMatchParams > maxMatchParams) {
                    maxMatchParams = currentMatchParams
                    resultMap["speech"] = elem.Response
                }
            }
        }

        // default fallback: direct to google search, get the first result
        if strings.Compare(resultMap["speech"].(string), "Response not found") == 0 {
            if(spellCheckApplied) {
                resp, err := http.Get("https://www.googleapis.com/customsearch/v1?q=" + 
                    "ntu+singapore+" + strings.Replace(originalRequest.String(), " ", "+", -1) + "&cx=" + conf.SearchEngineID + "&key=" + conf.ApiKey)
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
            } else {
                // TODO: move to config
                url := "https://www.pieceofcode.org:8080/spellcheck"

                var jsonStr = []byte(`{"Query":"`+ originalRequest.String() +`"}`)
                req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
                req.Header.Set("Content-Type", "application/json")

                client := &http.Client{}
                resp, err := client.Do(req)
                if err != nil {
                    panic(err)
                }
                defer resp.Body.Close()

                body, _ := ioutil.ReadAll(resp.Body)
                resultMap["speech"] = gjson.Get(string(body[:]), "Result")
            }
        }

        f.WriteString("Response:\r\n")
        f.WriteString(strings.Replace(resultMap["speech"].(string), "\n", "\r\n", -1) + "\r\n")
        f.WriteString("----------\r\n")

        resultJson, _ := json.Marshal(resultMap)
        
        rw.Header().Set("Content-Type", "application/json")
            
        rw.Write(resultJson)
    }
}