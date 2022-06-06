CREATE TABLE IF NOT EXISTS "book" (
  "id" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "year" int NOT NULL,
  "author" varchar NOT NULL,
  "summary" varchar,
  "publisher" varchar,
  "page_count" int,
  "read_page" int,
  "finished" bool,
  "reading" bool,
  "created_at" varchar DEFAULT (extract(epoch from now()) * 1000),
  "updated_at" varchar DEFAULT (extract(epoch from now()) * 1000)
);
