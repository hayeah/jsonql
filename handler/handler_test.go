package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hayeah/jsonql/fixture"

	"github.com/hayeah/jsonql"
)

var db *jsonql.DB

var handler *HTTPHandler

var TestAddress = ":8899"

func TestMain(m *testing.M) {
	err := fixture.CreateDB()
	defer fixture.DestroyDB()
	if err != nil {
		log.Fatal(err)
	}

	db, err = jsonql.OpenDB("sqlite3", fixture.TestDBDataSource)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	handler = NewHTTPHandler(db)

	os.Exit(m.Run())
}

func TestHttpHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "http://test.local/", nil)
	assert.NoError(t, err)
	vals := make(url.Values)
	vals.Set("q", `{"from": "users"}`)
	req.URL.RawQuery = vals.Encode()

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	body := w.Body.String()
	// log.Println("body:\n", )
	decoder := json.NewDecoder(strings.NewReader(body))
	users := make([]interface{}, 4)
	decoder.Decode(&users)
	assert.Len(t, users, 4, "Should return 4 json objects")
	for _, user := range users {
		assert.NotNil(t, user)
		// log.Println(user)
	}
}
