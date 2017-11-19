package main

import "flag"
import "fmt"
import "log"
import "net/http"
import "strconv"
import "strings"

func main() {
    // Get the port to listen on
    portPtr := flag.Int("port", 8080, "the port to run the listener on")
    flag.Parse()
    fmt.Println("Listener port:",*portPtr)
    
    // Basic structure of TCP listener curtesy of https://coderwall.com/p/wohavg/creating-a-simple-tcp-server-in-go
    
    CONN_PORT := strconv.Itoa(*portPtr)
    
    http.HandleFunc("/", handleRequest)
    
    log.Fatal(http.ListenAndServe(":"+CONN_PORT, nil))
}

// Handles incoming requests.
func handleRequest(w http.ResponseWriter, r *http.Request) {

    if r.Method == "GET" {
        fmt.Fprintf(w, "You must post to this service")
    } else {
        value := strings.TrimSpace(r.FormValue("value"))
        
        fmt.Fprintf(w, "Recieved value: " + value)
    }   
}