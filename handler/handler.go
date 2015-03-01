package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/hayeah/jsonql"
)

var (
	ErrJSONParse = errors.New("Error parsing JSON")
)

// type errorResponse struct {
// 	err error
// }

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
		replyError(res, err)
		return
	}

	records, err := h.db.Query(query)
	if err != nil {
		replyError(res, err)
		return
	}

	// FIXME: what do do with io error? Probably just kill the client.
	first := true
	res.Write([]byte("[\n"))
	for records.Next() {
		if first {
			res.Write([]byte(" "))
		} else {
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

func replyError(res http.ResponseWriter, err error) {
	res.WriteHeader(400)
	encoder := json.NewEncoder(res)
	encoder.Encode(&struct {
		Error string
	}{
		Error: err.Error(),
	})
}
