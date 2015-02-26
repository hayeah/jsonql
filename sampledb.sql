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
INSERT INTO "users" VALUES(4,'jane','2015-02-1 13:47:08');

COMMIT;