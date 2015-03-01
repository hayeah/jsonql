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

func TestDbRelation(t *testing.T) {
	relations := make(map[string]RelateQuery)
	relations["followers"] = RelateQuery{Using: "user_id", Order: "id"}
	relations["followings"] = RelateQuery{Using: "user_id", Order: "id"}
	q := &Query{From: "users", Relate: relations, Where: "id = 1"}

	records, err := db.Query(q)
	assert.NoError(t, err)
	assert.True(t, records.Next())

	user := records.GetMap()
	assert.Equal(t, 1, user["id"])

	followers := user["followers"].([]map[string]interface{})
	followings := user["followings"].([]map[string]interface{})
	assert.Equal(t, 2, len(followers))
	assert.Equal(t, 1, len(followings))

	var u map[string]interface{}
	u = followers[0]
	assert.Equal(t, 2, u["id"])
	assert.Equal(t, "jerry", u["username"])

	u = followers[1]
	assert.Equal(t, 3, u["id"])
	assert.Equal(t, "jenny", u["username"])

	u = followings[0]
	assert.Equal(t, 4, u["id"])
	assert.Equal(t, "jane", u["username"])

	// fmt.Println("relate:", records.GetJSON())

}

func TestDbGetAll(t *testing.T) {
	q := &Query{From: "users"}
	records, err := db.Query(q)
	assert.NoError(t, err)

	maps, err := records.All()
	assert.NoError(t, err)
	assert.Equal(t, 4, len(maps))

	// null result should return empty slice instead of nil
	q = &Query{From: "users", Where: "1 = 0"}
	records, err = db.Query(q)
	assert.NoError(t, err)

	maps, err = records.All()
	assert.NoError(t, err)
	assert.NotNil(t, maps)
	assert.Equal(t, 0, len(maps))

}

func TestDbRelationParentKey(t *testing.T) {
	// The implicit parent key of a relation is "id".
	relations := make(map[string]RelateQuery)
	relations["followers"] = RelateQuery{Using: "user_id", Order: "id"}
	q := &Query{Select: []string{"username"}, From: "users", Relate: relations}

	records, err := db.Query(q)
	assert.NoError(t, err)
	assert.False(t, records.Next())
	assert.Equal(t, records.Err(), ErrRelationNoParentId)

	// can specify an alternate parent join key
	relations = make(map[string]RelateQuery)
	relations["followers"] = RelateQuery{Using: "user_id", Order: "id", ParentKey: "foobar"}
	q = &Query{Select: []string{"id as foobar"}, From: "users", Relate: relations}
	records, err = db.Query(q)
	assert.NoError(t, err)
	assert.True(t, records.Next())
}
