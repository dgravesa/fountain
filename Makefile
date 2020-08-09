#
build:
	mkdir -p bin
	go build -o bin ./...

install:
	go install ./...

clean:
	rm -r bin
