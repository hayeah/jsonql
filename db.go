package jsonql

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
)

type DB struct {
	db *sql.DB
}

func OpenDB(driverName string, dataSourceName string) (*DB, error) {
	sqldb, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db: sqldb}, nil
}

func (d *DB) Close() {
	d.db.Close()
}

func (d *DB) Query(q *Query) (records *Records, err error) {
	sql, err := q.ToSql()
	if err != nil {
		return
	}

	rows, err := d.db.Query(sql)
	if err != nil {
		return
	}

	records = &Records{rows: rows}
	return
}

type Records struct {
	rows *sql.Rows
	// Map to scan row data into
	dest    map[string]interface{}
	scanErr error
}

func (r *Records) Next() bool {
	next := r.rows.Next()
	if !next {
		return false
	}

	r.scanErr = r.scanData()
	return r.scanErr == nil
}

func (r *Records) scanData() error {
	if r.dest == nil {
		r.dest = make(map[string]interface{})
	}

	err := MapScan(r.rows, r.dest)
	return err
}

// The returned *Record is only valid until the call to Next()
func (r *Records) GetJSON() string {
	buf := bytes.NewBufferString("")
	err := r.WriteJSON(buf)
	if err != nil {
		panic(err.Error())
	}

	return buf.String()
}

// valid until Next()
func (r *Records) GetMap() map[string]interface{} {
	return r.dest
}

func (r *Records) CopyMap() map[string]interface{} {
	map2 := make(map[string]interface{})
	for k, v := range r.dest {
		map2[k] = v
	}
	return map2
}

func (r *Records) WriteJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(r.dest)
	return err
}

func (r *Records) Err() error {
	if err := r.rows.Err(); err != nil {
		return err
	}
	return r.scanErr
}

// Adapted from https://github.com/jmoiron/sqlx/blob/69738bd209812c5a381786011aff17944a384554/sqlx.go#L777
//
// MapScan scans a single Row into the dest map[string]interface{}.
// Use this to get results for SQL that might not be under your control
// (for instance, if you're building an interface for an SQL server that
// executes SQL from input).  Please do not use this as a primary interface!
// This will modify the map sent to it in place, so reuse the same map with
// care.  Columns which occur more than once in the result will overwrite
// eachother!
func MapScan(r *sql.Rows, dest map[string]interface{}) error {
	// ignore r.started, since we needn't use reflect for anything.
	columns, err := r.Columns()
	if err != nil {
		return err
	}

	// TODO can reuse this slice
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	err = r.Scan(values...)
	if err != nil {
		return err
	}

	for i, column := range columns {
		val := *(values[i].(*interface{}))
		switch val := val.(type) {
		// FIXME: sqlite3 returns []byte for TEXT. Cannot distinguish between Text and BLOB
		// will need to inspect schema to decide whether casting string is necessary.
		// JSON encoder uses base64 to encode []byte
		case []byte:
			dest[column] = string(val)
		default:
			dest[column] = val
		}
	}

	return r.Err()
}
