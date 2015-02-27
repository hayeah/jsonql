package jsonql

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

var (
	ErrEmptyTableName = errors.New("Table name cannot be empty")
)

type Query struct {
	From   string   `json:"from"`
	Select []string `json:"select"`
	Where  string   `json:"where"`
	Limit  int      `json:"limit"`
}

// can't bind column and table names as parameters... need to escape. Can't find an escape function offhand.
// FIXME: prevent SQL injection, yo
func (q *Query) ToSql() (string, error) {
	sqlBuf := bytes.NewBufferString("SELECT ")

	// SELECT ...
	if q.Select == nil {
		sqlBuf.WriteString("*")
	} else {
		// Probably don't want to escape column names? this would allow `count(*)` in SELECT
		sqlBuf.WriteString(strings.Join(q.Select, ", "))
	}

	if q.From == "" {
		return "", ErrEmptyTableName
	}

	// FROM
	sqlBuf.WriteString(" FROM ")
	sqlBuf.WriteString(q.From)

	// WHERE
	if q.Where != "" {
		sqlBuf.WriteString(" WHERE ")
		sqlBuf.WriteString(q.Where)
	}

	if q.Limit != 0 {
		sqlBuf.WriteString(" LIMIT ")
		sqlBuf.WriteString(strconv.Itoa(q.Limit))
	}

	sqlBuf.WriteString(";")

	return sqlBuf.String(), nil
}
