ALTER TABLE IF EXISTS "posts" DROP CONSTRAINT IF EXISTS "owner_title_key";

ALTER TABLE IF EXISTS "posts" DROP CONSTRAINT IF EXISTS "posts_owner_fkey";

DROP TABLE IF EXISTS "users";