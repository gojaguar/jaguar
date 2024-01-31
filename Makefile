tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet -v ./...

test:
	go test -race -covermode=atomic -coverprofile=coverage.tx -v ./...
	go tool cover -func=coverage.tx -o=coverage.out

test-html:
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o=coverage.html