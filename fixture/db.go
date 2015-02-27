package fixture

import (
	"os"
	"os/exec"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var TestDBSQL = `
PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;

CREATE TABLE users (
  id integer primary key autoincrement,
  username text unique not null,
  created_at datetime not null
);

INSERT INTO "users" VALUES(1,'john','2015-02-26 13:47:08');
INSERT INTO "users" VALUES(2,'jerry','2015-02-21 13:47:08');
INSERT INTO "users" VALUES(3,'jenny','2015-02-10 13:47:08');
INSERT INTO "users" VALUES(4,'jane','2015-02-01 13:47:08');

COMMIT;
`

var (
	TestDBDataSource = "testdb.sqlite3"
)

func CreateDB() error {
	var err error
	DestroyDB()
	cmd := exec.Command("sqlite3", TestDBDataSource)

	// pipe sql into db
	r := strings.NewReader(TestDBSQL)
	cmd.Stdin = r
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func DestroyDB() {
	os.Remove(TestDBDataSource)
}
