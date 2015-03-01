package jsonql

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrEmptyTableName = errors.New("Table name cannot be empty")
	ErrEmptyJoinKey   = errors.New("Must specify a join key with `using`")
)

type Query struct {
	From   string                 `json:"from"`
	Select []string               `json:"select"`
	Where  string                 `json:"where"`
	Limit  int                    `json:"limit"`
	Order  string                 `json:"order"`
	Relate map[string]RelateQuery `json:"relate"`
}

type RelateQuery struct {
	// same fields as Query
	From   string   `json:"from"`
	Select []string `json:"select"`
	Where  string   `json:"where"`
	Limit  int      `json:"limit"`
	Order  string   `json:"order"`

	// join related fields
	ParentKey string `json:"parent_key"`
	Using     string `json:"using"`
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

	if q.Order != "" {
		sqlBuf.WriteString(" ORDER BY ")
		sqlBuf.WriteString(q.Order)
	}

	if q.Limit != 0 {
		sqlBuf.WriteString(" LIMIT ")
		sqlBuf.WriteString(strconv.Itoa(q.Limit))
	}

	sqlBuf.WriteString(";")

	return sqlBuf.String(), nil
}

// Convert a relation query to an ordinary select
func (r *RelateQuery) ToQuery(relationName string, parentId int64) (*Query, error) {
	q := &Query{}

	// must specify join key
	if r.Using == "" {
		return nil, ErrEmptyJoinKey
	}

	q.From = r.From
	if r.From == "" {
		q.From = relationName
	}

	// join condition
	joinCond := fmt.Sprintf("%s = %d", r.Using, parentId)
	// merge join condition with relation's where condition
	if r.Where != "" {
		q.Where = fmt.Sprintf("(%s) AND %s", r.Where, joinCond)
	} else {
		q.Where = joinCond
	}

	// Copy other fields
	q.Limit = r.Limit
	q.Order = r.Order

	return q, nil
}

func (r *RelateQuery) ToSql(relationName string, parentId int64) (string, error) {
	q, err := r.ToQuery(relationName, parentId)
	if err != nil {
		return "", err
	}
	return q.ToSql()
}
