package main

import "crypto/sha512"
import "encoding/base64"
import "flag"
import "fmt"
import "log"
import "net/http"
import "strconv"
import "strings"
import "time"

const(
    SLEEP_SECONDS = 5
)

func main() {
    // Get the port to listen on
    portPtr := flag.Int("port", 8080, "the port to run the listener on")
    flag.Parse()
    fmt.Println("Listener port:",*portPtr)

    http.HandleFunc("/", handleRequest)
    
    log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*portPtr), nil))
}

// Handles incoming requests.
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Remember when we got the request so we can respond after time == SLEEP_SECONDS
    start := time.Now()
    if r.Method != "POST" {
        fmt.Fprintf(w, "You must post to this service")
    } else {
        // Do hashing and encoding right away so we can respond in as close to SLEEP_SECONDS as possible
        password := strings.TrimSpace(r.FormValue("password"))
        hashed := hashAndEncode(password)
        
        // Calculate remaining sleep time after calculations are done
        remainingSleepTime := (time.Duration(SLEEP_SECONDS) * time.Second) - time.Since(start)
        
        // Sleep for that time, if at the time or past this will return immediately
        time.Sleep(remainingSleepTime)
        
        // Return the hashed value
        fmt.Fprintf(w, "%s", hashed)
    }
}

// Hash and encode the input
func hashAndEncode(s string) string {
    sha_512 := sha512.New()
    sha_512.Write([]byte(s))
    return base64.StdEncoding.EncodeToString(sha_512.Sum(nil))
}
