GO ?= go

.PHONY: test
.DEFAULT_GOAL := test

test:
	$(GO) test -v ./...