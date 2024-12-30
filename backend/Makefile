dep:
	go mod tidy

run:
	go run cmd/app/main.go

test:
	go test -short -cover ./...

build:
	go build -o bin/server cmd/app/main.go

docker-image:
	docker build -t server:v1 .

docker-run:
	docker run -it -d -p 3000:3000 --name server
