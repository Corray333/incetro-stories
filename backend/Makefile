.SILENT:
build:
	cd cmd && go build main.go
lint:
	golangci-lint run
run: build
	cd cmd && ./main ../.env
goose-up:
	cd migrations && goose postgres "user=postgres password=root host=db port=5432 dbname=stories sslmode=disable" up
goose-down:
	cd migrations && goose postgres "user=postgres password=root host=db port=5432 dbname=stories sslmode=disable" down
goose-down-all:
	cd migrations && goose postgres "user=postgres password=root host=db port=5432 dbname=stories sslmode=disable" down-to 0
