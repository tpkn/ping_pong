package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

var version = "0.0.0"
var help = fmt.Sprintf(Help, version)

var args = struct {
	Port    string
	Json    bool
	Silent  bool
	Help    bool
	Version bool
	CPU     int
}{}

type Response struct {
	Timestamp int64  `json:"timestamp"`
	Method    string `json:"method"`
	Url       string `json:"url"`
	Body      string `json:"body"`
}

// RequestHandler default handler for all requests
func RequestHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data = Response{
		Timestamp: time.Now().UnixMicro(),
		Method:    r.Method,
		Url:       "http://" + r.Host + r.RequestURI,
		Body:      "",
	}

	if r.Body != nil {
		body, err := io.ReadAll(r.Body)
		if err == nil {
			data.Body = string(body)
		}
	}

	var response_data []byte
	if args.Json {
		w.Header().Set("Content-Type", "application/json")
		response_data, _ = json.Marshal(data)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		response_data = []byte(fmt.Sprintf("%v\t%v\t%v\t%v", data.Timestamp, data.Method, data.Url, data.Body))
	}

	if !args.Silent {
		fmt.Println(string(response_data))
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write(response_data)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	flag.StringVar(&args.Port, "p", "8181", "Server port (default: 8181")
	flag.BoolVar(&args.Json, "j", false, "Pong in JSON format (default: tsv)")
	flag.BoolVar(&args.Silent, "s", false, "Don't print requests to Stdout")
	flag.StringVar(&args.Port, "port", "8181", "Alias for -p")
	flag.BoolVar(&args.Json, "json", false, "Alias for -j")
	flag.BoolVar(&args.Silent, "silent", false, "Alias for -s")
	flag.IntVar(&args.CPU, "cpu", runtime.NumCPU(), "Maximum number of CPU cores used by server")
	flag.BoolVar(&args.Help, "h", false, "Help")
	flag.BoolVar(&args.Help, "help", false, "Alias for -h")
	flag.BoolVar(&args.Version, "v", false, "Version")
	flag.BoolVar(&args.Version, "version", false, "Alias for -v")
	flag.Parse()

	runtime.GOMAXPROCS(args.CPU)

	log.SetFlags(0)
	log.SetPrefix("Error: ")

	if args.Help {
		fmt.Print(help)
		os.Exit(0)
	}

	if args.Version {
		fmt.Print(version)
		os.Exit(0)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", RequestHandler)
	fmt.Println("✓ Server is up and running: http://localhost:" + args.Port)
	err := http.ListenAndServe(":"+args.Port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
