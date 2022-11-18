validate:
	go mod tidy -v && \
	go mod verify && \
	golangci-lint run