build:
	env GOOS=linux GOARCH=amd64 go build -o ./bin/twitter-api-client_linux_amd64
	env GOOS=darwin GOARCH=amd64 go build -o ./bin/twitter-api-client_darwin_amd64
	env GOOS=darwin GOARCH=arm64 go build -o ./bin/twitter-api-client_darwin_arm64