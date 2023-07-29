postgres:
	docker run --name postgres15	-p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15

dropdb:
	docker exec -it postgres15 dropdb -f --username root tigerhall_kittens

createdb:
		docker exec -it postgres15 createdb --username=root --owner=root tigerhall_kittens

migratedbup:
	migrate -path database/migration -database "postgresql://root:secret@localhost:5432/tigerhall_kittens?sslmode=disable" -verbose up

migratedbdown:
	migrate -path database/migration -database "postgresql://root:secret@localhost:5432/tigerhall_kittens?sslmode=disable" -verbose down


.PHONY: postgres	dropdb	createdb	migratedbup	migratedbdown

