BINARY_NAME='mscache'

build:
	go build -o bin/${BINARY_NAME}

runleader: build 
	./bin/${BINARY_NAME} --listenaddr :3000

runfollower: build 
	./bin/${BINARY_NAME} --listenaddr :4000 --leaderaddr :3000 

test: 
	go test -v ./...