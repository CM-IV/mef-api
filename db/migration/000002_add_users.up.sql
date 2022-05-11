CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
  "id" uuid DEFAULT uuid_generate_v4(),
  "user_name" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY (id)
);

ALTER TABLE "posts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("user_name");

-- CREATE UNIQUE INDEX ON "posts" ("title", "owner");
ALTER TABLE "posts" ADD CONSTRAINT "owner_title_key" UNIQUE ("title", "owner");