# jsonQL

Inspired by GraphQL, build JSON documents from a declarative JSON-based spec.

The design goals are:

+ Obvious query generation. What you see is what you get.
+ Cover 90% of common sql query usage.
+ Anti-ORM: take advantage sql features like views.
  + Or denormalized data
+ Should be easy to programmatically generate jsonQL queries.
+ Straight-forward mapping to nested JSON structures.


# Select

Simple query:

```
// select * from users;
{from: "users"}
```

Select specific columns:

```
// select id, username from users;
{from: "users", select: ["id","username"]}
```

# Where

Where condition can take a string literal:

```
// select * from users where 1 < id and id < 4 and id != 3
{from: "users", where: "1 < id and id < 4 and id != 3"}
```

Or it can use reversed polish notation to build complex conditions:

```
// select * from users where 1 < id and id < 4 and id != 3
{from: "users", where: ["and",[">",1,"id","4"],"id != 3"]}
```

# Limit

```
// array of at most 10 element
{limit: 10}
// array of at most 1 element
{limit: 1}
// 1 element, or null
{limit: "first"}
```

# Named Parameter

(NOT IMPLEMENTED YET)

jsonQL doesn't escape anything. Consider it as potentially dangerous as writing raw SQL queries. To protect against SQL injection attack, you can use named parameters in `where`, like so:

```
["id > :lower_id and id < :upper_id"]
```

# Relations

Only support simple foreign key join on a single key. Use views to handle more complex join conditions.

```
{ from: "users",
  join: {
     followers: {limit: 10},
     followers_count: {table: "followers", select: "count(*)", limit: "first"},
     following: {from: "followers", key: "follower_id", limit: 10},
  }
}
```

Nested relations. Geting the latest 3 tweets of all following users:

```
{
  from: "users",
  relate: {
    "following": {
      relate: {"tweets": {limit: 3, order: "id desc"}}
    }
  },
  where: {"id = :id"}
}
```

Join is not used. A select query will be issued per relation per record.


