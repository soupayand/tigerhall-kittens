postgres:
	docker run --name postgres	-p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15

kafka:
	docker run -d --name kafka -p 9092:9092 -e KAFKA_CFG_ZOOKEEPER_CONNECT=zookeeper:2181 -e ALLOW_PLAINTEXT_LISTENER=yes bitnami/kafka:3.5.1

dropdb:
	docker exec -it postgres dropdb -f --username root tigerhall_kittens

createdb:
		docker exec -it postgres createdb --username=root --owner=root tigerhall_kittens

migratedbup:
	migrate -path database/migration -database "postgresql://root:secret@localhost:5432/tigerhall_kittens?sslmode=disable" -verbose up

migratedbdown:
	migrate -path database/migration -database "postgresql://root:secret@localhost:5432/tigerhall_kittens?sslmode=disable" -verbose down


.PHONY: postgres	dropdb	createdb	migratedbup	migratedbdown kafka

