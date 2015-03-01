package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/hayeah/jsonql"
	"github.com/hayeah/jsonql/handler"

	kjson "github.com/klauspost/json"

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

	case "q", "query":
		if len(args) < 1 {
			help()
		}
		baseURL := args[0]

		var err error
		var input io.Reader
		if len(args) == 2 {
			input, err = os.Open(args[2])
			if err != nil {
				log.Fatal(err)
			}
		} else {
			input = os.Stdin
		}

		query(baseURL, input)
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

func query(baseURL string, r io.Reader) {
	matched, err := regexp.MatchString("^https?://", baseURL)
	if err != nil {
		log.Fatal(err)
	}

	input, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	jsonQuery := string(input)

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

	sbuf := bytes.NewBufferString("")
	io.Copy(sbuf, res.Body)
	// kjson.IndentStream swallows the last char, because at the end res.Body
	// responds with the last char and EOF, and IndentStream immediately
	// returns without processing the last char.
	err = kjson.IndentStream(os.Stdout, sbuf, "", "  ")
	if err != nil {
		log.Println(err)
	}

	// io.Copy(os.Stdout, res.Body)
}

func help() {
	intro := `
Start a HTTP jsonql service:
  jsonql server [-p 4000] <sqlite.db>

Make a jsonql request:
  jsonql [q|query] <host> <jsonql.json>
  cat jsonql.json | jsonql query
  `
	fmt.Println(intro)
	os.Exit(1)
	// flag.PrintDefaults()
}
