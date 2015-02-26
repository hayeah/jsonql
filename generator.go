package jsonql

import (
	"bytes"
	"errors"
)

var (
	ErrEmptyTableName = errors.New("Table name cannot be empty")
)

type Query struct {
	From string
}

// can't bind column and table names as parameters... need to escape. Can't find an escape function offhand.
// FIXME: prevent SQL injection, yo
func (q *Query) ToSql() (string, error) {
	sqlBuf := bytes.NewBufferString("SELECT *")

	if q.From == "" {
		return "", ErrEmptyTableName
	}
	// FROM
	sqlBuf.WriteString(" FROM ")
	sqlBuf.WriteString(q.From)

	sqlBuf.WriteString(";")

	return sqlBuf.String(), nil
}
