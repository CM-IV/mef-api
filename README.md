# Monero Economic Forum API
## Written in *Gin* Go Web Framework

By [Charlie](https://civdev.rocks) (CM-IV)

Dockerized PostgreSQL database used featuring sqlc type-safe Go generator.

---

**Table of Contents**

1. [CRUD Functionality Quick Peek](#crud-functionality-quick-peek) (You are here)
2. [Makefile Usage](#makefile-usage)
3. [PostgreSQL Queries and Migrations](#postgresql-queries-and-migrations)
4. [sqlc Generator](#sqlc-generator)
5. [Custom CRUD Implementation](#custom-crud-implementation)
6. [Unit Tests and RSG Util](#unit-tests-and-rsg-util)

---

### CRUD Functionality Quick Peek

HTTP protocol used to exchange representations of resources between the client frontend and the server API backend.  Post data is retrieved and accessed using URIs, here are the endpoints along with their operations:

*server.go*
```
router.POST("/api/posts", server.createPost)
router.GET("/api/posts/:id", server.getPost)
router.GET("/api/posts", server.listPost)
router.PUT("/api/posts/:id", server.updatePost)
router.DELETE("/api/posts/:id", server.deletePost)
```

---

### Makefile Usage

The Makefile contains all the needed commands for everything from bringing up the database (DB) to seeding it with data and performing up/down migrations.

Bring up the database container with the Makefile command `make postgres-up`.

The container can be stopped with the `make postgres-stop` command.  This frees up the terminal to be used for other things once the DB is restarted with `make postgres-start`.

Create the Postgres DB itself with a root username and owner.  The DB is named "meforum", this command can be utilized with `make createdb`.  In order to drop the DB, `make dropdb` is used.

Additional commands for DB migrations, sqlc code generation, and testing will be explained later on in their respective sections.

*Makefile*

```
postgres-up:
	docker-compose up

postgres-start:
	docker-compose start

postgres-stop:
	docker-compose stop	

postgres-down:
	docker-compose down

createdb:
	docker exec -it mef-api_db_1 createdb --username=root --owner=root meforum

dropdb:
	docker exec -it mef-api_db_1 dropdb meforum

migrateup:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/meforum?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/meforum?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -count=1 -v ./db/sqlc

server:
	go run main.go


.PHONY: postgres-up postgres-down postgres-start postgres-stop createdb dropdb migrateup migratedown sqlc test server
```

---

### PostgreSQL Queries and Migrations

**Migrations**

The DB up/down Migrations are an automated and useful way of updating or changing the DB SQL tables themselves.

Shown within the Makefile in the previous section, the DB table is created with the `make migrateup` command and the table can be dropped with the `make migratedown` command.  The migration commands have the path to the local Postgres instance with the ports used within them.

The up migration creates the "posts" table along with a "title" index:

*000001_init_schema.up.sql*
```
CREATE TABLE "posts" (
  "id" bigserial PRIMARY KEY,
  "image" varchar NOT NULL,
  "title" varchar NOT NULL,
  "subtitle" varchar NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "posts" ("title");
```

This will be extended to include tables for the other pages (Agenda page) of the website, and custom CRUD functionality will be implemented as well.

Some sort of user entity with admin permissions would most likely be implemented as well so that he/she can create the posts.

Currently each post shows the timestamp along with the timezone of when it was created.  This appears on each card within the frontend client on the homepage.

**Queries**

The [sqlc documentation](https://docs.sqlc.dev/en/latest/tutorials/getting-started.html) was very helpful in the creation of the SQL queries.  Sqlc itself will be expanded on in the next section, but this should be mentioned.  The `CreatePost`, `GetPost`, `ListPosts`, `UpdatePost`, and `DeletePost` SQL queries are all included here:


*post.sql*
```
-- name: CreatePost :one
INSERT INTO posts (
  image,
  title,
  subtitle,
  content
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;


-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1 LIMIT 1;

-- name: ListPosts :many
SELECT * FROM posts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdatePost :one
UPDATE posts
SET content = $2
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;

```

Shown within this file are the crucial comments that help sqlc do its work in generating the type-safe and idiomatic interfaces to the SQL queries.  The query annotations `one`, `many`, and `exec` tell sqlc that the query returns that many rows.

---

### sqlc Generator

Sqlc is configured using the sqlc.yaml file which must be in the directory that the sqlc command itself is run.  See their [config reference documentation](https://docs.sqlc.dev/en/latest/reference/config.html) for more details about what each key does.

*sqlc.yaml*
```
version: "1"
packages:
  - name: "db"
    path: "./db/sqlc"
    queries: "./db/query/"
    schema: "./db/migration/"
    engine: "postgresql"
    emit_json_tags: true
    emit_prepared_queries: false
    emit_interface: false
    emit_exact_table_names: false
    emit_empty_slices: true

```

The aforementioned sqlc [docs](https://docs.sqlc.dev/en/latest/howto/select.html) shed light on a Queries struct with DB access methods which is created using the `New` method.  This is located within the db.go file.

*db.go*
```
package db

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}
```

The pointer to sql.DB `*sql.DB` is within the store.go file to be used with the `NewStore(db *sql.DB)` method which returns a pointer reference in the Store instance that is created.

*store.go*
```
package db

import (
	"database/sql"
)

//Store will allow DB execute queries and transactions for all functions
//Composition extending struct functionality
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {

	return &Store{

		db:      db,
		Queries: New(db),
	}

}

```

The models.go file generated by sqlc sets up the type Post struct and the resulting json formatting that comes with it.  The json is configured to use lowercase representations of the rows.

*models.go*
```
package db

import (
	"time"
)

type Post struct {
	ID        int64     `json:"id"`
	Image     string    `json:"image"`
	Title     string    `json:"title"`
	Subtitle  string    `json:"subtitle"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
```

The post.sql.go file that interfaces with the PostgreSQL DB is generated along with the previously shown files.  This go code is both type-safe and idiomatic, allowing the programmer to write his/her own custom API application code.  The `CreatePost`, `GetPost`, `ListPosts`, `UpdatePosts`, and `DeletePost` is within this file.

*post.sql.go*
```
// source: post.sql

package db

import (
	"context"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts (
  image,
  title,
  subtitle,
  content
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, image, title, subtitle, content, created_at
`

type CreatePostParams struct {
	Image    string `json:"image"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Content  string `json:"content"`
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.Image,
		arg.Title,
		arg.Subtitle,
		arg.Content,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Image,
		&i.Title,
		&i.Subtitle,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}

const deletePost = `-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1
`

func (q *Queries) DeletePost(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deletePost, id)
	return err
}

const getPost = `-- name: GetPost :one
SELECT id, image, title, subtitle, content, created_at FROM posts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPost(ctx context.Context, id int64) (Post, error) {
	row := q.db.QueryRowContext(ctx, getPost, id)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Image,
		&i.Title,
		&i.Subtitle,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}

const listPosts = `-- name: ListPosts :many
SELECT id, image, title, subtitle, content, created_at FROM posts
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListPostsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListPosts(ctx context.Context, arg ListPostsParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, listPosts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Post{}
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.Image,
			&i.Title,
			&i.Subtitle,
			&i.Content,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePost = `-- name: UpdatePost :one
UPDATE posts
SET content = $2
WHERE id = $1
RETURNING id, image, title, subtitle, content, created_at
`

type UpdatePostParams struct {
	ID      int64  `json:"id"`
	Content string `json:"content"`
}

func (q *Queries) UpdatePost(ctx context.Context, arg UpdatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, updatePost, arg.ID, arg.Content)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Image,
		&i.Title,
		&i.Subtitle,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}
```