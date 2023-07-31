docker-compose -f docker-compose.yml up -d
make createdb
make migratedbup
go build main.go
go build worker.go
./main &
./worker &
