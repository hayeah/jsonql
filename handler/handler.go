package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/hayeah/jsonql"
)

type HTTPHandler struct {
	db *jsonql.DB
}

func NewHTTPHandler(db *jsonql.DB) *HTTPHandler {
	return &HTTPHandler{db: db}
}

// TODO 404 for any path other than root
// TODO error handling: http://nesv.blogspot.com/2012/09/super-easy-json-http-responses-in-go.html
func (h *HTTPHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	jsonqlString := req.URL.Query().Get("q")

	decoder := json.NewDecoder(strings.NewReader(jsonqlString))
	query := &jsonql.Query{}
	err := decoder.Decode(query)
	if err != nil {
		log.Fatal(err)
	}

	records, err := h.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	// FIXME: what do do with io error? Probably just kill the client.
	first := true
	res.Write([]byte("["))
	for records.Next() {
		if !first {
			res.Write([]byte(","))
		}
		first = false
		err := records.WriteJSON(res)
		if err != nil {
			break
		}
	}
	res.Write([]byte("]"))

	if err != nil {
		panic(err)
	}

	err = records.Err()
	if err != nil {
		panic(err)
	}
}
