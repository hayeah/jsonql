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

CREATE TABLE follows (
  user_id integer not null,
  follower_id integer not null,
  created_at datetime not null,
  UNIQUE (user_id, follower_id)
);

INSERT INTO follows VALUES
(1, 2,'2015-03-01 13:47:08'),
(1, 3,'2015-03-02 13:47:08'),
(2, 3,'2015-03-03 13:47:08'),
(4, 3,'2015-03-04 13:47:08'),
(4, 1,'2015-03-05 13:47:08');

-- followers - those following a user
-- user_id, followed_at, id, username, created_at
CREATE VIEW followers AS
SELECT users.id as user_id, followers.*, follows.created_at as followed_at FROM users
INNER JOIN follows ON users.id = follows.user_id
INNER JOIN users as followers ON followers.id = follows.follower_id;

-- followings - those followed by a user
CREATE VIEW followings AS
SELECT users.id as user_id, followings.*, follows.created_at as followed_at FROM users
INNER JOIN follows ON users.id = follows.follower_id
INNER JOIN users as followings ON followings.id = follows.user_id;

COMMIT;