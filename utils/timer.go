package utils
import (
	"time"
	"log"
	"os"
	"fmt"
)
func TimeFunction(start time.Time, name string) {
    elapsed := time.Since(start)
    log.Printf("%s %s", name, elapsed)
    f, err := os.OpenFile("log-alpha2.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if(err != nil) {
        fmt.Printf("error %v", err)
    }
    defer f.Close()
    f.WriteString("time: " + elapsed.String() + "\r\n")
    f.WriteString("----------\r\n")
    // log.Printf("%s", elapsed)
}