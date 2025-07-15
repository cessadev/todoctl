BINARY=todoctl

build:
    go build -o ${BINARY} main.go

run:
    go run main.go

test:
    go test ./...

clean:
    rm -f ${BINARY}