package jsonql

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hayeah/jsonql/fixture"

	_ "github.com/mattn/go-sqlite3"
)

var db *DB

func TestMain(m *testing.M) {
	err := fixture.CreateDB()
	defer fixture.DestroyDB()
	if err != nil {
		log.Fatal(err)
	}

	db, err = OpenDB("sqlite3", fixture.TestDBDataSource)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func selectUsers() *Records {
	q := &Query{From: "users"}
	records, err := db.Query(q)
	if err != nil {
		panic(err.Error())
	}
	return records
}

func TestDbQuery(t *testing.T) {
	records := selectUsers()
	count := 0
	for records.Next() {
		count += 1
	}
	assert.NoError(t, records.Err())
	assert.Equal(t, 4, count, "Should iterate over 4 records")

	q := &Query{From: "user"}
	_, err := db.Query(q)
	assert.Error(t, err)
}

func TestDbRecords(t *testing.T) {
	records := selectUsers()

	// get data as map
	assert.True(t, records.Next())
	john := records.GetMap()
	assert.Equal(t, john["id"], 1)
	assert.Equal(t, john["username"], "john")
	_, isTime := john["created_at"].(time.Time)
	assert.True(t, isTime, "created_at should be a time.Time")

	// get data as json string
	jsonString := records.GetJSON()
	decoder := json.NewDecoder(strings.NewReader(jsonString))
	decodedJSON := make(map[string]interface{})
	err := decoder.Decode(&decodedJSON)
	assert.NoError(t, err)
	assert.Equal(t, john["id"], decodedJSON["id"], "JSON and record data should be the same")
	assert.Equal(t, john["username"], decodedJSON["username"], "JSON and record data should be the same")
}

func TestDbRecordsCopyMap(t *testing.T) {
	records := selectUsers()

	records.Next()
	john := records.CopyMap()

	records.Next()
	assert.Equal(t, john["id"], 1)
	assert.Equal(t, john["username"], "john")
}
