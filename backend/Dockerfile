FROM golang:alpine

RUN apk update && apk add --no-cache gcc libc-dev

WORKDIR /work

COPY go.mod go.sum ./

RUN go mod tidy

ADD . .

RUN go test -short ./...

RUN go build -o /bin/app cmd/app/main.go

WORKDIR /

RUN rm -r /work

COPY config/config.json /config/config.json

EXPOSE 4000

CMD ["/bin/app"]
