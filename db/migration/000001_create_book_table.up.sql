CREATE TABLE IF NOT EXISTS "book" (
  "id" varchar PRIMARY KEY,
  "name" varchar,
  "year" int,
  "author" varchar,
  "summary" varchar,
  "publisher" varchar,
  "page_count" int,
  "read_page" int,
  "finished" bool,
  "reading" bool,
  "created_at" varchar,
  "updated_at" varchar
);
