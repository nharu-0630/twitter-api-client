build: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64

build-linux-amd64:
	env GOOS=linux GOARCH=amd64 go build -o ./bin/twitter-api-client_linux_amd64

build-linux-arm64:
	env GOOS=linux GOARCH=arm64 go build -o ./bin/twitter-api-client_linux_arm64

build-darwin-amd64:
	env GOOS=darwin GOARCH=amd64 go build -o ./bin/twitter-api-client_darwin_amd64

build-darwin-arm64:
	env GOOS=darwin GOARCH=arm64 go build -o ./bin/twitter-api-client_darwin_arm64

clean:
	rm -rf ./bin