
postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=$${POSTGRES_USER} -e POSTGRES_PASSWORD=$${POSTGRES_PASSWORD} -d postgres

createdb:
	docker exec -it postgres createdb --username=$${POSTGRES_USER} --owner=$${POSTGRES_USER} $${POSTGRES_DB}

dropdb:
	docker exec -it postgres dropdb $${POSTGRES_DB}

migrateup:	
	migrate -path db/migration -database "$${DB_SOURCE}" -verbose up

migrateup1:
	migrate -path db/migration -database "$${DB_SOURCE}" -verbose up 1

migratedown:
	migrate -path db/migration -database "$${DB_SOURCE}" -verbose down

# Roll back only the last N migrations (e.g. make migratedown1 to undo only users migration)
migratedown1:
	migrate -path db/migration -database "$${DB_SOURCE}" -verbose down 1
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go	

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server