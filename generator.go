package jsonql

import (
	"bytes"
	"errors"
	"strings"
)

var (
	ErrEmptyTableName = errors.New("Table name cannot be empty")
)

type Query struct {
	From   string   `json:"from"`
	Select []string `json:"select"`
}

// can't bind column and table names as parameters... need to escape. Can't find an escape function offhand.
// FIXME: prevent SQL injection, yo
func (q *Query) ToSql() (string, error) {
	sqlBuf := bytes.NewBufferString("SELECT ")

	if q.Select == nil {
		sqlBuf.WriteString("*")
	} else {
		sqlBuf.WriteString(strings.Join(q.Select, ", "))
	}

	if q.From == "" {
		return "", ErrEmptyTableName
	}
	// FROM
	sqlBuf.WriteString(" FROM ")
	sqlBuf.WriteString(q.From)

	sqlBuf.WriteString(";")

	return sqlBuf.String(), nil
}
