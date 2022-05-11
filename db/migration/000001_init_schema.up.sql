CREATE TABLE "posts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "image" varchar NOT NULL,
  "title" varchar NOT NULL,
  "subtitle" varchar NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "posts" ("title");