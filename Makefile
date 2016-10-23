
all: build

build: deps
	govvv build -o "dtv-combust" main.go

deps:
	go get github.com/Masterminds/glide
	glide install

run: build
	./dtv-combust

consumer: .PHONY
	go run consumer/main.go

.PHONY:
