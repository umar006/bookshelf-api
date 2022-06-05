postgres:
	docker run --name bookshelf-db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:14.0-alpine

createdb:
	docker exec -it bookshelf-db createdb bookshelf

dropdb:
	docker exec -it bookshelf-db dropdb bookshelf

migrateup:
	migrate -database "postgres://root:root@localhost:5432/bookshelf?sslmode=disable" -path db/migration up

migratedown:
	migrate -database "postgres://root:root@localhost:5432/bookshelf?sslmode=disable" -path db/migration down

.PHONY: postgres createdb dropdb migrateup migratedown
