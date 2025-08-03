all: build run

rebuild: clean build run

build:
	go build -o main main.go

run:
	./main

clean:
	rm -rf ./build

clean-windows:
	rd "./build" /s /q
