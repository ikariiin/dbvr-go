BINARY_NAME=dbvr

build:
	go build -o bin/${BINARY_NAME} main.go

run: build
	./bin/${BINARY_NAME}

clean-bin:
	rm bin/*

clean: clean-bin
	go clean