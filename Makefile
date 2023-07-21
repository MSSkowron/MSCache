BINARY_NAME=mscache

buildserver:
	go build -o server/bin/${BINARY_NAME} ./server/cmd

runleader: buildserver 
	./server/bin/${BINARY_NAME} --listenaddr :3000

runfollower: buildserver 
	./server/bin/${BINARY_NAME} --listenaddr :4000 --leaderaddr :3000 

buildclient:
	go build -o client/bin/${BINARY_NAME} ./client/runtest

runclient: buildclient
	./client/bin/${BINARY_NAME} --endpoint :3000

test: 
	go test ./...

lint:
	staticcheck ./...
	golint ./...

clean:
	@rm -rf ./server/bin/$(BINARY_NAME)
	@rm -rf ./client/bin/$(BINARY_NAME)

.PHONY: build runleader runfollower test lint clean