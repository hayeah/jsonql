package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/hayeah/jsonql"
	"github.com/hayeah/jsonql/handler"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		help()
	}

	cmd := args[1]
	args = args[2:]

	switch cmd {
	case "server":
		flags := flag.NewFlagSet("server", flag.ExitOnError)

		var port int
		flags.IntVar(&port, "p", 4000, "HTTP server port")
		flags.Parse(args)

		args = flags.Args()

		if len(args) != 1 {
			help()
		}
		dataSource := args[0]
		server(port, dataSource)

	case "get":
		if len(args) != 2 {
			help()
		}
		baseURL := args[0]
		jsonQuery := args[1]
		get(baseURL, jsonQuery)
		// flags := flag.NewFlagSet("server", flag.ExitOnError)
		// req.
	default:
		log.Println("Unrecognized command:", cmd)
		help()
	}

}

func server(port int, dataSource string) {
	db, err := jsonql.OpenDB("sqlite3", dataSource)
	if err != nil {
		log.Fatal(err)
	}

	jsonqlService := handler.NewHTTPHandler(db)

	log.Println("jsonql HTTP started on port", port)
	http.Handle("/", jsonqlService)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func get(baseURL string, jsonQuery string) {
	matched, err := regexp.MatchString("^https?://", baseURL)
	if err != nil {
		log.Fatal(err)
	}

	if !matched {
		baseURL = "http://" + baseURL
	}

	requestURL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	if requestURL.Scheme == "" {
		requestURL.Scheme = "http"
	}

	// add the q query parameter
	urlq := requestURL.Query()
	urlq.Add("q", jsonQuery)
	requestURL.RawQuery = urlq.Encode()

	log.Println("GET:", requestURL.String())
	res, err := http.Get(requestURL.String())
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(os.Stdout, res.Body)
}

func help() {
	intro := `
Start a HTTP jsonql service:
  jsonql server [-p 4000] <sqlite.db>

Make a jsonql request:
  jsonql get <host> <jsonql>
  `
	fmt.Println(intro)
	os.Exit(1)
	// flag.PrintDefaults()
}
