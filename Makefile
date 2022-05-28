BINARY_NAME=mqtts-test

build:
	GOARCH=arm64 GOOS=linux go build -o ${BINARY_NAME}-linux
	GOARCH=arm64 GOOS=darwin go build -o ${BINARY_NAME}-darwin

clean:
	go clean
	rm -f ${BINARY_NAME}-linux
	rm -f ${BINARY_NAME}-darwin