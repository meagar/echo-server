.PHONY: build
build:
	docker build --platform linux/amd64 . -t meagar/echo-server:latest

.PHONY: push
push:
	docker push meagar/echo-server:latest