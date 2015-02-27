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
