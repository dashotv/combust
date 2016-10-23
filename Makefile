
all: build

build:
	go build -o "dtv-combust" main.go

deps:
	glide install

run:
	go run main.go

consumer: .PHONY
	go run consumer/main.go

.PHONY:
