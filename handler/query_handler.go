package handler
import (
    "net/http"
    "io/ioutil"
    "encoding/json"
	"../utils"
	"../course"
	"strings"
	"sort"
	"fmt"
	"os"
	"github.com/marcossegovia/apiai-go"    
	"strconv"
)

func NewQueryHandler(c *course.Course) func(http.ResponseWriter, *http.Request) {
    return func (rw http.ResponseWriter, req *http.Request) {
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
        f, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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
                // fmt.Printf("Param %v %v\n", paramValue, qr.Result.Params[paramValue])
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
                if !utils.Contains(courseIntent, param) {
                    courseCode, auxCourseCode = c.GetCourseCode(param)
                } else {
                    courseAttr = param
                }
            }
            
            if(courseAttr == "venue" || courseAttr == "time") {
                for _, class := range c.Classes {
                    if strings.ToLower(class.Code) == strings.ToLower(courseCode) &&
                    class.Index == number {
                        resultMap["Result"] = resultMap["Result"].(string) + course.GetSchedulePrint(class) + "\n"
                    }
                }
            } else {
                // TODO: Fixed white space trimming properly
                for _, mod := range c.Modules {
                    if strings.ToLower(mod.Code) == strings.ToLower(courseCode) {
                        if(courseAttr == "course description") {
                            resultMap["Result"] = mod.Description
                        } else if(courseAttr == "course name") {
                            resultMap["Result"] = mod.Name
                        } else if(courseAttr == "au") {
                            resultMap["Result"] = strconv.Itoa(mod.AU)
                        } else if(courseAttr == "prereq") {
                            resultMap["Result"] = mod.PreReq
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
                resultMap["Result"] = "Please specify your index number.\n" + c.GetIndexString(courseCode)
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
}