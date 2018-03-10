package handler
import (
    "net/http"
    "io/ioutil"
    "encoding/json"
	"../utils"
	"../course"
	"time"
	"../storage"
	"github.com/tidwall/gjson"
	"strings"
	"sort"
	"reflect"
)

// TODO: Please remove if obsolete
func WebhookHandler(rw http.ResponseWriter, req *http.Request) {
    defer utils.TimeFunction(time.Now(), "w")
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
    // fmt.Println(originalRequest)
    // fmt.Println(params)
    // fmt.Println(intent)

    // The following fields are not used for now
    // resultMap["displayText"] = "Test Response"
    // resultMap["data"] = ""
    resultMap["speech"] = "Response not found"
    resultMap["contextOut"] = []string{}
    resultMap["source"] = "golang_server"

    if(strings.Compare(intent.String(), "Course") == 0) {
        // course related
        for _, elem := range params {
            courseCode := course.ParseCourseCode(elem)
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
            "ntu+singapore+" + strings.Replace(originalRequest.String(), " ", "+", -1) + "&cx=000348109821987500770%3Ar1ufthpxqxg&key=AIzaSyDW0l64m7xweAo28Z_q3yAskU_d5fbevGw")
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