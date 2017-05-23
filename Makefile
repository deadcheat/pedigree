#! /usr/bin/make
test:
	go test -cover -v `glide novendor`

build: deps
	go build -o pedigree main.go

deps:
	glide install