package jsonql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratorFrom(t *testing.T) {
	q := &Query{From: "foo"}
	sql, err := q.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM foo;", sql)
}

func TestGeneratorFromEmpty(t *testing.T) {
	q := &Query{From: ""}
	_, err := q.ToSql()
	assert.Error(t, err)
}

func TestGeneratorSelect(t *testing.T) {
	q := &Query{From: "foo", Select: []string{"id", "username"}}
	sql, err := q.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT id, username FROM foo;", sql)
}

func TestGeneratorWhere(t *testing.T) {
	q := &Query{From: "foo", Where: "id = 1 or id = 3"}
	sql, err := q.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM foo WHERE id = 1 or id = 3;", sql)
}

func TestGeneratorLimit(t *testing.T) {
	q := &Query{From: "foo", Limit: 2}
	sql, err := q.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM foo LIMIT 2;", sql)
}

func TestGeneratorOrder(t *testing.T) {
	q := &Query{From: "foo", Order: "id DESC"}
	sql, err := q.ToSql()
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM foo ORDER BY id DESC;", sql)
}

func TestGeneratorRelate(t *testing.T) {
	var r *RelateQuery
	var sql string
	var err error

	// error if Using is not specified
	r = &RelateQuery{}
	_, err = r.ToSql("followers", 42)
	assert.Error(t, err)

	r = &RelateQuery{Using: "user_id"}
	sql, err = r.ToSql("followers", 42)
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM followers WHERE user_id = 42;", sql)

	// merge where
	r = &RelateQuery{Using: "user_id", Where: "id > 10"}
	sql, err = r.ToSql("followers", 42)
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM followers WHERE (id > 10) AND user_id = 42;", sql)

	// override relation name
	r = &RelateQuery{From: "followers", Using: "user_id"}
	sql, err = r.ToSql("foobar", 42)
	assert.NoError(t, err)
	assert.Equal(t, "SELECT * FROM followers WHERE user_id = 42;", sql)
}
