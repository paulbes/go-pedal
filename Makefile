cli-build:
	docker build --build-arg command=cmd/pedal/main.go --tag pedal-cli .
.PHONY: cli-build

cli-run:
	docker run pedal-cli:latest -client-identifier $(CLIENT_IDENTIFIER)
.PHONY: cli-run

api-build:
	docker build --build-arg command=cmd/api/main.go --tag pedal-api .
.PHONY: api-build

api-run:
	docker run --expose 8080 --publish 8080:8080 pedal-api:latest -client-identifier $(CLIENT_IDENTIFIER)
.PHONY: api-run

include project.mk
include base.mk
