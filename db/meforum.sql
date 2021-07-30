CREATE TABLE "posts" (
  "id" bigserial PRIMARY KEY,
  "image" varchar NOT NULL,
  "title" varchar UNIQUE NOT NULL,
  "subtitle" varchar NOT NULL,
  "content" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "posts" ("title");
