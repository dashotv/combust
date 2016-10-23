
all: build

build:
	go build -o "dtv-torrents" main.go client.go

deps:
	glide install

run:
	go run main.go

consumer: .PHONY
	go run consumer/main.go

.PHONY:
