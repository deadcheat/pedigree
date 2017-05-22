#! /usr/bin/make
test:
	go test -cover -v `glide novendor`