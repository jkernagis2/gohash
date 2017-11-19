package main

import "context"
import "crypto/sha512"
import "encoding/base64"
import "flag"
import "fmt"
import "log"
import "net/http"
import "os"
import "strconv"
import "strings"
import "os/signal"
import "syscall"
import "sync"
import "time"

const(
    SLEEP_SECONDS = 5
)

type Server struct {
	logger *log.Logger
	mux    *http.ServeMux
    stop   chan os.Signal
}

func NewServer(options ...func(*Server)) *Server {
	s := &Server{
		logger: log.New(os.Stdout, "", 0),
		mux:    http.NewServeMux(),
        stop:   make(chan os.Signal),
	}

	for _, f := range options {
		f(s)
	}

	s.mux.HandleFunc("/", s.handleRequest)

	return s
}

// Handles incoming requests.
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
    // Remember when we got the request so we can respond after time == SLEEP_SECONDS
    start := time.Now()
    if r.Method != "POST" {
        s.logger.Printf("Non post method attempted")
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprintf(w, "You must post to this service")
    } else {
        if(r.FormValue("graceful shutdown") == "true"){
            s.logger.Printf("graceful shutdown request received")
            
            time.Sleep(time.Duration(SLEEP_SECONDS) * time.Second)
            
            w.WriteHeader(http.StatusOK)
            fmt.Fprintf(w, "Server will shut down momentarily")
            
            s.stop <- syscall.SIGINT
        }else if(r.FormValue("password") != ""){
            s.logger.Printf("Hash request received")
            
            // Do hashing and encoding right away so we can respond in as close to SLEEP_SECONDS as possible
            password := strings.TrimSpace(r.FormValue("password"))
            retVal := hashAndEncode(password)
            
            // Calculate remaining sleep time after calculations are done
            remainingSleepTime := (time.Duration(SLEEP_SECONDS) * time.Second) - time.Since(start)
            
            // Sleep for that time, if at the time or past this will return immediately
            time.Sleep(remainingSleepTime)
            
            // Return
            fmt.Fprintf(w, "%s", retVal)
        }else{
            s.logger.Printf("Invalid request received")
            
            time.Sleep(time.Duration(SLEEP_SECONDS) * time.Second)
            
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(w, "Invalid request")
        }
    }
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func main() {
    // Get the port to listen on
    portPtr := flag.Int("port", 8080, "the port to run the listener on")
    flag.Parse()
    
    logger := log.New(os.Stdout, "", 0)
    
    addr := ":"+strconv.Itoa(*portPtr)
    
    stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT)
    
    s := NewServer(func(s *Server) {
                    s.logger = logger
                    s.stop = stop
                })
    h := &http.Server{Addr: addr, Handler: s}
    
    logger.Printf("Listening on http://localhost%s\n", addr)
    
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        h.ListenAndServe()
    }()
    
    select {
        case signal := <-stop:
            logger.Printf("Got signal:%v\n", signal)
	}
    
    logger.Printf("Stopping listener\n")
    // Start graceful shutdown
    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	h.Shutdown(ctx)
	logger.Printf("Waiting on server\n")
	wg.Wait()
}

// Hash and encode the input
func hashAndEncode(s string) string {
    sha_512 := sha512.New()
    sha_512.Write([]byte(s))
    return base64.StdEncoding.EncodeToString(sha_512.Sum(nil))
}
