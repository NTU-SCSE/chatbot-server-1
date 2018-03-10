package handler
import (
    "net/http"
    "io/ioutil"
    "fmt"
    "encoding/json"
    "github.com/marcossegovia/apiai-go"    
    "../config"
)

type query_struct struct {
    Query string
    SessionID string
    Enum []string   `json:",omitempty"`
}

func NewInternalHandler(conf *config.AgentConfig) func(http.ResponseWriter, *http.Request) {
    return func(rw http.ResponseWriter, req *http.Request) {
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
                Token:      conf.Token,
                QueryLang:  conf.QueryLang,
                SpeechLang: conf.SpeechLang,
            },
        )
        
        if err != nil {
            fmt.Printf("%v", err)
        }

        // log into file
        // f, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
        // if(err != nil) {
        //     fmt.Printf("error %v", err)
        // }
        // defer f.Close()

        // f.WriteString("Query from: " + t.SessionID + "\r\n")
        // f.WriteString(t.Query + "\r\n")

        //Set the query string and your current user identifier.
        var qr *apiai.QueryResponse
        if(t.Query == "reset") {
            qr, err = client.Query(apiai.Query{Query: []string{t.Query}, SessionId: t.SessionID, ResetContexts: true})
            resultMap["Result"] = "reset"
            resultJson, _ := json.Marshal(resultMap)
        
            rw.Header().Set("Content-Type", "application/json")
            rw.Write(resultJson)
            // f.WriteString("----------\r\n")
            return
        } else {
            // if ind, err := strconv.Atoi(t.Query); err == nil && len(t.Enum) > 0 && ind > 0 && ind <= len(t.Enum) {
            //     qr, err = client.Query(apiai.Query{Query: []string{t.Enum[ind - 1]}, SessionId: t.SessionID})
            // } else {
            qr, err = client.Query(apiai.Query{Query: []string{t.Query}, SessionId: t.SessionID})
            // }        
        }
        fmt.Println(qr.Result.Fulfillment.Speech)
        resultMap["Result"] = qr.Result.Fulfillment.Speech

        resultJson, _ := json.Marshal(resultMap)

        rw.Header().Set("Content-Type", "application/json")
        
        rw.Write(resultJson)
    }
}