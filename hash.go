package main

import "flag"
import "fmt"

func main() {
    portPtr := flag.Int("port", 8080, "the port to run the listener on")
    
    flag.Parse()
    
    fmt.Println("listener port:",*portPtr)
}