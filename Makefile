BINARY_NAME=mscache

buildserver:
	go build -o /bin/mscache/${BINARY_NAME} ./cmd/mscache/main.go

runleader: buildserver 
	./bin/mscache/${BINARY_NAME}  --listenaddr :3000

runfollower: buildserver 
	./bin/mscache/${BINARY_NAME} --listenaddr :4000 --leaderaddr :3000 

test: 
	go test ./...

lint:
	staticcheck ./...
	golint ./...

clean:
	@rm -rf ./server/bin/$(BINARY_NAME)
	@rm -rf ./client/bin/$(BINARY_NAME)

.PHONY: build runleader runfollower test lint clean
