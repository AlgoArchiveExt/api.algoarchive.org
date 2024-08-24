all: clean build run

build:
	go build -o build/main.exe main.go

run:
	./build/main

clean:
	rm -rf ./build

clean-windows:
	rd "./build" /s /q
