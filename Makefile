build:
	docker build --tag pedal .
.PHONY: build

run:
	docker run pedal:latest -client-identifier $(CLIENT_IDENTIFIER)
.PHONY: run